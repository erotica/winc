package main_test

import (
	"encoding/json"
	"os/exec"
	"path/filepath"

	winc "code.cloudfoundry.org/winc/cmd/winc"
	"code.cloudfoundry.org/winc/command"
	"code.cloudfoundry.org/winc/container"
	"code.cloudfoundry.org/winc/hcsclient"
	"code.cloudfoundry.org/winc/sandbox"
	ps "github.com/mitchellh/go-ps"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

var _ = Describe("State", func() {
	Context("given an existing container id", func() {
		var (
			containerId string
			cm          winc.ContainerManager
			actualState *specs.State
			client      hcsclient.Client
		)

		BeforeEach(func() {
			containerId = filepath.Base(bundlePath)

			client = &hcsclient.HCSClient{}
			sm := sandbox.NewManager(client, &command.Command{}, bundlePath)
			nm := networkManager(client)
			cm = container.NewManager(client, sm, nm, containerId)

			bundleSpec := runtimeSpecGenerator(rootfsPath)
			Expect(cm.Create(&bundleSpec)).To(Succeed())
		})

		AfterEach(func() {
			Expect(cm.Delete()).To(Succeed())
		})

		Context("when the container has been created", func() {
			It("prints the state of the container to stdout", func() {
				cmd := exec.Command(wincBin, "state", containerId)
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))

				actualState = &specs.State{}
				Expect(json.Unmarshal(session.Out.Contents(), actualState)).To(Succeed())

				Expect(actualState.Status).To(Equal("created"))
				Expect(actualState.Version).To(Equal(specs.Version))
				Expect(actualState.ID).To(Equal(containerId))
				Expect(actualState.Bundle).To(Equal(bundlePath))

				p, err := ps.FindProcess(actualState.Pid)
				Expect(err).ToNot(HaveOccurred())
				Expect(p.Executable()).To(Equal("wininit.exe"))
			})
		})
	})

	Context("given a nonexistent container id", func() {
		It("errors", func() {
			cmd := exec.Command(wincBin, "state", "doesntexist")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			expectedError := &hcsclient.NotFoundError{Id: "doesntexist"}
			Eventually(session.Err).Should(gbytes.Say(expectedError.Error()))
		})
	})
})
