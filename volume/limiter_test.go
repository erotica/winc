package volume_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"code.cloudfoundry.org/winc/volume"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Limiter", func() {
	const (
		volumeSize  = 20 * 1024 * 1024
		volumeLimit = 5 * 1024 * 1024
	)

	var (
		mountPath  string
		volumeDir  string
		volumeGuid string
	)

	BeforeEach(func() {
		var err error
		volumeDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		mountPath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		vhdxPath := filepath.Join(volumeDir, "vol.vhdx")
		volumeGuid = createVolume(vhdxPath, volumeSize)

		Expect(exec.Command("mountvol", mountPath, volumeGuid).Run()).To(Succeed())
	})

	AfterEach(func() {
		destroyVolume(filepath.Join(volumeDir, "vol.vhdx"))
		Expect(exec.Command("mountvol", mountPath, "/D").Run()).To(Succeed())
		Expect(os.RemoveAll(mountPath)).To(Succeed())
		Expect(os.RemoveAll(volumeDir)).To(Succeed())
	})

	It("applies the limit to the volume", func() {
		limiter := &volume.Limiter{}
		Expect(limiter.SetDiskLimit(volumeGuid, volumeLimit)).To(Succeed())

		largeFilePath := filepath.Join(mountPath, "file.txt")
		Expect(exec.Command("fsutil", "file", "createnew", largeFilePath, strconv.Itoa(volumeLimit+1)).Run()).ToNot(Succeed())
		Expect(largeFilePath).ToNot(BeAnExistingFile())
	})
})

func powershell(cmd ...string) string {
	args := append([]string{"-Command"}, cmd...)
	output, err := exec.Command("powershell.exe", args...).CombinedOutput()
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "Powershell command output: "+string(output))
	return string(output)
}

func createVolume(vhdxPath string, volumeSize int) string {
	powershell(fmt.Sprintf("New-VHD -Path %s -SizeBytes %d", vhdxPath, volumeSize))
	powershell("Mount-VHD " + vhdxPath)
	powershell(fmt.Sprintf("Get-VHD %s | Initialize-Disk -PartitionStyle MBR", vhdxPath))
	powershell(fmt.Sprintf("Get-DiskImage %s | Get-Disk | New-Partition -UseMaximumSize", vhdxPath))
	powershell(fmt.Sprintf("Get-DiskImage %s | Get-Disk | Get-Partition | Get-Volume | Format-Volume", vhdxPath))
	volumeGuid := powershell(fmt.Sprintf("(Get-DiskImage %s | Get-Disk | Get-Partition | Get-Volume).Path", vhdxPath))

	return strings.TrimSpace(volumeGuid)
}

func destroyVolume(vhdxPath string) {
	powershell("Dismount-VHD " + vhdxPath)
	Expect(os.RemoveAll(vhdxPath)).To(Succeed())
}
