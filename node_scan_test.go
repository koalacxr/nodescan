package nodescan

import (
	"fmt"
	"testing"
	"time"
)

func TestLocalNewNodeScanner(t *testing.T) {
	newNodeScanner := NewNodeScanner(
		SetIsLocal(true),
		SetTimeout(300*time.Millisecond),
		SetPorts([]int{3306, 8082}),
		SetNetworkTypes([]NetworkType{NetWorkTCP4}),
		SetMaxChannel(100))
	ipsa := newNodeScanner.Scan()
	for _, ipa := range ipsa.Values {
		fmt.Println(fmt.Sprintf("%s,%v:%d", ipa.Network, ipa.IP, ipa.Port))
	}
}

func TestRemoteNewNodeScanner(t *testing.T) {
	newNodeScanner := NewNodeScanner(
		SetIsLocal(false),
		SetIps([]string{"10.128.61.138"}),
		SetTimeout(300*time.Millisecond),
		SetPorts([]int{80, 443, 3306, 8082}),
		SetNetworkTypes([]NetworkType{NetWorkTCP4}),
		SetMaxChannel(100))
	ipsa := newNodeScanner.Scan()
	for _, ipa := range ipsa.Values {
		fmt.Println(fmt.Sprintf("%s,%v:%d", ipa.Network, ipa.IP, ipa.Port))
	}
}

func TestNodeScanner_LocalCIDRs(t *testing.T) {
	newNodeScanner := NewNodeScanner(SetIsLocal(true))
	LocalCIDRs := newNodeScanner.LocalCIDRs()
	for _, locCIDR := range LocalCIDRs {
		fmt.Println(fmt.Sprintf("%s", locCIDR))
	}
}

func TestNodeScanner_LocalIPs(t *testing.T) {
	newNodeScanner := NewNodeScanner(SetIsLocal(true))
	LocalIPs := newNodeScanner.LocalIPs()
	for _, locIP := range LocalIPs {
		fmt.Println(fmt.Sprintf("%s", locIP))
	}
}

func TestNodeScanner_LocalScan(t *testing.T) {
	newNodeScanner := NewNodeScanner(SetIsLocal(true), SetPorts([]int{80, 443, 3306, 8082}))
	LocalIPs := newNodeScanner.LocalScan()
	for _, locIP := range LocalIPs.Values {
		fmt.Println(fmt.Sprintf("%s:%v", locIP.IP, locIP.Port))
	}
}
