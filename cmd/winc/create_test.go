package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"code.cloudfoundry.org/winc/container"
	"code.cloudfoundry.org/winc/hcsclient"
	"code.cloudfoundry.org/winc/mounter"
	"code.cloudfoundry.org/winc/sandbox"

	"github.com/Microsoft/hcsshim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

var _ = Describe("Create", func() {
	var (
		config      []byte
		containerId string
		client      hcsclient.Client
		cm          container.ContainerManager
		bundleSpec  specs.Spec
		err         error
		stdOut      *bytes.Buffer
		stdErr      *bytes.Buffer
	)

	BeforeEach(func() {
		containerId = filepath.Base(bundlePath)

		client = &hcsclient.HCSClient{}
		sm := sandbox.NewManager(client, &mounter.Mounter{}, bundlePath)
		nm := networkManager(client)
		cm = container.NewManager(client, sm, nm, containerId)

		bundleSpec = runtimeSpecGenerator(rootfsPath)

		stdOut = new(bytes.Buffer)
		stdErr = new(bytes.Buffer)
	})

	JustBeforeEach(func() {
		config, err = json.Marshal(&bundleSpec)
		Expect(err).NotTo(HaveOccurred())

		Expect(ioutil.WriteFile(filepath.Join(bundlePath, "config.json"), config, 0666)).To(Succeed())
	})

	Context("when provided valid arguments", func() {
		AfterEach(func() {
			Expect(cm.Delete()).To(Succeed())
		})

		It("creates and starts a container", func() {
			err := exec.Command(wincBin, "create", "-b", bundlePath, containerId).Run()
			Expect(err).ToNot(HaveOccurred())

			Expect(containerExists(containerId)).To(BeTrue())

			sandboxVHDX := filepath.Join(bundlePath, "sandbox.vhdx")
			Expect(sandboxVHDX).To(BeAnExistingFile())

			err = exec.Command("powershell", "-Command", "Test-VHD", sandboxVHDX).Run()
			Expect(err).ToNot(HaveOccurred())

			sandboxInitialized := filepath.Join(bundlePath, "initialized")
			Expect(sandboxInitialized).To(BeAnExistingFile())

			layerChainPath := filepath.Join(bundlePath, "layerchain.json")
			Expect(layerChainPath).To(BeAnExistingFile())

			layerChain, err := ioutil.ReadFile(layerChainPath)
			Expect(err).ToNot(HaveOccurred())

			layers := []string{}
			err = json.Unmarshal(layerChain, &layers)
			Expect(err).ToNot(HaveOccurred())

			Expect(layers[0]).To(Equal(rootfsPath))

			state, err := cm.State()
			Expect(err).ToNot(HaveOccurred())
			Expect(state.Pid).ToNot(Equal(-1))
		})

		It("attaches a network endpoint with a port mapping", func() {
			err := exec.Command(wincBin, "create", "-b", bundlePath, containerId).Run()
			Expect(err).ToNot(HaveOccurred())

			endpoints := allEndpoints(containerId)
			Expect(len(endpoints)).To(Equal(1))

			endpoint, err := client.GetHNSEndpointByID(endpoints[0])
			Expect(err).To(Succeed())
			Expect(endpoint.Name).To(Equal(containerId))

			natPolicies := []hcsshim.NatPolicy{}
			for _, pol := range endpoint.Policies {
				natPolicy := hcsshim.NatPolicy{}

				err := json.Unmarshal(pol, &natPolicy)
				Expect(err).To(Succeed())
				if natPolicy.Type != "NAT" {
					continue
				}

				natPolicies = append(natPolicies, natPolicy)
			}

			Expect(len(natPolicies)).To(Equal(2))
			sort.Slice(natPolicies, func(i, j int) bool { return natPolicies[i].InternalPort < natPolicies[j].InternalPort })
			Expect(natPolicies[0].InternalPort).To(Equal(uint16(2222)))
			Expect(natPolicies[0].ExternalPort).To(BeNumerically(">=", 40000))
			Expect(natPolicies[0].Protocol).To(Equal("TCP"))

			Expect(natPolicies[1].InternalPort).To(Equal(uint16(8080)))
			Expect(natPolicies[1].ExternalPort).To(BeNumerically(">=", 40000))
			Expect(natPolicies[1].Protocol).To(Equal("TCP"))

			Expect(natPolicies[0].ExternalPort).NotTo(Equal(natPolicies[1].ExternalPort))
		})

		It("mounts the sandbox.vhdx at C:\\proc\\<pid>\\root", func() {
			err := exec.Command(wincBin, "create", "-b", bundlePath, containerId).Run()
			Expect(err).ToNot(HaveOccurred())

			state, err := cm.State()
			Expect(err).ToNot(HaveOccurred())

			Expect(ioutil.WriteFile(filepath.Join("c:\\", "proc", strconv.Itoa(state.Pid), "root", "test.txt"), []byte("contents"), 0644)).To(Succeed())
			cmd := exec.Command(wincBin, "exec", containerId, "cmd.exe", "/C", "type", "test.txt")
			cmd.Stdout = stdOut

			Expect(cmd.Run()).To(Succeed())

			Expect(stdOut.String()).To(ContainSubstring("contents"))
		})

		Context("when the bundle path is not provided", func() {
			It("uses the current directory as the bundle path", func() {
				cmd := exec.Command(wincBin, "create", containerId)
				cmd.Dir = bundlePath
				Expect(cmd.Run()).To(Succeed())

				state, err := cm.State()
				Expect(err).ToNot(HaveOccurred())
				Expect(state.Pid).ToNot(Equal(-1))
			})
		})

		Context("when the bundle path ends with a \\", func() {
			It("creates a container sucessfully", func() {
				err := exec.Command(wincBin, "create", "-b", bundlePath+"\\", containerId).Run()
				Expect(err).ToNot(HaveOccurred())

				state, err := cm.State()
				Expect(err).ToNot(HaveOccurred())
				Expect(state.Pid).ToNot(Equal(-1))
			})
		})

		Context("when the '--pid-file' flag is provided", func() {
			var pidFile string

			BeforeEach(func() {
				f, err := ioutil.TempFile("", "pid")
				Expect(err).ToNot(HaveOccurred())
				Expect(f.Close()).To(Succeed())
				pidFile = f.Name()
			})

			AfterEach(func() {
				Expect(os.RemoveAll(pidFile)).To(Succeed())
			})

			It("creates and starts the container and writes the container pid to the specified file", func() {
				err := exec.Command(wincBin, "create", "-b", bundlePath, "--pid-file", pidFile, containerId).Run()
				Expect(err).ToNot(HaveOccurred())

				state, err := cm.State()
				Expect(err).ToNot(HaveOccurred())
				Expect(state.Pid).ToNot(Equal(-1))

				pidBytes, err := ioutil.ReadFile(pidFile)
				Expect(err).ToNot(HaveOccurred())
				pid, err := strconv.ParseInt(string(pidBytes), 10, 64)
				Expect(err).ToNot(HaveOccurred())
				Expect(int(pid)).To(Equal(state.Pid))
			})
		})

		Context("when the '--no-new-keyring' flag is provided", func() {
			It("ignores it and creates and starts a container", func() {
				err := exec.Command(wincBin, "create", "-b", bundlePath, "--no-new-keyring", containerId).Run()
				Expect(err).ToNot(HaveOccurred())

				state, err := cm.State()
				Expect(err).ToNot(HaveOccurred())
				Expect(state.Pid).ToNot(Equal(-1))
			})
		})

		Context("when the bundle config.json specifies bind mounts", func() {
			var (
				mountSource string
				mountDest   string
			)

			BeforeEach(func() {
				var err error
				mountSource, err = ioutil.TempDir("", "mountsource")
				Expect(err).ToNot(HaveOccurred())
				Expect(ioutil.WriteFile(filepath.Join(mountSource, "sentinel"), []byte("hello"), 0644)).To(Succeed())

				mountDest = "C:\\mountdest"

				mount := specs.Mount{Destination: mountDest, Source: mountSource}
				bundleSpec.Mounts = []specs.Mount{mount}
			})

			AfterEach(func() {
				Expect(os.RemoveAll(mountSource)).To(Succeed())
			})

			It("creates a container with the specified directories as mounts", func() {
				err := exec.Command(wincBin, "create", "-b", bundlePath, containerId).Run()
				Expect(err).ToNot(HaveOccurred())

				cmd := exec.Command(wincBin, "exec", containerId, "cmd.exe", "/C", "type", filepath.Join(mountDest, "sentinel"))
				cmd.Stdout = stdOut
				Expect(cmd.Run()).To(Succeed())

				Expect(stdOut.String()).To(ContainSubstring("hello"))
			})

			It("the mounted directories are read only", func() {
				err := exec.Command(wincBin, "create", "-b", bundlePath, containerId).Run()
				Expect(err).ToNot(HaveOccurred())

				cmd := exec.Command(wincBin, "exec", containerId, "cmd.exe", "/C", "echo hello > "+filepath.Join(mountDest, "sentinel2"))
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session).Should(gexec.Exit(1))
			})

			Context("when the destination is /tmp/", func() {
				BeforeEach(func() {
					mountDest = "/tmp/mountdest"

					mount := specs.Mount{Destination: mountDest, Source: mountSource}
					bundleSpec.Mounts = []specs.Mount{mount}
				})
				It("mounts the specified directories", func() {
					err := exec.Command(wincBin, "create", "-b", bundlePath, containerId).Run()
					Expect(err).ToNot(HaveOccurred())

					cmd := exec.Command(wincBin, "exec", containerId, "cmd.exe", "/C", "type", filepath.Join(mountDest, "sentinel"))
					cmd.Stdout = stdOut
					Expect(cmd.Run()).To(Succeed())

					Expect(stdOut.String()).To(ContainSubstring("hello"))
				})

				Context("when calling the mounted executable", func() {
					BeforeEach(func() {
						Expect(copy(filepath.Join(mountSource, "cmd.exe"), "C:\\Windows\\System32\\cmd.exe")).To(Succeed())

					})
					Context("when using the windows path", func() {
						It("mounts the specified directories", func() {
							err := exec.Command(wincBin, "create", "-b", bundlePath, containerId).Run()
							Expect(err).ToNot(HaveOccurred())

							cmd := exec.Command(wincBin, "exec", containerId, filepath.Join(mountDest, "cmd"), "/C", "type", filepath.Join(mountDest, "sentinel"))
							cmd.Stdout = stdOut
							Expect(cmd.Run()).To(Succeed())

							Expect(stdOut.String()).To(ContainSubstring("hello"))
						})
					})
					Context("when using the unix path", func() {
						It("mounts the specified directories", func() {
							err := exec.Command(wincBin, "create", "-b", bundlePath, containerId).Run()
							Expect(err).ToNot(HaveOccurred())

							cmd := exec.Command(wincBin, "exec", containerId, mountDest+"/cmd", "/C", "type", filepath.Join(mountDest, "sentinel"))
							cmd.Stdout = stdOut
							Expect(cmd.Run()).To(Succeed())

							Expect(stdOut.String()).To(ContainSubstring("hello"))
						})
					})
				})
			})

			Context("when a file is supplied as a mount", func() {
				var (
					logFile   string
					mountFile string
				)

				BeforeEach(func() {
					l, err := ioutil.TempFile("", "winc.log")
					Expect(err).ToNot(HaveOccurred())
					Expect(l.Close()).To(Succeed())
					logFile = l.Name()

					m, err := ioutil.TempFile("", "mountfile")
					Expect(err).ToNot(HaveOccurred())
					Expect(m.Close()).To(Succeed())
					mountFile = m.Name()

					bundleSpec.Mounts = append(bundleSpec.Mounts, specs.Mount{
						Source:      mountFile,
						Destination: "C:\\foobar",
					})
				})

				AfterEach(func() {
					Expect(os.RemoveAll(logFile)).To(Succeed())
					Expect(os.RemoveAll(mountFile)).To(Succeed())
				})

				It("ignores it and logs that it did so", func() {
					err := exec.Command(wincBin, "--debug", "--log", logFile, "create", "-b", bundlePath, containerId).Run()
					Expect(err).ToNot(HaveOccurred())

					contents, err := ioutil.ReadFile(logFile)
					Expect(err).ToNot(HaveOccurred())
					Expect(string(contents)).To(ContainSubstring("mount is not a directory, ignoring"))
					Expect(string(contents)).To(ContainSubstring(fmt.Sprintf("mount=\"%s\"", mountFile)))
				})
			})
		})

		FContext("when the bundle config.json specifies a container memory limit", func() {
			var memLimitMB = uint64(128)

			BeforeEach(func() {
				memLimitBytes := memLimitMB * 1024 * 1024
				bundleSpec.Windows = &specs.Windows{
					Resources: &specs.WindowsResources{
						Memory: &specs.WindowsMemoryResources{
							Limit: &memLimitBytes,
						},
					},
				}
			})

			JustBeforeEach(func() {
				err := exec.Command(wincBin, "create", "-b", bundlePath, containerId).Run()
				Expect(err).ToNot(HaveOccurred())

				state, err := cm.State()
				Expect(err).ToNot(HaveOccurred())
				err = copy(consumeBin, filepath.Join("c:\\", "proc", strconv.Itoa(state.Pid), "root", "consume.exe"))
				Expect(err).NotTo(HaveOccurred())
			})

			grabMemory := func(mem int, exitCode int) string {
				cmd := exec.Command(wincBin, "exec", containerId, "c:\\consume.exe", strconv.Itoa(mem*1024*1024))
				session, err := gexec.Start(cmd, stdOut, stdErr)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, defaultTimeout*2).Should(gexec.Exit())
				if session.ExitCode() != exitCode {
					fmt.Println("STDOUT: " + stdOut.String())
					fmt.Println("STDERR: " + stdErr.String())
				}
				Expect(session.ExitCode()).To(Equal(exitCode))
				return stdErr.String()
			}

			It("is not constrained by smaller memory limit", func() {
				Expect(grabMemory(10, 0)).To(Equal(""))
			})

			It("is constrained by hitting the memory limit", func() {
				Expect(grabMemory(int(memLimitMB), 1)).To(ContainSubstring("Exception of type 'System.OutOfMemoryException' was thrown"))
			})
		})
	})

	Context("when the mount source does not exist", func() {
		BeforeEach(func() {
			mountDest := "C:\\mnt"
			mountSource := "C:\\not\\a\\directory\\mountsource"

			mount := specs.Mount{Destination: mountDest, Source: mountSource}
			bundleSpec.Mounts = []specs.Mount{mount}
		})

		It("errors and does not create the container", func() {
			cmd := exec.Command(wincBin, "create", "-b", bundlePath, containerId)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))

			Expect(containerExists(containerId)).To(BeFalse())
		})
	})

	Context("when provided a container id that already exists", func() {
		BeforeEach(func() {
			Expect(cm.Create(&bundleSpec)).To(Succeed())
		})

		AfterEach(func() {
			Expect(cm.Delete()).To(Succeed())
		})

		It("errors", func() {
			cmd := exec.Command(wincBin, "create", "-b", bundlePath, containerId)
			session, err := gexec.Start(cmd, stdOut, stdErr)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			expectedError := &hcsclient.AlreadyExistsError{Id: containerId}
			Expect(stdErr.String()).To(ContainSubstring(expectedError.Error()))
		})
	})

	Context("when the bundle directory name and container id do not match", func() {
		It("errors and does not create the container", func() {
			containerId = "doesnotmatchbundle"
			cmd := exec.Command(wincBin, "create", "-b", bundlePath, containerId)
			session, err := gexec.Start(cmd, stdOut, stdErr)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			expectedError := &hcsclient.InvalidIdError{Id: containerId}
			Expect(stdErr.String()).To(ContainSubstring(expectedError.Error()))

			Expect(containerExists(containerId)).To(BeFalse())
		})
	})

	Context("when using a custom rootfs", func() {
		var generatedTag string

		BeforeEach(func() {
			generatedTag = fmt.Sprintf("tag-%d", rand.Int())
			err := exec.Command("docker", "build", "-t", generatedTag, "-f", "fixtures\\Dockerfile.custom", "fixtures").Run()
			Expect(err).To(Succeed())

			dockerCmd := fmt.Sprintf("(docker inspect %s | ConvertFrom-Json).GraphDriver.Data.Dir", generatedTag)
			customRootfsPath, err := exec.Command("powershell.exe", dockerCmd).CombinedOutput()
			Expect(err).To(Succeed())

			bundleSpec = runtimeSpecGenerator(strings.TrimSpace(string(customRootfsPath)))
		})

		It("should find hello.txt in custom rootfs", func() {
			err := exec.Command(wincBin, "create", "-b", bundlePath, containerId).Run()
			Expect(err).ToNot(HaveOccurred())

			Expect(containerExists(containerId)).To(BeTrue())

			cmd := exec.Command(wincBin, "exec", containerId, "cmd.exe", "/C", "type C:\\hello.txt")
			cmd.Stdout = stdOut
			Expect(cmd.Run()).To(Succeed())
			Expect(stdOut.String()).To(ContainSubstring("hello from a text file"))
		})

		AfterEach(func() {
			Expect(cm.Delete()).To(Succeed())
			err := exec.Command("docker", "rmi", generatedTag).Run()
			Expect(err).To(Succeed())
		})
	})
})
