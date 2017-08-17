package layer_test

import (
	"errors"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/winc/layer"
	"code.cloudfoundry.org/winc/layer/layerfakes"

	"github.com/Microsoft/hcsshim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager", func() {
	var (
		client       *layerfakes.FakeHCSClient
		m            *layer.Manager
		parentLayers []string
		storeDir     string
		err          error
		containerId  string
	)

	const expectedVolumeGuid = `\\?\Volume{some-guid}\`

	BeforeEach(func() {
		storeDir, err = ioutil.TempDir("", "layer-store-dir")
		Expect(err).NotTo(HaveOccurred())

		parentLayers = []string{"layer-2", "layer-1"}

		client = &layerfakes.FakeHCSClient{}
		m = layer.NewManager(client, storeDir)
	})

	AfterEach(func() {
		Expect(os.RemoveAll(storeDir)).To(Succeed())
	})

	Describe("CreateLayer", func() {
		BeforeEach(func() {
			client.GetLayerMountPathReturns(expectedVolumeGuid, nil)
		})

		It("creates the layer", func() {
			volumeGuid, err := client.CreateLayer(driverInfo, containerId, rootfsPath, sandboxLayers)
			Expect(err).ToNot(HaveOccurred())

			Expect

			Expect(volumeGuid).To(Equal(expectedVolumeGuid))
		})

		Context("when the layer has been created but not activated", func() {
			BeforeEach(func() {
				Expect(hcsshim.CreateSandboxLayer(driverInfo, containerId, rootfsPath, sandboxLayers)).To(Succeed())
				_, _ = client.CreateLayer(driverInfo, containerId, rootfsPath, sandboxLayers)
			})

			It("continues and creates the layer", func() {
				volumeGuid, err := client.CreateLayer(driverInfo, containerId, rootfsPath, sandboxLayers)
				Expect(err).ToNot(HaveOccurred())

				expectedVolumeGuid, err := hcsshim.GetLayerMountPath(driverInfo, containerId)
				Expect(err).ToNot(HaveOccurred())
				Expect(volumeGuid).ToNot(BeEmpty())
				Expect(volumeGuid).To(Equal(expectedVolumeGuid))
			})
		})

		Context("when the layer has been created and activated but not prepared", func() {
			BeforeEach(func() {
				Expect(hcsshim.CreateSandboxLayer(driverInfo, containerId, rootfsPath, sandboxLayers)).To(Succeed())
				Expect(hcsshim.ActivateLayer(driverInfo, containerId)).To(Succeed())
			})

			It("continues and creates the layer", func() {
				volumeGuid, err := client.CreateLayer(driverInfo, containerId, rootfsPath, sandboxLayers)
				Expect(err).ToNot(HaveOccurred())

				expectedVolumeGuid, err := hcsshim.GetLayerMountPath(driverInfo, containerId)
				Expect(err).ToNot(HaveOccurred())
				Expect(volumeGuid).ToNot(BeEmpty())
				Expect(volumeGuid).To(Equal(expectedVolumeGuid))
			})
		})

		Context("when the layer has been created, activated, and prepared", func() {
			BeforeEach(func() {
				Expect(hcsshim.CreateSandboxLayer(driverInfo, containerId, rootfsPath, sandboxLayers)).To(Succeed())
				Expect(hcsshim.ActivateLayer(driverInfo, containerId)).To(Succeed())
				Expect(hcsshim.PrepareLayer(driverInfo, containerId, sandboxLayers)).To(Succeed())
			})

			It("continues and creates the layer", func() {
				volumeGuid, err := client.CreateLayer(driverInfo, containerId, rootfsPath, sandboxLayers)
				Expect(err).ToNot(HaveOccurred())

				expectedVolumeGuid, err := hcsshim.GetLayerMountPath(driverInfo, containerId)
				Expect(err).ToNot(HaveOccurred())
				Expect(volumeGuid).ToNot(BeEmpty())
				Expect(volumeGuid).To(Equal(expectedVolumeGuid))
			})
		})
	})

	Describe("DestroyLayer", func() {
		Context("when the layer exists", func() {
			BeforeEach(func() {
				Expect(hcsshim.CreateSandboxLayer(driverInfo, containerId, rootfsPath, sandboxLayers)).To(Succeed())
				Expect(hcsshim.ActivateLayer(driverInfo, containerId)).To(Succeed())
				Expect(hcsshim.PrepareLayer(driverInfo, containerId, sandboxLayers)).To(Succeed())
			})

			FIt("destroys the layer", func() {
				Expect(client.DestroyLayer(driverInfo, containerId)).To(Succeed())
				Expect(hcsshim.LayerExists(driverInfo, containerId)).To(BeFalse())
			})
		})

		Context("when the layer exists but is not prepared", func() {
			BeforeEach(func() {
				Expect(hcsshim.CreateSandboxLayer(driverInfo, containerId, rootfsPath, sandboxLayers)).To(Succeed())
				Expect(hcsshim.ActivateLayer(driverInfo, containerId)).To(Succeed())
			})

			It("destroys the layer", func() {
				Expect(client.DestroyLayer(driverInfo, containerId)).To(Succeed())
				Expect(hcsshim.LayerExists(driverInfo, containerId)).To(BeFalse())
			})
		})

		Context("when the layer exists but is not activated", func() {
			BeforeEach(func() {
				Expect(hcsshim.CreateSandboxLayer(driverInfo, containerId, rootfsPath, sandboxLayers)).To(Succeed())
			})

			It("destroys the layer", func() {
				Expect(client.DestroyLayer(driverInfo, containerId)).To(Succeed())
				Expect(hcsshim.LayerExists(driverInfo, containerId)).To(BeFalse())
			})
		})

		Context("when the layer does not exist", func() {
			It("succeeds", func() {
				Expect(client.DestroyLayer(driverInfo, containerId)).To(Succeed())
			})
		})
	})

	Describe("Retryable", func() {
		Context("when the error is a timeout error", func() {
			It("returns true", func() {
				err := errors.New("Some operation failed: This operation returned because the timeout period expired")
				Expect(client.Retryable(err)).To(BeTrue())
			})
		})

		Context("when the error is something else", func() {
			It("returns false", func() {
				err := errors.New("some other error")
				Expect(client.Retryable(err)).To(BeFalse())
			})
		})
	})
})
