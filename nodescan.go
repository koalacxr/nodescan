package nodescan

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type NodeScanner struct {
	isLocal      bool
	timeout      time.Duration
	ips          []string
	ports        []int
	networkTypes []NetworkType
	maxChannel   int
}

type IpAttribute struct {
	Network NetworkType
	IP      string
	Port    int
}

type IpAttributes struct {
	sync.Mutex
	Values []IpAttribute
}

type NodeScannerFunc func(scanner *NodeScanner)

func SetIsLocal(isLocal bool) NodeScannerFunc {
	return func(scanner *NodeScanner) {
		scanner.isLocal = isLocal
	}
}

func SetTimeout(timeout time.Duration) NodeScannerFunc {
	return func(scanner *NodeScanner) {
		scanner.timeout = timeout
	}
}

func SetIps(ips []string) NodeScannerFunc {
	return func(scanner *NodeScanner) {
		scanner.ips = ips
	}
}

func SetPorts(ports []int) NodeScannerFunc {
	return func(scanner *NodeScanner) {
		scanner.ports = ports
	}
}

func SetNetworkTypes(networkTypes []NetworkType) NodeScannerFunc {
	return func(scanner *NodeScanner) {
		scanner.networkTypes = networkTypes
	}
}

func SetMaxChannel(maxChannel int) NodeScannerFunc {
	return func(scanner *NodeScanner) {
		scanner.maxChannel = maxChannel
	}
}

func NewNodeScanner(options ...NodeScannerFunc) *NodeScanner {
	newNodeScanner := &NodeScanner{
		isLocal:      true,
		timeout:      300 * time.Millisecond,
		ips:          make([]string, 0),
		ports:        []int{80, 443},
		networkTypes: []NetworkType{NetWorkTCP4},
		maxChannel:   10,
	}

	for _, op := range options {
		op(newNodeScanner)
	}

	return newNodeScanner
}

func (ns *NodeScanner) LocalCIDRs() []string {
	ns.isLocal = true
	localCIDRs, err := GetLocalNetworkIpAddresses(NetWorkTCP4, "")
	if err != nil {
		log.Printf("nodescan LocalIPs error:%v", err.Error())
	}
	return localCIDRs
}

func (ns *NodeScanner) LocalIPs() []string {
	ns.isLocal = true
	localCIDRs, err := GetLocalNetworkIpAddresses(NetWorkTCP4, "")
	if err != nil {
		log.Printf("nodescan LocalIPs error:%v", err.Error())
	}
	for i, itemCIDR := range localCIDRs {
		splitIndex := strings.Index(itemCIDR, "/")
		if splitIndex >= 0 {
			localCIDRs[i] = itemCIDR[:splitIndex]
		}
	}

	return localCIDRs
}

func (ns *NodeScanner) Scan() *IpAttributes {
	ias := &IpAttributes{
		Values: make([]IpAttribute, 0),
	}

	opIPs := ns.getIPs()
	opPorts := ns.getPorts()
	opNetworkTypes := ns.getNetworkTypes()

	var wg sync.WaitGroup
	workChannel := make(chan struct{}, ns.maxChannel)

	for _, networkType := range opNetworkTypes {
		for _, ip := range opIPs {
			for _, port := range opPorts {
				wg.Add(1)
				workChannel <- struct{}{}

				go func(locNetworkType NetworkType, locIP string, locPort int) {
					defer func() {
						wg.Done()
						<-workChannel
					}()

					if conn, err := net.DialTimeout(string(locNetworkType), fmt.Sprintf("%s:%d", locIP, locPort), ns.timeout); err == nil {
						_ = conn.Close()

						newIPAttribute := IpAttribute{
							Network: locNetworkType,
							IP:      locIP,
							Port:    locPort,
						}
						ias.Lock()
						ias.Values = append(ias.Values, newIPAttribute)
						ias.Unlock()
					}
				}(networkType, ip, port)
			}
		}
	}

	wg.Wait()

	return ias
}

func (ns *NodeScanner) getIPs() []string {
	ips := make([]string, 0)

	if ns.isLocal {
		localIps, err := GetLocalNetworkIpAddresses(NetWorkTCP4, "")
		if err != nil {
			log.Printf("nodescan GetLocalNetworkIpAddresses error:%v", err.Error())
		}
		for _, network := range localIps {
			hosts, err := GetLocalSubNetIps(network)
			if err != nil {
				log.Printf("nodescan GetLocalSubNetIps error:%v", err.Error())
			}

			ips = append(ips, hosts...)
		}
	} else {
		for _, ipItem := range ns.ips {
			tmpIP := net.ParseIP(ipItem)
			if tmpIP != nil {
				ips = append(ips, tmpIP.String())
			} else {
				log.Printf("nodescan getIPs ip invalid:%v", ipItem)
			}
		}
	}

	return ips
}

func (ns *NodeScanner) getPorts() []int {
	ports := make([]int, 0)

	for _, itemPort := range ns.ports {
		if itemPort >= 0 && itemPort <= 65535 {
			ports = append(ports, itemPort)
		} else {
			log.Printf("nodescan getPorts port invalid:%v", itemPort)
		}
	}

	return ports
}

func (ns *NodeScanner) getNetworkTypes() []NetworkType {
	networkTypes := make([]NetworkType, 0)

	for _, itemNetworkType := range ns.networkTypes {
		if CheckNetworkType(itemNetworkType) {
			networkTypes = append(networkTypes, itemNetworkType)
		} else {
			log.Printf("nodescan getNetworkTypes networkType invalid:%v", itemNetworkType)
		}
	}

	return networkTypes
}
