package wg

//source: https://gist.github.com/kotakanbe/d3059af990252ba89a82
import (
	"fmt"
	"log"
	"net"
)

type HostAndClients struct {
	HostIP      string
	HostClients []string
}

//GenerateHostsAndClients Please this function is shit, don't read it
func GenerateHostsAndClients(cidr string) ([]HostAndClients, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Println(cidr)
		return nil, err
	}
	var hclients []HostAndClients
	var ips []string
	var ipsNotString []net.IP
	var hc HostAndClients
	hindex := -1
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
		ipsNotString = append(ipsNotString, ip)
		if ip[3] == 1 {
			hc = HostAndClients{}
			hc.HostIP = fmt.Sprintf("%s/24", ip.String())
			hclients = append(hclients, hc)
			hindex++
		} else {
			if hindex > -1 && ip[3] != 255 {
				hclients[hindex].HostClients = append(hclients[hindex].HostClients, fmt.Sprintf("%s/32", ip.String()))
			}

		}
	}

	return hclients, nil
}

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
