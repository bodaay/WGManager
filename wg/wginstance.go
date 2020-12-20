package wg

import (
	"WGManager/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

type WGInstanceConfig struct {
	InstanceNameReadOnly         string      `json:"instance_name_read_only"`
	InstanceServerIPCIDRReadOnly string      `json:"instance_server_ipcidr_read_only"`
	InstanceServerPortReadOnly   uint16      `json:"instance_server_port_read_only"`
	InstanceEndPointHostname     string      `json:"instance_end_point_hostname"`
	ClientInstanceDNSServers     []string    `json:"client_instance_dns_servers"`
	InstanceFireWallPostUP       string      `json:"instance_fire_wall_post_up"`
	InstanceFireWallPostDown     string      `json:"instance_fire_wall_post_down"`
	InstancePubKey               string      `json:"instance_pub_key"`
	InstancePriKey               string      `json:"instance_pri_key"`
	ClientKeepAlive              uint64      `json:"client_keep_alive"`
	ClientAllowedIPsCIDR         []string    `json:"client_allowed_i_ps_cidr"`
	MaxClientsPerInstance        uint64      `json:"max_clients_per_instance"`
	WGClients                    []*WGClient `json:"WGClients"`
	//Added by imran
	InstanceUseNAT       bool   `json:"InstanceUseNAT"`
	InstanceEthernetName string `json:"InstanceEthernetName"`
}

func (wi *WGInstanceConfig) save(instancePath string) error {
	jsondata, err := json.MarshalIndent(wi, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(instancePath, jsondata, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
func (wi *WGInstanceConfig) load(instancePath string) error {
	fdata, err := ioutil.ReadFile(instancePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fdata, &wi)
	if err != nil {
		return err
	}
	return nil
}
func (wi *WGInstanceConfig) remove(confFilePath string, instancePath string) error {
	err := os.Remove(confFilePath)
	if err != nil {
		return err
	}
	err = os.Remove(instancePath)
	if err != nil {
		return err
	}
	return nil
}
func (wi *WGInstanceConfig) deploy(confpath string) error {
	confFileName := fmt.Sprintf("%s.conf", wi.InstanceNameReadOnly)
	confFileNameAndPath := path.Join(confpath, confFileName)
	err := utils.CreateFolderIfNotExists(confpath)
	if err != nil {
		return err
	}
	var sb strings.Builder
	sb.WriteString("[interface]\n")
	sb.WriteString(fmt.Sprintf("Address = %s\n", wi.InstanceServerIPCIDRReadOnly))
	sb.WriteString(fmt.Sprintf("ListenPort = %d\n", wi.InstanceServerPortReadOnly))
	sb.WriteString(fmt.Sprintf("PrivateKey = %s\n", wi.InstancePriKey))
	sb.WriteString(fmt.Sprintf("PostUp = %s\n", wi.InstanceFireWallPostUP))
	sb.WriteString(fmt.Sprintf("PostDown = %s\n", wi.InstanceFireWallPostDown))

	sb.WriteString("\n")
	sb.WriteString("\n")

	for _, wc := range wi.WGClients {
		sb.WriteString("[Peer]\n")
		sb.WriteString(fmt.Sprintf("# ClientUUID: %s, IsAllocated: %t, Allocated Timestamp:%s\n", wc.ClientUUID, wc.IsAllocated, wc.AllocatedTimestamp))
		sb.WriteString(fmt.Sprintf("PublicKey = %s\n", wc.ClientPubKey))
		sb.WriteString(fmt.Sprintf("AllowedIPs = %s\n", wc.ClientIPCIDR))
		// tempAIPSLine := ""
		// if len(wi.ClientAllowedIPsCIDR) > 0 {
		// 	for _, d := range wi.ClientAllowedIPsCIDR {
		// 		tempAIPSLine += d
		// 		tempAIPSLine += ","
		// 	}
		// 	tempAIPSLine = tempAIPSLine[:len(tempAIPSLine)-1]
		// 	sb.WriteString(fmt.Sprintf("AllowedIPs = %s\n", tempAIPSLine))
		// }
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
func (wi *WGInstanceConfig) generateServerAndClients(ipcidr string) error {
	possibleHosts, err := GenerateHostsForCIDR(ipcidr)
	if err != nil {
		return err
	}
	_, ipnet, err := net.ParseCIDR(ipcidr)
	if err != nil {
		return err
	}

	counter := 0
	for _, c := range possibleHosts {
		if counter == 0 { //skip first ip
			sz, _ := ipnet.Mask.Size()
			wi.InstanceServerIPCIDRReadOnly = fmt.Sprintf("%s/%d", c, sz)
			counter++
			continue
		}
		if counter == int(wi.MaxClientsPerInstance) {
			break
		}
		pkey, err := newPrivateKey()
		if err != nil {
			return err
		}
		wc := WGClient{
			ClientIPCIDR:       fmt.Sprintf("%s/32", c),
			GeneratedTimestamp: time.Now().Format(utils.MyTimeFormatWithoutTimeZone),
			IsAllocated:        false,
			ClientPubKey:       pkey.Public().String(),
			ClientPriKey:       pkey.String(),
		}
		// err = wgdb.InsertUpdateClient(wc)
		// if err != nil {
		// 	return err
		// }

		wi.WGClients = append(wi.WGClients, &wc)
	}
	return nil

}

func (wi *WGInstanceConfig) findClientBYIPCIDR(IPCIDR string) (*WGClient, error) {
	for _, wc := range wi.WGClients {
		if wc.ClientIPCIDR == IPCIDR {
			return wc, nil
		}
	}
	return nil, errors.New("Client Not Found")
}

func (wi *WGInstanceConfig) allocateClient(ClientUUID string, instancePath string, wgconfigPath string) (*WGClient, error) {
	foundAvailable := false
	allocatedWC := &WGClient{}
	//Check if he has been asigned an IP before
	for _, wc := range wi.WGClients {
		if wc.ClientUUID == ClientUUID {
			return nil, fmt.Errorf("ClientUUID Exists to Another IP CIDDR: %s\tinstance name: %s", wc.ClientIPCIDR, wi.InstanceNameReadOnly)
		}

	}
	for _, wc := range wi.WGClients {
		if !wc.IsAllocated {
			wc.ClientUUID = ClientUUID
			wc.IsAllocated = true
			wc.AllocatedTimestamp = time.Now().Format(utils.MyTimeFormatWithoutTimeZone)
			foundAvailable = true
			allocatedWC = wc
			break
		}
	}
	if !foundAvailable {
		return nil, fmt.Errorf("No Free IPs Available in instance: %s", wi.InstanceNameReadOnly)
	}
	instanceFileName := fmt.Sprintf("%s.json", wi.InstanceNameReadOnly)
	finalFileNameAndPath := path.Join(instancePath, instanceFileName)
	err := wi.save(finalFileNameAndPath)
	if err != nil {
		return nil, err
	}
	err = wi.deploy(wgconfigPath)
	if err != nil {
		return nil, err
	}

	return allocatedWC, nil
}
func (wi *WGInstanceConfig) revokeClientByUUID(ClientUUID string, instancePath string, wgconfigPath string) error {
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
			break
		}
	}
	instanceFileName := fmt.Sprintf("%s.json", wi.InstanceNameReadOnly)
	finalFileNameAndPath := path.Join(instancePath, instanceFileName)
	err := wi.save(finalFileNameAndPath)
	if err != nil {
		return err
	}
	err = wi.deploy(wgconfigPath)
	if err != nil {
		return err
	}
	return nil
}
func (wi *WGInstanceConfig) revokeClientByIPCIDR(IPCIDR string) error {
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
