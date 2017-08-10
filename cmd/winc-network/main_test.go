package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"code.cloudfoundry.org/winc/network"

	"github.com/Microsoft/hcsshim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

var _ = Describe("up", func() {
	var (
		config      []byte
		containerId string
		bundleSpec  specs.Spec
		err         error
		stdOut      *bytes.Buffer
		stdErr      *bytes.Buffer
	)

	BeforeEach(func() {
		containerId = filepath.Base(bundlePath)
		bundleSpec = runtimeSpecGenerator(createSandbox("C:\\run\\winc", rootfsPath, containerId), containerId)
		config, err = json.Marshal(&bundleSpec)
		Expect(err).NotTo(HaveOccurred())
		Expect(ioutil.WriteFile(filepath.Join(bundlePath, "config.json"), config, 0666)).To(Succeed())

		err := exec.Command(wincBin, "create", "-b", bundlePath, containerId).Run()
		Expect(err).ToNot(HaveOccurred())

		stdOut = new(bytes.Buffer)
		stdErr = new(bytes.Buffer)
	})

	AfterEach(func() {
		err := exec.Command(wincBin, "delete", containerId).Run()
		Expect(err).ToNot(HaveOccurred())
		Expect(exec.Command(wincImageBin, "--store", "C:\\run\\winc", "delete", containerId).Run()).To(Succeed())
	})

	Context("stdin contains a port mapping request", func() {
		It("prints the correct port mapping for the container", func() {
			cmd := exec.Command(wincNetworkBin, "--action", "up", "--handle", containerId)
			cmd.Stdin = strings.NewReader(`{"Pid": 123, "Properties": {} ,"netin": [{"host_port": 0, "container_port": 8080}]}`)
			output, err := cmd.CombinedOutput()
			Expect(err).To(Succeed())

			regex := `{"properties":{"garden\.network\.container-ip":"\d+\.\d+\.\d+\.\d+","garden\.network\.host-ip":"255\.255\.255\.255","garden\.network\.mapped-ports":"\[{\\"HostPort\\":\d+,\\"ContainerPort\\":8080}\]"}}`
			Expect(string(output)).To(MatchRegexp(regex))
		})

		It("outputs the host's public IP as the container IP", func() {
			cmd := exec.Command(wincNetworkBin, "--action", "up", "--handle", containerId)
			cmd.Stdin = strings.NewReader(`{"Pid": 123, "Properties": {} ,"netin": [{"host_port": 0, "container_port": 8080}]}`)
			output, err := cmd.CombinedOutput()
			Expect(err).To(Succeed())

			regex := regexp.MustCompile(`"garden\.network\.container-ip":"(\d+\.\d+\.\d+\.\d+)"`)
			matches := regex.FindStringSubmatch(string(output))
			Expect(len(matches)).To(Equal(2))

			cmd = exec.Command("powershell", "-Command", "Get-NetIPAddress", matches[1])
			output, err = cmd.CombinedOutput()
			Expect(err).To(BeNil())
			Expect(string(output)).NotTo(ContainSubstring("Loopback"))
			Expect(string(output)).NotTo(ContainSubstring("HNS Internal NIC"))
			Expect(string(output)).To(MatchRegexp("AddressFamily.*IPv4"))
		})
	})

	Context("stdin contains a port mapping request with two ports", func() {
		It("prints the correct port mapping for the container", func() {
			cmd := exec.Command(wincNetworkBin, "--action", "up", "--handle", containerId)
			cmd.Stdin = strings.NewReader(`{"Pid": 123, "Properties": {} ,"netin": [{"host_port": 0, "container_port": 8080}, {"host_port": 0, "container_port": 2222}]}`)
			output, err := cmd.CombinedOutput()
			Expect(err).To(Succeed())

			regex := `{"properties":{"garden\.network\.container-ip":"\d+\.\d+\.\d+\.\d+","garden\.network\.host-ip":"255\.255\.255\.255","garden\.network\.mapped-ports":"\[{\\"HostPort\\":\d+,\\"ContainerPort\\":8080},{\\"HostPort\\":\d+,\\"ContainerPort\\":2222}\]"}}`
			Expect(string(output)).To(MatchRegexp(regex))
		})
	})

	Context("stdin does not contain a port mapping request", func() {
		It("prints an empty list of mapped ports", func() {
			cmd := exec.Command(wincNetworkBin, "--action", "up", "--handle", containerId)
			cmd.Stdin = strings.NewReader(`{"Pid": 123, "Properties": {} }`)
			output, err := cmd.CombinedOutput()
			Expect(err).To(Succeed())

			regex := `{"properties":{"garden\.network\.container-ip":"\d+\.\d+\.\d+\.\d+","garden\.network\.host-ip":"255\.255\.255\.255","garden\.network\.mapped-ports":"\[\]"}}`
			Expect(string(output)).To(MatchRegexp(regex))
		})
	})

	Context("stdin contains an invalid port mapping request", func() {
		It("errors", func() {
			cmd := exec.Command(wincNetworkBin, "--action", "up", "--handle", containerId)
			cmd.Stdin = strings.NewReader(`{"Pid": 123, "Properties": {} ,"netin": [{"host_port": 0, "container_port": 1234}, {"host_port": 0, "container_port": 2222}]}`)
			session, err := gexec.Start(cmd, stdOut, stdErr)
			Expect(err).To(Succeed())

			Eventually(session).Should(gexec.Exit(1))
			Expect(stdErr.String()).To(ContainSubstring("invalid port mapping"))
		})
	})

	Context("stdin contains a net out rule request with single ip/port", func() {
		type firewall struct {
			Protocol   string `json:"Protocol"`
			RemotePort string `json:"RemotePort"`
		}

		var containerIp string

		const getContainerFirewall = `Get-NetFirewallAddressFilter | ?{$_.RemoteAddress -eq "8.8.8.8" -and $_.LocalAddress -eq "%s"} | Get-NetFirewallRule | Get-NetFirewallPortFilter | ConvertTo-Json`

		AfterEach(func() {
			Expect(exec.Command(wincNetworkBin, "--action", "down", "--handle", containerId).Run()).To(Succeed())
			parsedCmd := fmt.Sprintf(getContainerFirewall, containerIp)
			output, err := exec.Command("powershell.exe", "-Command", parsedCmd).CombinedOutput()
			Expect(err).To(Succeed())
			Expect(string(output)).To(BeEmpty())
		})

		It("creates the correct firewall rule", func() {
			cmd := exec.Command(wincNetworkBin, "--action", "up", "--handle", containerId)
			netOutRule := network.NetOutRule{
				Protocol: network.ProtocolTCP,
				Networks: []network.IPRange{network.IPRangeFromIP(net.ParseIP("8.8.8.8"))},
				Ports:    []network.PortRange{network.PortRangeFromPort(80)},
			}
			netOutRuleStr, err := json.Marshal(&netOutRule)
			Expect(err).ToNot(HaveOccurred())

			cmd.Stdin = strings.NewReader(fmt.Sprintf(`{"Pid": 123, "Properties": {}, "netout_rules": [%s]}`, string(netOutRuleStr)))
			Expect(cmd.Run()).To(Succeed())

			containerIp = getContainerIp(containerId).String()
			parsedCmd := fmt.Sprintf(getContainerFirewall, containerIp)
			output, err := exec.Command("powershell.exe", "-Command", parsedCmd).CombinedOutput()
			Expect(err).To(Succeed())

			var f firewall
			Expect(json.Unmarshal(output, &f)).To(Succeed())

			Expect(f.Protocol).To(Equal("TCP"))
			Expect(f.RemotePort).To(Equal("80"))
		})
	})

})

func getContainerIp(containerId string) net.IP {
	container, err := hcsshim.OpenContainer(containerId)
	Expect(err).ToNot(HaveOccurred(), "no containers with id: "+containerId)

	stats, err := container.Statistics()
	Expect(err).ToNot(HaveOccurred())

	Expect(stats.Network).ToNot(BeEmpty(), "container has no networks attached: "+containerId)
	endpoint, err := hcsshim.GetHNSEndpointByID(stats.Network[0].EndpointId)
	Expect(err).ToNot(HaveOccurred())

	return endpoint.IPAddress
}
