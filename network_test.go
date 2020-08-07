package nodescan

import (
	"fmt"
	"testing"
)

func TestGetNetworkIpAddresses(t *testing.T) {
	hosts , err := GetLocalNetworkIpAddresses(NetWorkTCP4, "")
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(hosts)
}

func TestGetSubNetIps(t *testing.T) {
	hosts, err := GetLocalSubNetIps("10.128.58.157/20")
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(hosts)
}

func TestGetNetworkSubNetIps(t *testing.T) {
	localIps, err := GetLocalNetworkIpAddresses(NetWorkTCP4, "")
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, network := range localIps {
		hosts, err := GetLocalSubNetIps(network)
		if err != nil {
			t.Fatal(err.Error())
		}

		fmt.Println(hosts)
	}
}