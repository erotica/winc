package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"code.cloudfoundry.org/winc/container"
	"code.cloudfoundry.org/winc/hcsclient"
	"code.cloudfoundry.org/winc/mounter"
	"code.cloudfoundry.org/winc/sandbox"
	"github.com/Microsoft/hcsshim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Delete", func() {
	var (
		stdOut *bytes.Buffer
		stdErr *bytes.Buffer
	)

	BeforeEach(func() {
		stdOut = new(bytes.Buffer)
		stdErr = new(bytes.Buffer)
	})

	Context("when provided an existing container id", func() {
		var (
			containerId string
			cm          container.ContainerManager
		)

		BeforeEach(func() {
			containerId = filepath.Base(bundlePath)

			client := hcsclient.HCSClient{}
			sm := sandbox.NewManager(&client, &mounter.Mounter{}, "", bundlePath)
			nm := networkManager(&client)
			cm = container.NewManager(&client, sm, nm, containerId)

			bundleSpec := runtimeSpecGenerator(rootfsPath)
			Expect(cm.Create(&bundleSpec)).To(Succeed())
		})

		Context("when the container is running", func() {
			It("deletes the container", func() {
				err := exec.Command(wincBin, "delete", containerId).Run()
				Expect(err).ToNot(HaveOccurred())

				Expect(containerExists(containerId)).To(BeFalse())
			})

			It("deletes the container endpoints", func() {
				containerEndpoints := allEndpoints(containerId)

				err := exec.Command(wincBin, "delete", containerId).Run()
				Expect(err).ToNot(HaveOccurred())

				existingEndpoints, err := hcsshim.HNSListEndpointRequest()
				Expect(err).NotTo(HaveOccurred())

				for _, ep := range containerEndpoints {
					for _, existing := range existingEndpoints {
						Expect(ep).NotTo(Equal(existing.Id))
					}
				}
			})

			It("does not delete the bundle directory", func() {
				err := exec.Command(wincBin, "delete", containerId).Run()
				Expect(err).ToNot(HaveOccurred())

				Expect(bundlePath).To(BeADirectory())
			})

			It("unmounts sandbox.vhdx", func() {
				state, err := cm.State()
				Expect(err).NotTo(HaveOccurred())
				rootPath := filepath.Join("c:\\", "proc", strconv.Itoa(state.Pid), "root")
				_, err = os.Lstat(rootPath)
				Expect(err).NotTo(HaveOccurred())

				err = exec.Command(wincBin, "delete", containerId).Run()
				Expect(err).ToNot(HaveOccurred())

				Expect(rootPath).NotTo(BeADirectory())

				// if not cleanly unmounted, the mount point is left as a symlink
				_, err = os.Lstat(rootPath)
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Context("when provided a nonexistent container id", func() {
		It("errors", func() {
			cmd := exec.Command(wincBin, "delete", "nonexistentcontainer")
			session, err := gexec.Start(cmd, stdOut, stdErr)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			expectedError := &hcsclient.NotFoundError{Id: "nonexistentcontainer"}
			Expect(stdErr.String()).To(ContainSubstring(expectedError.Error()))
		})
	})
})
