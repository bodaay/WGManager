package wg

import (
	"WGManager/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"runtime"
	"sync"
)

const defaultAPIListenAdderss = "0.0.0.0"
const defaultAPIListenPort = 6969
const defaultAdminAPIListenAdderss = "0.0.0.0"
const defaultAdminAPIListenPort = 1129
const defaultAPIUseTLS = false
const defaultAPICertFile = "/etc/ssl/wgman/wgman.cert"
const defaultAPIKeyFile = "/etc/ssl/wgman/wgman.key"

var defaultAllowedIPsCIDR = []string{"0.0.0.0/32"}

const defaultFirstInstanceCIDR = "172.27.36.0/22"
const defaultFirstInstancePort = 22200
const defaultInstanceEndPointHostName = "165.22.30.240"
const defaultInstanceConfigPath = "wginstance"

//WGConfig Global Configuration For WGManager
type WGConfig struct {
	sync.Mutex
	APIListenAddress           string              `json:"api_listen_address"`
	APIListenPort              uint16              `json:"api_listen_port"`
	AdminAPIListenAddress      string              `json:"admin_api_listen_address"`
	AdminAPIListenPort         uint16              `json:"admin_api_listen_port"`
	APIUseTLS                  bool                `json:"api_use_tls"`
	APITLSCert                 string              `json:"apitls_cert"`
	APITLSKey                  string              `json:"apitls_key"`
	APIAllowedIPS              []string            `json:"api_allowed_ips"`
	InstancesConfigPath        string              `json:"instances_config_path"`
	WGInsatncesServiceFilePath string              `json:"wg_insatnces_service_file_path"`
	WGInstances                []*WGInstanceConfig `json:"wg_instances"`
}

//WGInstanceConfig Per Instance Configuration

//CreateDefaultconfig Create Default Config file based on our constants
func (w *WGConfig) CreateDefaultconfig(configpath string) (*WGConfig, error) {
	var wgdefault WGConfig
	wgdefault.APIListenAddress = defaultAPIListenAdderss
	wgdefault.APIListenPort = defaultAPIListenPort
	wgdefault.AdminAPIListenAddress = defaultAdminAPIListenAdderss
	wgdefault.AdminAPIListenPort = defaultAdminAPIListenPort
	wgdefault.APIUseTLS = defaultAPIUseTLS
	wgdefault.APITLSCert = defaultAPICertFile
	wgdefault.APITLSKey = defaultAPIKeyFile
	wgdefault.APIAllowedIPS = defaultAllowedIPsCIDR
	wgdefault.InstancesConfigPath = defaultInstanceConfigPath

	if runtime.GOOS == "windows" {
		wgdefault.WGInsatncesServiceFilePath = "etc/wireguard/"
	} else {
		wgdefault.WGInsatncesServiceFilePath = "/etc/wireguard/"
	}
	err := utils.CreateFolderAllIfNotExists(wgdefault.WGInsatncesServiceFilePath)
	if err != nil {
		return nil, err
	}
	err = wgdefault.SaveConfigFile(configpath)
	if err != nil {
		return nil, err
	}
	return &wgdefault, nil
}

func (w *WGConfig) FindInstanceByIPAndPort(ipcidr string, port uint16) (*WGInstanceConfig, error) {
	newInstanceIP, _, err := net.ParseCIDR(ipcidr)
	if err != nil {
		return nil, err
	}
	for _, i := range w.WGInstances {
		_, ipnet, err := net.ParseCIDR(i.InstanceServerIPCIDRReadOnly)
		if err != nil {
			return nil, err
		}
		if ipnet.Contains(newInstanceIP) {
			return i, errors.New("IP confilct with instance")
		}
		if i.InstanceServerPortReadOnly == port {
			return i, errors.New("Port Conflict with instance")
		}
	}
	return nil, nil
}
func (w *WGConfig) FindInstanceByName(instanceName string) (*WGInstanceConfig, error) {
	for _, i := range w.WGInstances {
		if i.InstanceNameReadOnly == instanceName {
			return i, nil
		}
	}
	return nil, errors.New("Could not find instance with the name: %s" + instanceName)
}
func (w *WGConfig) LoadInstancesFiles() error {
	instacesFiles := utils.GetMeFileListInFolders(w.InstancesConfigPath, ".json", true, false, true)
	for _, ifile := range instacesFiles {
		var wginstance WGInstanceConfig
		err := wginstance.load(ifile)
		if err != nil {
			return err
		}
		w.WGInstances = append(w.WGInstances, &wginstance)

	}
	return nil
}

func (w *WGConfig) AllocateClient(instanceName string, clientuuid string) ([]byte, error) {
	wi, err := w.FindInstanceByName(instanceName)
	if err != nil {
		return nil, fmt.Errorf("instance name not found: %s", instanceName)
	}
	w.Lock()
	defer w.Unlock()
	wc, err := wi.allocateClient(clientuuid, w.InstancesConfigPath, w.WGInsatncesServiceFilePath)
	if err != nil {
		return nil, err
	}
	if wc != nil {
		qrcontent, err := wc.createClientConfigString(wi.InstanceServerIPCIDRReadOnly, wi.InstancePubKey, wi.ClientInstanceDNSServers, wi.ClientAllowedIPsCIDR, wi.InstanceEndPointHostname, uint16(wi.ClientKeepAlive), wi.InstanceServerPortReadOnly)
		if err != nil {
			return nil, err
		}
		qrbytes, err := wc.createClientConfigQRCode(qrcontent)
		if err != nil {
			return nil, err
		}
		return qrbytes, nil
	}
	return nil, errors.New("Client Allocaiton failed")
}

