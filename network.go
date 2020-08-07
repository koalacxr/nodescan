package nodescan

import (
	"fmt"
	"math"
	"net"
	"strings"
)

type NetworkType string

const (
	NetWorkTCP        NetworkType = "tcp"
	NetWorkTCP4       NetworkType = "tcp4"
	NetWorkTCP6       NetworkType = "tcp6"
	NetWorkUDP        NetworkType = "udp"
	NetWorkUDP4       NetworkType = "udp4"
	NetWorkUDP6       NetworkType = "udp6"
	NetWorkIP         NetworkType = "ip"
	NetWorkIP4        NetworkType = "ip4"
	NetWorkIP6        NetworkType = "ip6"
	NetWorkUNIX       NetworkType = "unix"
	NetWorkUNIXGRAM   NetworkType = "unixgram"
	NetWorkUNIXPACKET NetworkType = "unixpacket"

	ConstLocalIP   = "127.0.0.1/8"
	ConstLocalhost = "localhost"
)

var networks map[NetworkType]struct{}

func init() {
	networks = make(map[NetworkType]struct{}, 0)
	networks[NetWorkTCP] = struct{}{}
	networks[NetWorkTCP4] = struct{}{}
	networks[NetWorkTCP6] = struct{}{}
	networks[NetWorkUDP] = struct{}{}
	networks[NetWorkUDP4] = struct{}{}
	networks[NetWorkUDP6] = struct{}{}
	networks[NetWorkIP] = struct{}{}
	networks[NetWorkIP4] = struct{}{}
	networks[NetWorkIP6] = struct{}{}
	networks[NetWorkUNIX] = struct{}{}
	networks[NetWorkUNIXGRAM] = struct{}{}
	networks[NetWorkUNIXPACKET] = struct{}{}
}

func CheckNetworkType(value NetworkType) bool {
	if _, ok := networks[value]; ok {
		return true
	}

	return false
}

func GetLocalNetworkIpAddresses(network NetworkType, ipSegment string) ([]string, error) {
	returnAddresses := make([]string, 0)

	if !CheckNetworkType(network) {
		return returnAddresses, fmt.Errorf("invalid network %s", network)
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return returnAddresses, err
	}

	for _, itemInterface := range interfaces {
		addresses, err := itemInterface.Addrs()
		if err != nil {
			continue
		}

		for _, itemAddress := range addresses {
			ip, _, err := net.ParseCIDR(itemAddress.String())
			if err != nil {
				continue
			}

			if strings.Contains(string(network), "4") && ip.To4() == nil {
				continue
			}

			if strings.Contains(string(network), "6") && ip.To16() == nil {
				continue
			}

			if IsLocalAddress(ip, ipSegment) {
				returnAddresses = append(returnAddresses, itemAddress.String())
			}
		}
	}

	return returnAddresses, nil
}

func GetLocalSubNetIps(localIp string) ([]string, error) {
	ips := make([]string, 0)

	if localIp == ConstLocalIP {
		return ips, nil
	}

	_, ipNet, err := net.ParseCIDR(localIp)
	if err != nil {
		return ips, err
	}

	ip := ipNet.IP.To4()
	var minIP, maxIP uint32
	for i := 0; i < 4; i++ {
		tmpCal := uint32(ip[i] & ipNet.Mask[i])
		minIP += tmpCal << ((3 - i) * 8)
	}
	ones, _ := ipNet.Mask.Size()
	maxIP = minIP | uint32(math.Pow(2, float64(32-ones))-1)

	for i := minIP; i < maxIP; i++ {
		if i&0x000000ff == 0 {
			continue
		}
		ips = append(ips, Uint32ToIP(i).String())
	}

	return ips, nil
}

func IsLocalAddress(ip net.IP, ipSegment string) bool {
	if ipSegment == "" {
		return true
	}

	_, ipRange, err := net.ParseCIDR(ipSegment)
	if err != nil {
		return false
	}

	if ipRange.Contains(ip) && !ip.IsLoopback() {
		return true
	}

	return false
}

func Uint32ToIP(intIP uint32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}
