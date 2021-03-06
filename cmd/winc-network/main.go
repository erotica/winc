package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/winc/hcs"
	"code.cloudfoundry.org/winc/lib/filelock"
	"code.cloudfoundry.org/winc/lib/serial"
	"code.cloudfoundry.org/winc/netrules"
	"code.cloudfoundry.org/winc/netsh"
	"code.cloudfoundry.org/winc/network"
	"code.cloudfoundry.org/winc/port_allocator"
)

func main() {
	action, handle, configFile, err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid args: %s", err.Error())
		os.Exit(1)
	}

	config, err := parseConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "configFile: %s", err.Error())
		os.Exit(1)
	}

	nm := wireNetworkManager(config, handle)

	switch action {
	case "up":
		var inputs network.UpInputs
		if err := json.NewDecoder(os.Stdin).Decode(&inputs); err != nil {
			fmt.Fprintf(os.Stderr, "networkUp: %s", err.Error())
			os.Exit(1)
		}

		outputs, err := nm.Up(inputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "networkUp: %s", err.Error())
			os.Exit(1)
		}

		if err := json.NewEncoder(os.Stdout).Encode(outputs); err != nil {
			fmt.Fprintf(os.Stderr, "networkUp: %s", err.Error())
			os.Exit(1)
		}

	case "down":
		if err := nm.Down(); err != nil {
			fmt.Fprintf(os.Stderr, "networkDown: %s", err.Error())
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "invalid action: %s", action)
		os.Exit(1)
	}
}

func parseArgs(allArgs []string) (string, string, string, error) {
	var action, handle, configFile string
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)

	flagSet.StringVar(&action, "action", "", "")
	flagSet.StringVar(&handle, "handle", "", "")
	flagSet.StringVar(&configFile, "configFile", "", "")

	err := flagSet.Parse(allArgs[1:])
	if err != nil {
		return "", "", "", err
	}

	if handle == "" {
		return "", "", "", fmt.Errorf("missing required flag 'handle'")
	}

	if action == "" {
		return "", "", "", fmt.Errorf("missing required flag 'action'")
	}

	return action, handle, configFile, nil
}

func parseConfig(configFile string) (network.Config, error) {
	var config network.Config
	if configFile != "" {
		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			return config, err
		}

		if err := json.Unmarshal(content, &config); err != nil {
			return config, err
		}
	}

	return config, nil
}

func wireNetworkManager(config network.Config, handle string) *network.Manager {
	client := &hcs.Client{}
	runner := netsh.NewRunner(client, handle)
	applier := netrules.NewApplier(runner, handle)

	tracker := &port_allocator.Tracker{
		StartPort: 40000,
		Capacity:  5000,
	}

	locker := filelock.NewLocker("C:\\var\\vcap\\data\\winc\\port-state.json")

	pa := &port_allocator.PortAllocator{
		Tracker:    tracker,
		Serializer: &serial.Serial{},
		Locker:     locker,
	}

	return network.NewManager(
		client,
		pa,
		applier,
		config,
		handle,
	)
}
