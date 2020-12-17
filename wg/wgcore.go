package wg

import (
	"WGManager/utils"
	"fmt"
	"path"
)

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
		//generate instances keys
		pkey, err := newPrivateKey()
		if err != nil {
			return err
		}
		wgInstance.InstancePubKey = pkey.Public().String()
		wgInstance.InstancePriKey = pkey.String()
		wgInstance.ClientsIP = host.HostClients
		w.WGInstances = append(w.WGInstances, wgInstance)

	}
	// log.Println(hosts)
	return nil
}

//GenerateAllClients Generate The Clients for first time
func (w *WGConfig) GenerateAllClients() error {
	for _, instance := range w.WGInstances {

		dbFileName := fmt.Sprintf("%s.buntdb", instance.InstanceNameReadOnly)
		dbFileNameAndPath := path.Join(w.ClientDBPath, dbFileName)
		err := utils.CreateFolderIfNotExists(w.ClientDBPath)
		if err != nil {
			return err
		}
		var wgdb wgdb
		err = wgdb.openDB(dbFileNameAndPath)
		if err != nil {
			return err
		}
		for _, c := range instance.ClientsIP {
			pkey, err := newPrivateKey()
			if err != nil {
				return err
			}

			wc := &WGClient{
				ClientIPCIDR:      c,
				InsertedTimestamp: "",
				ClientPubKey:      pkey.Public().String(),
				ClientPriKey:      pkey.String(),
			}
			err = wgdb.InsertUpdateClient(wc)
			if err != nil {
				return err
			}
		}

		// dbClient, _ := wgdb.GetClient("172.27.32.8/32")
		// if dbClient != nil {
		// 	log.Println(dbClient)
		// }

	}
	return nil
}
func (w *WGConfig) LoadClients() error {
	for _, instance := range w.WGInstances {
		dbFileName := fmt.Sprintf("%s.buntdb", instance.InstanceNameReadOnly)
		dbFileNameAndPath := path.Join(w.ClientDBPath, dbFileName)
		err := utils.CreateFolderIfNotExists(w.ClientDBPath)
		if err != nil {
			return err
		}
		var wgdb wgdb
		err = wgdb.openDB(dbFileNameAndPath)
		if err != nil {
			return err
		}

	}
}
func WriteServiceConfigFiles(instance *WGInstanceConfig, clientDBPath string) error {

}
