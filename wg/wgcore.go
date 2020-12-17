package wg

import "fmt"

/*
"WGInstances":[{
	"InstanceNameReadOnly"						:"wg01",
	"InstanceServerIPReadOnly"					:"172.27.32.1",
	"InstanceServerPortReadOnly"				:22200,
	"ClientInstanceDNSServers"			:["1.1.1.1","8.8.8.8"],
	"InstanceFireWallPostUP"			:"",
	"InstanceFireWallPostDown"			:"",
	"ClientKeepAlive"					:10,
	"ClientAllowedIPsCIDR"				:["0.0.0.0/0"]
}],

*/
func (w *WGConfig) InitiateConfig(UseNAT bool, NATApaterName string, MaxInstances uint16) error {
	hclients, err := GenerateHostsAndClients(w.WGInstancesCIDR)
	if err != nil {
		return err
	}
	PortIncrement := w.WGInstancesStartPort
	w.LimitMaxInstances = MaxInstances
	for i, host := range hclients {
		if MaxInstances > 0 && i >= int(MaxInstances) {
			break
		}
		var wgInstance WGInstanceConfig
		wgInstance.InstanceNameReadOnly = fmt.Sprintf("wg%02d", i+1)
		wgInstance.InstanceServerIPCIDRReadOnly = host.HostIP
		wgInstance.InstanceServerPortReadOnly = PortIncrement
		PortIncrement++
		wgInstance.ClientInstanceDNSServers = []string{"1.1.1.1", "8.8.8.8"}
		if UseNAT {
			wgInstance.InstanceFireWallPostUP = fmt.Sprintf("iptables -A FORWARD -i %s -j ACCEPT; iptables -A FORWARD -o %s -j ACCEPT; iptables -t nat -A POSTROUTING -o %s -j MASQUERADE", wgInstance.InstanceNameReadOnly, wgInstance.InstanceNameReadOnly, NATApaterName)
			wgInstance.InstanceFireWallPostDown = fmt.Sprintf("iptables -A FORWARD -i %s -j ACCEPT; iptables -A FORWARD -o %s -j ACCEPT; iptables -t nat -A POSTROUTING -o %s -j MASQUERADE", wgInstance.InstanceNameReadOnly, wgInstance.InstanceNameReadOnly, NATApaterName)
		} else {
			wgInstance.InstanceFireWallPostUP = fmt.Sprintf("iptables -A FORWARD -i %s -j ACCEPT", wgInstance.InstanceNameReadOnly)
			wgInstance.InstanceFireWallPostDown = fmt.Sprintf("iptables -A FORWARD -i %s -j ACCEPT", wgInstance.InstanceNameReadOnly)
		}
		wgInstance.ClientKeepAlive = 10
		wgInstance.ClientAllowedIPsCIDR = []string{"0.0.0.0/0"}
		wgInstance.ClientsIP = host.HostClients
		w.WGInstances = append(w.WGInstances, wgInstance)

	}
	// log.Println(hosts)
	return nil
}

//GenerateClients Generate The Clients for first time
func (w *WGConfig) GenerateClients(clientFilePath string) error {

	return nil
}
