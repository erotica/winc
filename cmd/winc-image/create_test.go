package main_test

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Microsoft/hcsshim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	Context("when provided a rootfs and handle arguments", func() {
		var (
			stdOut *bytes.Buffer
			// stdErr      *bytes.Buffer
			containerId string
		)

		BeforeEach(func() {
			stdOut = new(bytes.Buffer)
			// stdErr = new(bytes.Buffer)
			rand.Seed(time.Now().UnixNano())
			containerId = strconv.Itoa(rand.Int())
		})

		type DesiredImageSpec struct {
			RootFS string `json:"rootfs,omitempty"`
			Image  struct {
				Config struct {
					Layers []string `json:"layers,omitempty"`
				} `json:"config,omitempty"`
			} `json:"image,omitempty"`
		}

		It("creates a sandbox layer and outputs the volume guid and layers on stdout", func() {
			cmd := exec.Command(wincImageBin, "create", rootfsPath, containerId)
			cmd.Stdout = stdOut
			err := cmd.Run()
			Expect(err).ToNot(HaveOccurred())

			var desiredImageSpec DesiredImageSpec
			Expect(json.Unmarshal(stdOut.Bytes(), &desiredImageSpec)).To(Succeed())

			Expect(desiredImageSpec.Image.Config.Layers).ToNot(BeEmpty())
			for _, layer := range desiredImageSpec.Image.Config.Layers {
				Expect(layer).To(BeADirectory())
			}

			sandboxPath := desiredImageSpec.Image.Config.Layers[0]
			Expect(filepath.Join(sandboxPath, "initialized")).To(BeAnExistingFile())

			err = exec.Command("powershell", "-Command", "Test-VHD", filepath.Join(sandboxPath, "sandbox.vhdx")).Run()
			Expect(err).ToNot(HaveOccurred())

			Expect(desiredImageSpec.RootFS).To(Equal(getVolumeGuid(containerId)))
		})

		Context("when the rootfs argument is invalid", func() {
			XIt("errors", func() {
			})
		})
	})
})

func getVolumeGuid(id string) string {
	driverInfo := hcsshim.DriverInfo{
		HomeDir: os.TempDir(),
		Flavour: 1,
	}
	volumePath, err := hcsshim.GetLayerMountPath(driverInfo, id)
	Expect(err).NotTo(HaveOccurred())
	return volumePath
}
