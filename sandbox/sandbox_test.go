package sandbox_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"code.cloudfoundry.org/winc/hcsclient"
	"code.cloudfoundry.org/winc/hcsclient/hcsclientfakes"
	"code.cloudfoundry.org/winc/sandbox"
	"code.cloudfoundry.org/winc/sandbox/sandboxfakes"
	"github.com/Microsoft/hcsshim"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sandbox", func() {
	var (
		containerId        string
		rootfs             string
		depotDir           string
		hcsClient          *hcsclientfakes.FakeClient
		sandboxManager     sandbox.SandboxManager
		expectedDriverInfo hcsshim.DriverInfo
		rootfsParents      []byte
		fakeMounter        *sandboxfakes.FakeMounter
		volumePath         = "some-volume-path"
	)

	BeforeEach(func() {
		var err error
		rootfs, err = ioutil.TempDir("", "rootfs")
		Expect(err).ToNot(HaveOccurred())

		depotDir, err = ioutil.TempDir("", "depotDir")
		Expect(err).ToNot(HaveOccurred())

		rand.Seed(time.Now().UnixNano())
		containerId = strconv.Itoa(rand.Int())

		hcsClient = &hcsclientfakes.FakeClient{}
		fakeMounter = &sandboxfakes.FakeMounter{}
		sandboxManager = sandbox.NewManager(hcsClient, fakeMounter, depotDir, containerId)

		expectedDriverInfo = hcsshim.DriverInfo{
			HomeDir: depotDir,
			Flavour: 1,
		}
		rootfsParents = []byte(`["path1", "path2"]`)

		hcsClient.GetLayerMountPathReturns(volumePath, nil)
	})

	JustBeforeEach(func() {
		Expect(ioutil.WriteFile(filepath.Join(rootfs, "layerchain.json"), rootfsParents, 0755)).To(Succeed())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(rootfs)).To(Succeed())
		Expect(os.RemoveAll(depotDir)).To(Succeed())
	})

	Context("Create", func() {
		Context("when provided a rootfs layer and handle", func() {
			It("creates and activates the sandbox", func() {
				imageSpec, err := sandboxManager.Create(rootfs)
				Expect(err).ToNot(HaveOccurred())
				Expect(imageSpec.RootFs).To(Equal(volumePath))
				layers := imageSpec.Image.Config.Layers

				sandboxLayer := filepath.Join(expectedDriverInfo.HomeDir, containerId)
				expectedParentLayers := []string{rootfs, "path1", "path2"}

				Expect(layers).To(Equal(append([]string{sandboxLayer}, expectedParentLayers...)))

				Expect(hcsClient.CreateSandboxLayerCallCount()).To(Equal(1))
				driverInfo, actualContainerId, parentLayer, parentLayers := hcsClient.CreateSandboxLayerArgsForCall(0)
				Expect(driverInfo).To(Equal(expectedDriverInfo))
				Expect(actualContainerId).To(Equal(containerId))
				Expect(parentLayer).To(Equal(rootfs))
				Expect(parentLayers).To(Equal(expectedParentLayers))

				Expect(hcsClient.ActivateLayerCallCount()).To(Equal(1))
				driverInfo, actualContainerId = hcsClient.ActivateLayerArgsForCall(0)
				Expect(driverInfo).To(Equal(expectedDriverInfo))
				Expect(actualContainerId).To(Equal(containerId))

				Expect(hcsClient.PrepareLayerCallCount()).To(Equal(1))
				driverInfo, actualContainerId, parentLayers = hcsClient.PrepareLayerArgsForCall(0)
				Expect(driverInfo).To(Equal(expectedDriverInfo))
				Expect(actualContainerId).To(Equal(containerId))
				Expect(parentLayers).To(Equal(expectedParentLayers))
			})

			Context("when creating the sandbox layer fails", func() {
				var createSandboxLayerError = errors.New("create sandbox failed")

				BeforeEach(func() {
					hcsClient.CreateSandboxLayerReturns(createSandboxLayerError)
				})

				It("errors", func() {
					_, err := sandboxManager.Create(rootfs)
					Expect(err).To(Equal(createSandboxLayerError))
				})
			})

			Context("when activating the sandbox layer fails", func() {
				var activateLayerError = errors.New("activate sandbox failed")

				BeforeEach(func() {
					hcsClient.ActivateLayerReturns(activateLayerError)
				})

				It("errors", func() {
					_, err := sandboxManager.Create(rootfs)
					Expect(err).To(Equal(activateLayerError))
				})
			})

			Context("when preparing the sandbox layer fails", func() {
				var prepareLayerError = errors.New("prepare sandbox failed")

				BeforeEach(func() {
					hcsClient.PrepareLayerReturns(prepareLayerError)
				})

				It("errors", func() {
					_, err := sandboxManager.Create(rootfs)
					Expect(err).To(Equal(prepareLayerError))
				})
			})
		})

		Context("when provided a nonexistent rootfs layer", func() {
			It("errors", func() {
				_, err := sandboxManager.Create("nonexistentrootfs")
				pathErr, isPathError := err.(*os.PathError)
				Expect(isPathError).To(BeTrue())
				Expect(pathErr.Path).To(Equal(filepath.Join("nonexistentrootfs", "layerchain.json")))
			})
		})

		Context("when provided a rootfs layer missing a layerchain.json", func() {
			JustBeforeEach(func() {
				Expect(os.RemoveAll(filepath.Join(rootfs, "layerchain.json"))).To(Succeed())
			})

			It("errors", func() {
				_, err := sandboxManager.Create(rootfs)
				pathErr, isPathError := err.(*os.PathError)
				Expect(isPathError).To(BeTrue())
				Expect(pathErr.Path).To(Equal(filepath.Join(rootfs, "layerchain.json")))
			})
		})

		Context("when the rootfs has a layerchain.json that is invalid JSON", func() {
			BeforeEach(func() {
				rootfsParents = []byte("[")
			})

			It("errors", func() {
				_, err := sandboxManager.Create(rootfs)
				Expect(err).To(Equal(&sandbox.InvalidRootfsLayerChainError{Path: rootfs}))
			})
		})

		Context("when getting the volume mount path of the container fails", func() {
			Context("when getting the volume returned an error", func() {
				var layerMountPathError = errors.New("could not get volume")

				BeforeEach(func() {
					hcsClient.GetLayerMountPathReturns("", layerMountPathError)
				})

				It("errors", func() {
					_, err := sandboxManager.Create(rootfs)
					Expect(err).To(Equal(layerMountPathError))
				})
			})

			Context("when the volume returned is empty", func() {
				BeforeEach(func() {
					hcsClient.GetLayerMountPathReturns("", nil)
				})

				It("errors", func() {
					_, err := sandboxManager.Create(rootfs)
					Expect(err).To(Equal(&hcsclient.MissingVolumePathError{Id: containerId}))
				})
			})
		})
	})

	XContext("Delete", func() {
		var bundlePath string
		It("unprepares and deactivates the bundlePath", func() {
			err := sandboxManager.Delete()
			Expect(err).ToNot(HaveOccurred())

			Expect(hcsClient.UnprepareLayerCallCount()).To(Equal(1))
			driverInfo, layerId := hcsClient.UnprepareLayerArgsForCall(0)
			Expect(driverInfo).To(Equal(expectedDriverInfo))
			Expect(layerId).To(Equal(containerId))

			Expect(hcsClient.DeactivateLayerCallCount()).To(Equal(1))
			driverInfo, layerId = hcsClient.DeactivateLayerArgsForCall(0)
			Expect(driverInfo).To(Equal(expectedDriverInfo))
			Expect(layerId).To(Equal(containerId))
		})

		It("only deletes the files that the container created", func() {
			sentinelPath := filepath.Join(bundlePath, "sentinel")
			f, err := os.Create(sentinelPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(f.Close()).To(Succeed())

			err = sandboxManager.Delete()
			Expect(err).ToNot(HaveOccurred())

			files, err := filepath.Glob(filepath.Join(bundlePath, "*"))
			Expect(err).ToNot(HaveOccurred())
			Expect(files).To(ConsistOf([]string{filepath.Join(bundlePath, "sentinel")}))
		})

		Context("when unpreparing the bundlePath fails", func() {
			var unprepareLayerError = errors.New("unprepare sandbox failed")

			BeforeEach(func() {
				hcsClient.UnprepareLayerReturns(unprepareLayerError)
			})

			It("errors", func() {
				err := sandboxManager.Delete()
				Expect(err).To(Equal(unprepareLayerError))
			})
		})

		Context("when deactivating the bundlePath fails", func() {
			var deactivateLayerError = errors.New("deactivate sandbox failed")

			BeforeEach(func() {
				hcsClient.DeactivateLayerReturns(deactivateLayerError)
			})

			It("errors", func() {
				err := sandboxManager.Delete()
				Expect(err).To(Equal(deactivateLayerError))
			})
		})
	})

	Context("Mount", func() {
		It("mounts the sandbox.vhdx at C:\\proc\\{{pid}}\\root", func() {
			pid := rand.Int()
			Expect(sandboxManager.Mount(pid, volumePath)).To(Succeed())

			rootPath := filepath.Join("c:\\", "proc", fmt.Sprintf("%d", pid), "root")
			Expect(rootPath).To(BeADirectory())

			Expect(fakeMounter.SetPointCallCount()).To(Equal(1))
			mp, vol := fakeMounter.SetPointArgsForCall(0)
			Expect(mp).To(Equal(rootPath))
			Expect(vol).To(Equal(volumePath))
		})
	})

	Context("Unmount", func() {
		var pid int
		var mountPath string
		var rootPath string

		BeforeEach(func() {
			pid = rand.Int()
			mountPath = filepath.Join("c:\\", "proc", fmt.Sprintf("%d", pid))
			rootPath = filepath.Join(mountPath, "root")
			Expect(os.MkdirAll(rootPath, 0755)).To(Succeed())
		})

		It("unmounts the sandbox.vhdx from c:\\proc\\<pid>\\mnt and removes the directory", func() {
			Expect(sandboxManager.Unmount(pid)).To(Succeed())

			Expect(fakeMounter.DeletePointCallCount()).To(Equal(1))
			mp := fakeMounter.DeletePointArgsForCall(0)
			Expect(mp).To(Equal(rootPath))

			Expect(mountPath).NotTo(BeADirectory())
		})
	})
})
