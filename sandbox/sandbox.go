package sandbox

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"code.cloudfoundry.org/winc/hcsclient"

	"github.com/Microsoft/hcsshim"
)

var sandboxFiles = []string{"Hives", "initialized", "sandbox.vhdx", "layerchain.json"}

type ImageSpec struct {
	RootFs string `json:"rootfs,omitempty"`
	Image  Image  `json:"image,omitempty"`
}

type Image struct {
	Config ImageConfig `json:"config,omitempty"`
}

type ImageConfig struct {
	Layers []string `json:"layers,omitempty"`
}

//go:generate counterfeiter . SandboxManager
type SandboxManager interface {
	Create(rootfsPath string) (*ImageSpec, error)
	Delete() error
	LayerFolderPath() string
	Mount(pid int, volumePath string) error
	Unmount(pid int) error
}

//go:generate counterfeiter . Mounter
type Mounter interface {
	SetPoint(string, string) error
	DeletePoint(string) error
}

type sandboxManager struct {
	hcsClient  hcsclient.Client
	id         string
	driverInfo hcsshim.DriverInfo
	mounter    Mounter
}

func NewManager(hcsClient hcsclient.Client, mounter Mounter, depotDir string, containerId string) SandboxManager {
	driverInfo := hcsshim.DriverInfo{
		HomeDir: depotDir,
		Flavour: 1,
	}

	return &sandboxManager{
		hcsClient:  hcsClient,
		mounter:    mounter,
		id:         containerId,
		driverInfo: driverInfo,
	}
}

func (s *sandboxManager) Create(rootfsPath string) (*ImageSpec, error) {
	parentLayerChain, err := ioutil.ReadFile(filepath.Join(rootfsPath, "layerchain.json"))
	if err != nil {
		return nil, err
	}

	parentLayers := []string{}
	if err := json.Unmarshal(parentLayerChain, &parentLayers); err != nil {
		return nil, &InvalidRootfsLayerChainError{Path: rootfsPath}
	}
	sandboxLayers := append([]string{rootfsPath}, parentLayers...)

	err = os.MkdirAll(s.driverInfo.HomeDir, 0755)
	if err != nil {
		return nil, err
	}

	if err := s.hcsClient.CreateSandboxLayer(s.driverInfo, s.id, rootfsPath, sandboxLayers); err != nil {
		return nil, err
	}

	if err := s.hcsClient.ActivateLayer(s.driverInfo, s.id); err != nil {
		return nil, err
	}

	if err := s.hcsClient.PrepareLayer(s.driverInfo, s.id, sandboxLayers); err != nil {
		return nil, err
	}

	volumePath, err := s.hcsClient.GetLayerMountPath(s.driverInfo, s.id)
	if err != nil {
		return nil, err
	} else if volumePath == "" {
		return nil, &hcsclient.MissingVolumePathError{Id: s.id}
	}

	sandboxLayers = append([]string{filepath.Join(s.driverInfo.HomeDir, s.id)}, sandboxLayers...)

	return &ImageSpec{
		RootFs: volumePath,
		Image: Image{
			Config: ImageConfig{
				Layers: sandboxLayers,
			},
		},
	}, nil
}

func (s *sandboxManager) Delete() error {
	if err := s.hcsClient.UnprepareLayer(s.driverInfo, s.id); err != nil {
		return err
	}

	if err := s.hcsClient.DeactivateLayer(s.driverInfo, s.id); err != nil {
		return err
	}

	for _, f := range sandboxFiles {
		layerFile := filepath.Join(s.driverInfo.HomeDir, f)
		if err := os.RemoveAll(layerFile); err != nil {
			return &UnableToDestroyLayerError{Msg: layerFile}
		}
	}

	return nil
}

func (s *sandboxManager) LayerFolderPath() string {
	return filepath.Join(s.driverInfo.HomeDir, s.id)
}

func (s *sandboxManager) mountPath(pid int) string {
	return filepath.Join("c:\\", "proc", strconv.Itoa(pid))
}

func (s *sandboxManager) rootPath(pid int) string {
	return filepath.Join(s.mountPath(pid), "root")
}

func (s *sandboxManager) Mount(pid int, volumePath string) error {
	if err := os.MkdirAll(s.rootPath(pid), 0755); err != nil {
		return err
	}

	return s.mounter.SetPoint(s.rootPath(pid), volumePath)
}

func (s *sandboxManager) Unmount(pid int) error {
	if err := s.mounter.DeletePoint(s.rootPath(pid)); err != nil {
		return err
	}

	return os.RemoveAll(s.mountPath(pid))
}
