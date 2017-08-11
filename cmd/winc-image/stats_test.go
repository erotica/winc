package main_test

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Stats", func() {
	var (
		storePath   string
		containerId string
	)

	BeforeEach(func() {
		var err error
		rand.Seed(time.Now().UnixNano())
		containerId = strconv.Itoa(rand.Int())
		storePath, err = ioutil.TempDir("", "container-store")
		Expect(err).ToNot(HaveOccurred())

		_, _, err = execute(wincImageBin, "--store", storePath, "create", rootfsPath, containerId)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(exec.Command(wincImageBin, "--store", storePath, "delete", containerId).Run()).To(Succeed())
		Expect(os.RemoveAll(storePath)).To(Succeed())
	})

	type DiskUsage struct {
		TotalBytesUsed     uint64 `json:"total_bytes_used"`
		ExclusiveBytesUsed uint64 `json:"exclusive_bytes_used"`
	}

	type ImageStats struct {
		Disk DiskUsage `json:"disk_usage"`
	}

	It("reports the image stats", func() {
		stdout, _, err := execute(wincImageBin, "--store", storePath, "stats", containerId)
		Expect(err).NotTo(HaveOccurred())
		var imageStats ImageStats
		Expect(json.Unmarshal(stdout.Bytes(), &imageStats)).To(Succeed())
		Expect(imageStats.Disk.TotalBytesUsed).To(Equal(uint64(1111)))
		Expect(imageStats.Disk.ExclusiveBytesUsed).To(Equal(uint64(2222)))
	})
})
