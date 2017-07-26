package container

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"code.cloudfoundry.org/winc/hcsclient"
	"code.cloudfoundry.org/winc/network"
	"code.cloudfoundry.org/winc/sandbox"
	"github.com/Microsoft/hcsshim"
	"github.com/Sirupsen/logrus"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

const destroyTimeout = time.Minute

type ContainerManager interface {
	Create(spec *specs.Spec) error
	Delete() error
	State() (*specs.State, error)
	Exec(*specs.Process) (hcsshim.Process, error)
}

type containerManager struct {
	hcsClient      hcsclient.Client
	sandboxManager sandbox.SandboxManager
	networkManager network.NetworkManager
	id             string
	bundlePath     string
}

func NewManager(hcsClient hcsclient.Client, sandboxManager sandbox.SandboxManager, networkManager network.NetworkManager, bundlePath string) ContainerManager {
	return &containerManager{
		hcsClient:      hcsClient,
		sandboxManager: sandboxManager,
		networkManager: networkManager,
		bundlePath:     bundlePath,
		id:             filepath.Base(bundlePath),
	}
}

func (c *containerManager) Create(spec *specs.Spec) error {
	_, err := c.hcsClient.GetContainerProperties(c.id)
	if err == nil {
		return &hcsclient.AlreadyExistsError{Id: c.id}
	}
	if _, ok := err.(*hcsclient.NotFoundError); !ok {
		return err
	}

	imageSpec, err := c.sandboxManager.Create(spec.Root.Path)
	if err != nil {
		return err
	}

	layerInfos := []hcsshim.Layer{}
	for _, layerPath := range imageSpec.Image.Config.Layers {
		layerId := filepath.Base(layerPath)
		layerGuid, err := c.hcsClient.NameToGuid(layerId)
		if err != nil {
			return err
		}

		layerInfos = append(layerInfos, hcsshim.Layer{
			ID:   layerGuid.ToString(),
			Path: layerPath,
		})
	}

	mappedDirs := []hcsshim.MappedDir{}
	for _, d := range spec.Mounts {
		fileInfo, err := os.Stat(d.Source)
		if err != nil {
			if deleteErr := c.sandboxManager.Delete(); deleteErr != nil {
				logrus.Error(deleteErr.Error())
			}
			return err
		}
		if !fileInfo.IsDir() {
			logrus.WithField("mount", d.Source).Error("mount is not a directory, ignoring")
			continue
		}

		mappedDirs = append(mappedDirs, hcsshim.MappedDir{
			HostPath:      d.Source,
			ContainerPath: destToWindowsPath(d.Destination),
			ReadOnly:      true,
		})
	}

	containerConfig := hcsshim.ContainerConfig{
		SystemType:        "Container",
		Name:              c.bundlePath,
		VolumePath:        imageSpec.RootFs,
		Owner:             "winc",
		LayerFolderPath:   c.sandboxManager.LayerFolderPath(),
		Layers:            layerInfos,
		MappedDirectories: mappedDirs,
	}

	containerConfig, err = c.networkManager.AttachEndpointToConfig(containerConfig, c.id)
	if err != nil {
		if deleteErr := c.sandboxManager.Delete(); deleteErr != nil {
			logrus.Error(deleteErr.Error())
		}

		return err
	}

	if spec.Windows != nil {
		if spec.Windows.Resources != nil {
			if spec.Windows.Resources.Memory != nil {
				if spec.Windows.Resources.Memory.Limit != nil {
					memBytes := *spec.Windows.Resources.Memory.Limit
					containerConfig.MemoryMaximumInMB = int64(memBytes / 1024 / 1024)
				}
			}
		}
	}

	container, err := c.hcsClient.CreateContainer(c.id, &containerConfig)
	if err != nil {
		if deleteErr := c.sandboxManager.Delete(); deleteErr != nil {
			logrus.Error(deleteErr.Error())
		}
		if deleteErr := c.networkManager.DeleteEndpointsById(containerConfig.EndpointList, c.id); deleteErr != nil {
			logrus.Error(deleteErr.Error())
		}

		return err
	}

	if err := container.Start(); err != nil {
		if deleteErr := c.deleteContainer(container); deleteErr != nil {
			logrus.Error(deleteErr.Error())
		}
		return err
	}

	pid, err := c.containerPid(c.id)
	if err != nil {
		if deleteErr := c.deleteContainer(container); deleteErr != nil {
			logrus.Error(deleteErr.Error())
		}
		return err
	}

	if err := c.sandboxManager.Mount(pid, imageSpec.RootFs); err != nil {
		if deleteErr := c.deleteContainer(container); deleteErr != nil {
			logrus.Error(deleteErr.Error())
		}
		return err
	}

	return nil
}

