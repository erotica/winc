package main

import (
	"fmt"

	"github.com/Microsoft/hcsshim"
)

func main() {
	eps, err := hcsshim.HNSListEndpointRequest()
	if err != nil {
		panic(err)
	}

	for _, ep := range eps {
		fmt.Printf("%s\n", ep.Id)
		fmt.Printf("\t%s\n", ep.Name)
		fmt.Printf("\t%s\n", ep.IPAddress.String())
		fmt.Printf("\t%s\n", ep.VirtualNetwork)
		fmt.Printf("\t%s\n", ep.VirtualNetworkName)
		for _, pol := range ep.Policies {
			fmt.Printf("\t%s\n", string(pol))
		}
		fmt.Printf("\n")

		// if ep.Name == "c2" {
		// 	ep.Delete()
		// }
	}
}
