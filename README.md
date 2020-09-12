# nodescan

```text
$$\   $$\                 $$\                  $$$$$$\                               
$$$\  $$ |                $$ |                $$  __$$\                              
$$$$\ $$ | $$$$$$\   $$$$$$$ | $$$$$$\        $$ /  \__| $$$$$$$\ $$$$$$\  $$$$$$$\  
$$ $$\$$ |$$  __$$\ $$  __$$ |$$  __$$\       \$$$$$$\  $$  _____|\____$$\ $$  __$$\ 
$$ \$$$$ |$$ /  $$ |$$ /  $$ |$$$$$$$$ |       \____$$\ $$ /      $$$$$$$ |$$ |  $$ |
$$ |\$$$ |$$ |  $$ |$$ |  $$ |$$   ____|      $$\   $$ |$$ |     $$  __$$ |$$ |  $$ |
$$ | \$$ |\$$$$$$  |\$$$$$$$ |\$$$$$$$\       \$$$$$$  |\$$$$$$$\\$$$$$$$ |$$ |  $$ |
\__|  \__| \______/  \_______| \_______|$$$$$$\\______/  \_______|\_______|\__|  \__|
                                        \______|                                     
                                                                                     

```

nodescan is a quick scan of local or remote IP and ports.

#### Installing
Using nodescan is easy. 
First, use go get to install the latest version of the library. 
This command will install the nodescan generator executable along with the library and its dependencies:
```commandline
$ go get -u github.com/koalacxr/nodescan
$ cd ./nodescan/cmds/nodescan
$ go install
```

##### Get nodescan help information
```commandline
$ nodescan --help

Usage:
  nodescan [command]

Available Commands:
  help        Help about any command
  l           Use commands(l or localIPs) to scan the locally IPs 
  lp          Use commands(lp or lanPorts) to scan the local network IP ports, multiple port Numbers are spaced by ','(ex:80,443) 
  p           Use commands(p or localPorts) to scan the locally IP ports, multiple port Numbers are spaced by ','(ex:80,443) 
  wp          Use commands(wp or wanIpPorts) to scan the remote network IP ports, multiple port Numbers are spaced by ','(ex: 127.0.0.1,10.128.51.187:80,443) 

Flags:
  -h, --help      help for nodescan
  -v, --version   version for nodescan

Use "nodescan [command] --help" for more information about a command.
```

### Example
#### Get the local IPs with the following command.
```commandline
$ nodescan l
$ nodescan localIPs
```

#### Scan for open ports corresponding to local IPs.
```commandline
$ nodescan ln 80,443,3306
$ nodescan lanPorts 80,443,3306
```

#### Scan the port corresponding to the local LAN. Multiple port Numbers separated by ','.
```commandline
$ nodescan p 80,443,3306
$ nodescan localPorts 80,443,3306
```

#### Scan the WAN IPs and port number. Multiple IPs and multiple ports separated by ','. IPs and ports separated by ':'.
```commandline
$ nodescan wp 127.0.0.1,10.128.11.187:80,443
$ nodescan wanIpPorts 127.0.0.1,10.128.11.187:80,443
```