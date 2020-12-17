package wg

import (
	"WGManager/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
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

		// for _, cip := range host.HostClients {
		// 	wgInstance.ClientsIP = append(wgInstance.ClientsIP, cip)
		// }
		// wgInstance.WGClients = make([]WGClient, 0)
		w.WGInstances = append(w.WGInstances, wgInstance)

	}
	// log.Println(hosts)
	return nil
}

//GenerateAllClients Generate The Clients for first time
func (w *WGConfig) GenerateAllClients() error {
	for k := range w.WGInstances {
		w.WGInstances[k].GenerateNewClients(w.ClientDBPath)
		// log.Printf("1-Total Generated for interface: %s is %d", w.WGInstances[k].InstanceNameReadOnly, len(w.WGInstances[k].WGClients))
	}
	return nil
}

//GenerateAllClients Generate The Clients for first time
func (w *WGConfig) ApplyAllConfigs() error {
	// log.Printf("2-Total Generated for interface: %s is %d", w.WGInstances[0].InstanceNameReadOnly, len(w.WGInstances[0].WGClients))
	for k := range w.WGInstances {
		w.WGInstances[k].Apply(w.InstancesConfigPath)
		// log.Printf("2-Total Generated for interface: %s is %d", w.WGInstances[k].InstanceNameReadOnly, len(w.WGInstances[k].WGClients))
	}

	return nil
}
func (wi *WGInstanceConfig) Save(dbPath string) error {
	return nil
}
func (wi *WGInstanceConfig) Load(dbPath string) error {
	return nil
}
func (wi *WGInstanceConfig) Apply(confpath string) error {
	confFileName := fmt.Sprintf("%s.conf", wi.InstanceNameReadOnly)
	confFileNameAndPath := path.Join(confpath, confFileName)
	err := utils.CreateFolderIfNotExists(confpath)
	if err != nil {
		return err
	}
	var sb strings.Builder
	sb.WriteString("[interface]\n")
	sb.WriteString(fmt.Sprintf("PrivateKey = %s\n", wi.InstancePriKey))
	sb.WriteString(fmt.Sprintf("Address = %s\n", wi.InstanceServerIPCIDRReadOnly))
	sb.WriteString(fmt.Sprintf("ListenPort = %d\n", wi.InstanceServerPortReadOnly))
	sb.WriteString(fmt.Sprintf("PostUp = %s\n", wi.InstanceFireWallPostUP))
	sb.WriteString(fmt.Sprintf("PostDown = %s\n", wi.InstanceFireWallPostDown))
	tempDNSLine := ""
	if len(wi.ClientInstanceDNSServers) > 0 {
		for _, d := range wi.ClientInstanceDNSServers {
			tempDNSLine += d
			tempDNSLine += ","
		}
		tempDNSLine = tempDNSLine[:len(tempDNSLine)-1]
		sb.WriteString(fmt.Sprintf("DNS = %s\n", tempDNSLine))
	}
	sb.WriteString("\n")
	sb.WriteString("\n")

	for _, wc := range wi.WGClients {
		sb.WriteString("[Peer]\n")
		sb.WriteString(fmt.Sprintf("# ClientUUID: %s, IsAllocated: %t, Allocated Timestamp:%s\n", wc.ClientUUID, wc.IsAllocated, wc.AllocatedTimestamp))
		sb.WriteString(fmt.Sprintf("PublicKey = %s\n", wc.ClientPubKey))
		tempAIPSLine := ""
		if len(wi.ClientAllowedIPsCIDR) > 0 {
			for _, d := range wi.ClientAllowedIPsCIDR {
				tempAIPSLine += d
				tempAIPSLine += ","
			}
			tempAIPSLine = tempAIPSLine[:len(tempAIPSLine)-1]
			sb.WriteString(fmt.Sprintf("AllowedIPs = %s\n", tempAIPSLine))
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
	}
	err = ioutil.WriteFile(confFileNameAndPath, []byte(sb.String()), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

//GenerateNewClients Generate The Clients for first time
func (wi *WGInstanceConfig) GenerateNewClients(dbPath string) error {

	// dbFileName := fmt.Sprintf("%s.buntdb", wi.InstanceNameReadOnly)
	// dbFileNameAndPath := path.Join(dbPath, dbFileName)
	// err := utils.CreateFolderIfNotExists(dbPath)
	// if err != nil {
	// 	return err
	// }
	// var wgdb wgdb
	// err = wgdb.openDB(dbFileNameAndPath)
	// if err != nil {
	// 	return err
	// }
	for _, c := range wi.ClientsIP {
		pkey, err := newPrivateKey()
		if err != nil {
			return err
		}
		wc := WGClient{
			ClientIPCIDR:       c,
			GeneratedTimestamp: time.Now().Format(utils.MyTimeFormatWithoutTimeZone),
			IsAllocated:        false,
			ClientPubKey:       pkey.Public().String(),
			ClientPriKey:       pkey.String(),
		}
		// err = wgdb.InsertUpdateClient(wc)
		// if err != nil {
		// 	return err
		// }
		wi.WGClients = append(wi.WGClients, wc)
	}

	return nil

}

func (wi *WGInstanceConfig) FindClientBYIPCIDR(IPCIDR string) (*WGClient, error) {
	for _, wc := range wi.WGClients {
		if wc.ClientIPCIDR == IPCIDR {
			return &wc, nil
		}
	}
	return nil, errors.New("Client Not Found")
}

func (wi *WGInstanceConfig) AllocateClient(ClientUUID string) error {
	foundAvailable := false
	//Check if he has been asigned an IP before
	for _, wc := range wi.WGClients {
		if wc.ClientUUID == ClientUUID {
			return fmt.Errorf("ClientUUID Exists to Another IP CIDDR: %s\tinstance name: %s", wc.ClientIPCIDR, wi.InstanceNameReadOnly)
		}

	}
	for _, wc := range wi.WGClients {
		if !wc.IsAllocated {
			wc.ClientUUID = ClientUUID
			wc.IsAllocated = true
			wc.AllocatedTimestamp = time.Now().Format(utils.MyTimeFormatWithoutTimeZone)
			foundAvailable = true
		}
	}
	if !foundAvailable {
		return fmt.Errorf("No Free IPs Available in instance: %s", wi.InstanceNameReadOnly)
	}
	return nil
}
func (wi *WGInstanceConfig) RevokeClientByUUID(ClientUUID string) error {
	for _, wc := range wi.WGClients {
		if wc.ClientUUID == ClientUUID {
			wc.ClientUUID = ""
			wc.ClientIPCIDR = ""
			wc.IsAllocated = false
			//we  have to change the keys
			pkey, err := newPrivateKey()
			if err != nil {
				return err
			}
			wc.ClientPubKey = pkey.Public().String()
			wc.ClientPriKey = pkey.String()
			wc.RevokedTimestamp = time.Now().Format(utils.MyTimeFormatWithoutTimeZone)
			return nil
		}
	}
	return nil
}
func (wi *WGInstanceConfig) RevokeClientByIPCIDR(IPCIDR string) error {
	for _, wc := range wi.WGClients {
		if wc.ClientIPCIDR == wc.ClientIPCIDR {
			wc.ClientUUID = ""
			wc.ClientIPCIDR = ""
			wc.IsAllocated = false
			//we  have to change the keys
			pkey, err := newPrivateKey()
			if err != nil {
				return err
			}
			wc.ClientPubKey = pkey.Public().String()
			wc.ClientPriKey = pkey.String()
			wc.RevokedTimestamp = time.Now().Format(utils.MyTimeFormatWithoutTimeZone)
		}
	}
	return nil
}

// func (w *WGConfig) LoadClients() error {
// 	for _, instance := range w.WGInstances {
// 		dbFileName := fmt.Sprintf("%s.buntdb", instance.InstanceNameReadOnly)
// 		dbFileNameAndPath := path.Join(w.ClientDBPath, dbFileName)
// 		err := utils.CreateFolderIfNotExists(w.ClientDBPath)
// 		if err != nil {
// 			return err
// 		}
// 		var wgdb wgdb
// 		err = wgdb.openDB(dbFileNameAndPath)
// 		if err != nil {
// 			return err
// 		}

// 	}
// }
// func WriteServiceConfigFiles(instance *WGInstanceConfig, clientDBPath string) error {

// }
