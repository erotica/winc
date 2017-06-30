package main_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

var _ = Describe("Spec", func() {
	Context("when the bundle path is not specified", func() {
		BeforeEach(func() {
			Expect(os.RemoveAll("config.json")).To(Succeed())
		})

		AfterEach(func() {
			Expect(os.RemoveAll("config.json")).To(Succeed())
		})

		It("creates a config.json in the current directory", func() {
			cmd := exec.Command(wincBin, "spec")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			config, err := ioutil.ReadFile("config.json")
			Expect(err).NotTo(HaveOccurred())
			var spec specs.Spec
			Expect(json.Unmarshal(config, &spec)).To(Succeed())
			Expect(spec.Version).To(Equal(specs.Version))
			Expect(spec.Platform.Arch).To(Equal(runtime.GOARCH))
			Expect(spec.Platform.OS).To(Equal(runtime.GOOS))
			Expect(spec.Process.Args).To(Equal([]string{"powershell"}))
			Expect(spec.Process.Cwd).To(Equal("/"))
			Expect(filepath.Join(spec.Root.Path, "layerchain.json")).To(BeAnExistingFile())
		})
	})

	Context("when the bundle path is specifided", func() {
		var bundlePath string

		BeforeEach(func() {
			var err error
			bundlePath, err = ioutil.TempDir("", "winc-spec-test")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(os.RemoveAll(bundlePath)).To(Succeed())
		})

		It("creates a config.json in the bundle path", func() {
			cmd := exec.Command(wincBin, "spec", "-b", bundlePath)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			config, err := ioutil.ReadFile(filepath.Join(bundlePath, "config.json"))
			Expect(err).NotTo(HaveOccurred())
			var spec specs.Spec
			Expect(json.Unmarshal(config, &spec)).To(Succeed())
		})
	})
})