func (w *WGConfig) RevokeClient(instanceName string, clientuuid string) error {
	wi, err := w.FindInstanceByName(instanceName)
	if err != nil {
		return fmt.Errorf("instance name not found: %s", instanceName)
	}
	w.Lock()
	defer w.Unlock()
	err = wi.revokeClientByUUID(clientuuid, w.InstancesConfigPath, w.WGInsatncesServiceFilePath)
	if err != nil {
		return err
	}
	output, err := restartWGInstance(instanceName)
	if err != nil {
		return err
	}
	log.Println(output)

	return nil
}
func (w *WGConfig) DeployAllInstances() error {
	for _, wi := range w.WGInstances {
		err := w.DeployInstanceByName(wi.InstanceNameReadOnly)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *WGConfig) DeployInstanceByName(instanceName string) error {
	wi, err := w.FindInstanceByName(instanceName)
	if err != nil {
		return fmt.Errorf("instance name not found: %s", instanceName)
	}
	w.Lock()
	defer w.Unlock()
	err = wi.deploy(w.WGInsatncesServiceFilePath)
	if err != nil {
		return err
	}
	output, err := restartWGInstance(instanceName)
	if err != nil {
		return err
	}
	log.Println(output)
	return nil
}
func (w *WGConfig) RemoveInstanceByName(instanceName string) error {
	wi, err := w.FindInstanceByName(instanceName)
	if err != nil {
		return fmt.Errorf("instance name not found: %s", instanceName)
	}
	w.Lock()
	defer w.Unlock()
	finalFileNameAndPath := path.Join(w.InstancesConfigPath, wi.InstanceNameReadOnly+".json")
	output, err := stopWGInstance(instanceName)
	if err != nil {
		return err
	}
	log.Println(output)
	confFilePath := path.Join(w.WGInsatncesServiceFilePath, fmt.Sprintf("%s.conf", instanceName))

	err = wi.remove(confFilePath, finalFileNameAndPath)
	for i, v := range w.WGInstances {
		if v == wi {
			w.WGInstances = append(w.WGInstances[:i], w.WGInstances[i+1:]...)
		}
	}
	if err != nil {
		return err
	}
	return nil
}

//ParseConfigFile Parse Config File by specified path
func (w *WGConfig) ParseConfigFile(configpath string) error {
	data, err := ioutil.ReadFile(configpath)
	if err != nil {
		return err
	}
	err = w.ParseConfig(string(data))
	if err != nil {
		return err
	}
	return nil
}

//ParseConfig Parse Config string
func (w *WGConfig) ParseConfig(configstring string) error {
	err := json.Unmarshal([]byte(configstring), w)
	if err != nil {
		return err
	}

	return nil
}

//SaveConfigFile Save the file into the specified path
func (w *WGConfig) SaveConfigFile(configpath string) error {
	jsondata, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configpath, jsondata, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (w *WGConfig) CreateNewInstance(instanceCIDR string, instancePort uint16, instanceDNS []string, UseNAT bool, EthernetAdapaterName string, MaxClients uint64) error {

	var wgInstance WGInstanceConfig
	_, err := w.FindInstanceByIPAndPort(instanceCIDR, instancePort)
	if err != nil {
		return err
	}
	wgInstance.InstanceNameReadOnly = fmt.Sprintf("wg%02d", len(w.WGInstances)+1)
	wgInstance.MaxClientsPerInstance = MaxClients
	wgInstance.InstanceServerPortReadOnly = instancePort
	wgInstance.ClientInstanceDNSServers = instanceDNS
	if UseNAT {
		wgInstance.InstanceFireWallPostUP = fmt.Sprintf("iptables -A FORWARD -i %s -j ACCEPT; iptables -t nat -A POSTROUTING -o %s -j MASQUERADE", wgInstance.InstanceNameReadOnly, EthernetAdapaterName)
		wgInstance.InstanceFireWallPostDown = fmt.Sprintf("iptables -D FORWARD -i %s -j ACCEPT; iptables -t nat -D POSTROUTING -o %s -j MASQUERADE", wgInstance.InstanceNameReadOnly, EthernetAdapaterName)
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
	wgInstance.InstanceEndPointHostname = defaultInstanceEndPointHostName
	wgInstance.InstancePubKey = pkey.Public().String()
	wgInstance.InstancePriKey = pkey.String()
	err = wgInstance.generateServerAndClients(instanceCIDR) //// wgInstance.InstanceServerIPCIDRReadOnly will be set using this function,
	if err != nil {
		return err
	}

	w.WGInstances = append(w.WGInstances, &wgInstance)

	//saving...
	err = utils.CreateFolderIfNotExists(w.InstancesConfigPath)
	if err != nil {
		return err
	}
	instanceFileName := fmt.Sprintf("%s.json", wgInstance.InstanceNameReadOnly)
	finalFileNameAndPath := path.Join(w.InstancesConfigPath, instanceFileName)
	w.Lock()
	defer w.Unlock()
	err = wgInstance.save(finalFileNameAndPath)
	if err != nil {
		return err
	}
	return nil
}
