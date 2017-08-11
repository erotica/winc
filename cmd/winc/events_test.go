package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"

	"code.cloudfoundry.org/winc/container"
	"code.cloudfoundry.org/winc/hcs"
	"code.cloudfoundry.org/winc/volume"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = FDescribe("Events", func() {
	var (
		stdOut *bytes.Buffer
		stdErr *bytes.Buffer
	)

	BeforeEach(func() {
		stdOut = new(bytes.Buffer)
		stdErr = new(bytes.Buffer)
	})

	Context("given an existing container id", func() {
		var (
			containerId string
			cm          *container.Manager
			client      *hcs.Client
		)

		BeforeEach(func() {
			containerId = filepath.Base(bundlePath)

			client = &hcs.Client{}
			nm := networkManager(client)
			cm = container.NewManager(client, &volume.Mounter{}, nm, rootPath, bundlePath)

			bundleSpec := runtimeSpecGenerator(createSandbox(rootPath, rootfsPath, containerId), containerId)
			Expect(cm.Create(&bundleSpec)).To(Succeed())
		})

		AfterEach(func() {
			Expect(cm.Delete()).To(Succeed())
			Expect(execute(wincImageBin, "--store", rootPath, "delete", containerId)).To(Succeed())
		})

		Context("when the container has been created", func() {
			It("exits without error", func() {
				cmd := exec.Command(wincBin, "events", containerId)
				cmd.Stdout = stdOut
				Expect(cmd.Run()).To(Succeed())
			})

			Context("when passed the --stats flag", func() {
				type stats struct {
					Memory struct {
						Stats struct {
							TotalRss uint64 `json:"total_rss"`
						} `json:"raw"`
					} `json:"memory"`
				}

				It("prints the container stats to stdout", func() {
					// TODO: run the consume.exe binary to consume some memory

					cmd := exec.Command(wincBin, "events", "--stats", containerId)
					cmd.Stdout = stdOut
					Expect(cmd.Run()).To(Succeed())

					var s stats
					Expect(json.Unmarshal(stdOut.Bytes(), &s)).To(Succeed())
					fmt.Println(stdOut.String())
					fmt.Printf("****\n%+v\n", s)
					Expect(s.Memory.Stats.TotalRss).ToNot(BeEquivalentTo(0))

					//  TODO: expect that the memory is at least as much as consume.exe took ^
				})
			})
		})
	})

	Context("given a nonexistent container id", func() {
		It("errors", func() {
			cmd := exec.Command(wincBin, "events", "doesntexist")
			session, err := gexec.Start(cmd, stdOut, stdErr)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			expectedError := &hcs.NotFoundError{Id: "doesntexist"}
			Expect(stdErr.String()).To(ContainSubstring(expectedError.Error()))
		})
	})
})