func (c *containerManager) Delete() error {
	pid, err := c.containerPid(c.id)
	if err != nil {
		return err
	}

	unmountErr := c.sandboxManager.Unmount(pid)
	if unmountErr != nil {
		logrus.Error(unmountErr.Error())
	}

	container, err := c.hcsClient.OpenContainer(c.id)
	if err != nil {
		return err
	}

	err = c.deleteContainer(container)
	if err != nil {
		return err
	}

	return unmountErr
}

func (c *containerManager) State() (*specs.State, error) {
	cp, err := c.hcsClient.GetContainerProperties(c.id)
	if err != nil {
		return nil, err
	}

	var status string
	if cp.Stopped {
		status = "stopped"
	} else {
		status = "created"
	}

	pid, err := c.containerPid(c.id)
	if err != nil {
		return nil, err
	}

	return &specs.State{
		Version: specs.Version,
		ID:      c.id,
		Status:  status,
		Bundle:  c.bundlePath,
		Pid:     pid,
	}, nil
}

func (c *containerManager) Exec(processSpec *specs.Process) (hcsshim.Process, error) {
	container, err := c.hcsClient.OpenContainer(c.id)
	if err != nil {
		return nil, err
	}

	env := map[string]string{}
	for _, e := range processSpec.Env {
		v := strings.Split(e, "=")
		env[v[0]] = strings.Join(v[1:], "=")
	}

	pc := &hcsshim.ProcessConfig{
		CommandLine:      makeCmdLine(processSpec.Args),
		CreateStdInPipe:  true,
		CreateStdOutPipe: true,
		CreateStdErrPipe: true,
		WorkingDirectory: processSpec.Cwd,
		User:             processSpec.User.Username,
		Environment:      env,
	}
	p, err := container.CreateProcess(pc)
	if err != nil {
		command := ""
		if len(processSpec.Args) != 0 {
			command = processSpec.Args[0]
		}
		return nil, &hcsclient.CouldNotCreateProcessError{Id: c.id, Command: command}
	}

	return p, nil
}

func (c *containerManager) containerPid(id string) (int, error) {
	container, err := c.hcsClient.OpenContainer(id)
	if err != nil {
		return -1, err
	}

	pl, err := container.ProcessList()
	if err != nil {
		return -1, err
	}

	var process hcsshim.ProcessListItem
	oldestTime := time.Now()
	for _, v := range pl {
		if v.ImageName == "wininit.exe" && v.CreateTimestamp.Before(oldestTime) {
			oldestTime = v.CreateTimestamp
			process = v
		}
	}

	return int(process.ProcessId), nil
}

func (c *containerManager) deleteContainer(container hcsshim.Container) error {
	if err := c.networkManager.DeleteContainerEndpoints(container, c.id); err != nil {
		logrus.Error(err.Error())
	}

	if err := c.shutdownContainer(container); err != nil {
		if err := c.terminateContainer(container); err != nil {
			return err
		}
	}

	return c.sandboxManager.Delete()
}

func (c *containerManager) shutdownContainer(container hcsshim.Container) error {
	if err := container.Shutdown(); err != nil {
		if c.hcsClient.IsPending(err) {
			if err := container.WaitTimeout(destroyTimeout); err != nil {
				logrus.Error("hcsContainer.WaitTimeout error after Shutdown", err)
				return err
			}
		} else {
			logrus.Error("hcsContainer.Shutdown error", err)
			return err
		}
	}

	return nil
}

func (c *containerManager) terminateContainer(container hcsshim.Container) error {
	if err := container.Terminate(); err != nil {
		if c.hcsClient.IsPending(err) {
			if err := container.WaitTimeout(destroyTimeout); err != nil {
				logrus.Error("hcsContainer.WaitTimeout error after Terminate", err)
				return err
			}
		} else {
			logrus.Error("hcsContainer.Terminate error", err)
			return err
		}
	}

	return nil
}

func destToWindowsPath(input string) string {
	vol := filepath.VolumeName(input)
	if vol == "" {
		input = filepath.Join("C:", input)
	}
	return filepath.Clean(input)
}

func makeCmdLine(args []string) string {
	if len(args) > 0 {
		args[0] = filepath.Clean(args[0])
		base := filepath.Base(args[0])
		match, _ := regexp.MatchString(`\.[a-zA-Z]{3}$`, base)
		if !match {
			args[0] += ".exe"
		}
	}
	var s string
	for _, v := range args {
		if s != "" {
			s += " "
		}
		s += syscall.EscapeArg(v)
	}

	return s
}
