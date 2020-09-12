package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/koalacxr/nodescan"
)

const ConstNodeScanLogo = `
$$\   $$\                 $$\                  $$$$$$\                               
$$$\  $$ |                $$ |                $$  __$$\                              
$$$$\ $$ | $$$$$$\   $$$$$$$ | $$$$$$\        $$ /  \__| $$$$$$$\ $$$$$$\  $$$$$$$\  
$$ $$\$$ |$$  __$$\ $$  __$$ |$$  __$$\       \$$$$$$\  $$  _____|\____$$\ $$  __$$\ 
$$ \$$$$ |$$ /  $$ |$$ /  $$ |$$$$$$$$ |       \____$$\ $$ /      $$$$$$$ |$$ |  $$ |
$$ |\$$$ |$$ |  $$ |$$ |  $$ |$$   ____|      $$\   $$ |$$ |     $$  __$$ |$$ |  $$ |
$$ | \$$ |\$$$$$$  |\$$$$$$$ |\$$$$$$$\       \$$$$$$  |\$$$$$$$\\$$$$$$$ |$$ |  $$ |
\__|  \__| \______/  \_______| \_______|$$$$$$\\______/  \_______|\_______|\__|  \__|
                                        \______|                                     
                                                                                     
                                                                                     
`

func main() {
	newNodeScanner := nodescan.NewNodeScanner(nodescan.SetIsLocal(true))
	LocalIPs := newNodeScanner.LocalIPs()

	rootCmd := &cobra.Command{
		Use: "nodescan",
		Short: ConstNodeScanLogo + "local IP:" + strings.Join(LocalIPs, ", ") + "\n\n\n" +
			"nodescan is a quick scan of local or remote IP and ports. ",
		Version: "1.0",
	}

	rootCmd.AddCommand(
		ScanLocalIPs(), ScanLocalIpPorts(),
		ScanLocalNetwork(),
		ScanRemoteNetwork())

	cobra.OnInitialize()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func ScanLocalIPs() *cobra.Command {
	return &cobra.Command{
		Use:     "l",
		Aliases: []string{"localIPs"},
		Short:   "Use commands(l or localIPs) to scan the locally IPs ",
		RunE: func(cmd *cobra.Command, args []string) error {
			newNodeScanner := nodescan.NewNodeScanner(nodescan.SetIsLocal(true))
			LocalIPs := newNodeScanner.LocalIPs()
			fmt.Println("Scanned the locally IPs as follows:")
			for _, localIp := range LocalIPs {
				fmt.Println(localIp)
			}
			return nil
		},
	}
}

func ScanLocalIpPorts() *cobra.Command {
	return &cobra.Command{
		Use:     "p",
		Aliases: []string{"localPorts"},
		Short:   "Use commands(p or localPorts) to scan the locally IP ports, multiple port Numbers are spaced by ','(ex:80,443) ",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			searchPorts := parsePortParm(args[0])
			if len(searchPorts) == 0 {
				fmt.Println("params[port] is empty")
				return nil
			}
			newNodeScanner := nodescan.NewNodeScanner(nodescan.SetIsLocal(true),
				nodescan.SetPorts(searchPorts))
			LocalIPs := newNodeScanner.LocalScan()
			fmt.Println("Scanned the locally available IP ports as follows:")
			for _, ipa := range LocalIPs.Values {
				fmt.Println(fmt.Sprintf("%v:%d", ipa.IP, ipa.Port))
			}
			return nil
		},
	}
}

func ScanLocalNetwork() *cobra.Command {
	return &cobra.Command{
		Use:     "lp",
		Aliases: []string{"lanPorts"},
		Short:   "Use commands(lp or lanPorts) to scan the local network IP ports, multiple port Numbers are spaced by ','(ex:80,443) ",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}

			searchPorts := parsePortParm(args[0])
			if len(searchPorts) == 0 {
				fmt.Println("params[port] is empty")
				return nil
			}

			newNodeScanner := nodescan.NewNodeScanner(
				nodescan.SetIsLocal(true),
				nodescan.SetTimeout(300*time.Millisecond),
				nodescan.SetPorts(searchPorts),
				nodescan.SetNetworkTypes([]nodescan.NetworkType{nodescan.NetWorkTCP4}),
				nodescan.SetMaxChannel(100))
			LocalIPs := newNodeScanner.Scan()
			fmt.Println("Scanned the local network IP ports as follows:")
			for _, localIp := range LocalIPs.Values {
				fmt.Println(fmt.Sprintf("%v:%d", localIp.IP, localIp.Port))
			}
			return nil
		},
	}
}

func ScanRemoteNetwork() *cobra.Command {
	return &cobra.Command{
		Use:     "wp",
		Aliases: []string{"wanIpPorts"},
		Short:   "Use commands(wp or wanIpPorts) to scan the remote network IP ports, multiple port Numbers are spaced by ','(ex: 127.0.0.1,10.128.51.187:80,443) ",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Println("params[ip or port] invalid, example 127.0.0.1,10.128.11.187:80,443 ")
				return nil
			}

			searchIPs, searchPorts := parseIpPortParm(args[0])
			if len(searchIPs) == 0 || len(searchPorts) == 0 {
				fmt.Println(fmt.Sprintf("params[ip or port] is empty, ips[%v], ports[%v]", searchIPs, searchPorts))
				return nil
			}

			newNodeScanner := nodescan.NewNodeScanner(
				nodescan.SetIsLocal(false),
				nodescan.SetTimeout(300*time.Millisecond),
				nodescan.SetIps(searchIPs),
				nodescan.SetPorts(searchPorts),
				nodescan.SetNetworkTypes([]nodescan.NetworkType{nodescan.NetWorkTCP4}),
				nodescan.SetMaxChannel(100))
			LocalIPs := newNodeScanner.Scan()
			fmt.Println("Scanned the remote network IP ports as follows:")
			for _, localIp := range LocalIPs.Values {
				fmt.Println(fmt.Sprintf("%v:%d", localIp.IP, localIp.Port))
			}
			return nil
		},
	}
}

func parseIpPortParm(value string) ([]string, []int) {
	ips := make([]string, 0)
	ports := make([]int, 0)
	if len(value) > 0 {
		param := strings.TrimSpace(value)
		ipPortParams := strings.Split(param, ":")

		if len(ipPortParams) == 2 {
			ipSplit := strings.Split(ipPortParams[0], ",")
			for _, itemIP := range ipSplit {
				ips = append(ips, itemIP)
			}

			portSplit := strings.Split(ipPortParams[1], ",")
			for _, itemPort := range portSplit {
				portNum, err := strconv.Atoi(strings.TrimSpace(itemPort))
				if err != nil {
					continue
				}
				ports = append(ports, portNum)
			}
		}
	}

	return ips, ports
}

func parsePortParm(value string) []int {
	ports := make([]int, 0)
	if len(value) > 0 {
		portParam := strings.TrimSpace(value)
		portSplit := strings.Split(portParam, ",")
		for _, itemPort := range portSplit {
			portNum, err := strconv.Atoi(strings.TrimSpace(itemPort))
			if err != nil {
				continue
			}
			ports = append(ports, portNum)
		}
	}

	return ports
}
