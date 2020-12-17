package wg

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const DefaultJSONFileString = `{
	"APIListenAddress":"0.0.0.0",
	"APIListenPort":6969,
	"APIUseTLS":true,
	"APITLSCert":"/etc/ssl/wgman.cert",
	"APITLSKey":"/etc/ssl/wgman.key",
	"APIAllowedIPSCIDR":["0.0.0.0/32"],
	"ClientDBPath":"wgman",
	"InstancesConfigPath":"/etc/wireguard/",
	"WGInstances":[],
	"WGInstancesCIDR":"172.27.32.0/20",
	"WGInstancesStartPort": 22200,
	"WGGlobalEndPointHostName": "wg.MyWireGuard.org",
	"LimitMaxInstances":0
}`

//WGConfig Global Configuration For WGManager
type WGConfig struct {
	APIListenAddress         string             `json:"APIListenAddress"`
	APIListenPort            uint16             `json:"APIListenPort"`
	APIUseTLS                bool               `json:"APIUseTLS"`
	APITLSCert               string             `json:"APITLSCert"`
	APITLSKey                string             `json:"APITLSKey"`
	APIAllowedIPS            []string           `json:"APIAllowedIPS"`
	ClientDBPath             string             `json:"ClientDBPath"`
	InstancesConfigPath      string             `json:"InstancesConfigPath"`
	WGInstances              []WGInstanceConfig `json:"WGInstances"`
	WGInstancesCIDR          string             `json:"WGInstancesCIDR"`
	WGInstancesStartPort     uint16             `json:"WGInstancesStartPort"`
	WGGlobalEndPointHostName string             `json:"WGGlobalEndPointHostName"`
	LimitMaxInstances        uint16             `json:"LimitMaxInstances"`
}

//WGInstanceConfig Per Instance Configuration
type WGInstanceConfig struct {
	InstanceNameReadOnly         string     `json:"InstanceNameReadOnly"`
	InstanceServerIPCIDRReadOnly string     `json:"InstanceServerIPCIDRReadOnly"`
	InstanceServerPortReadOnly   uint16     `json:"InstanceServerPortReadOnly"`
	ClientInstanceDNSServers     []string   `json:"ClientInstanceDNSServers"`
	InstanceFireWallPostUP       string     `json:"InstanceFireWallPostUP"`
	InstanceFireWallPostDown     string     `json:"InstanceFireWallPostDown"`
	InstancePubKey               string     `json:"InstancePubKey"`
	InstancePriKey               string     `json:"InstancePriKey"`
	ClientKeepAlive              uint64     `json:"ClientKeepAlive"`
	ClientAllowedIPsCIDR         []string   `json:"ClientAllowedIPsCIDR"`
	ClientsIP                    []string   `json:"-"`
	WGDB                         wgdb       `json:"-"`
	WGClients                    []WGClient `json:"-"`
}

type WGClient struct {
	ClientIPCIDR       string `json:"ClientIPCIDR"`
	ClientPubKey       string `json:"ClientPubKey"`
	ClientPriKey       string `json:"ClientPriKey"`
	IsAllocated        bool   `json:"IsAllocated"`
	ClientUUID         string `json:"ClientUUID"`
	InsertedTimestamp  string `json:"InsertedTimestamp"`
	AllocatedTimestamp string `json:"AllocatedTimestamp"`
	RevokedTimestamp   string `json:"RevokedTimestamp"`
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
