package main

const defaultJSONFileString = `{
	"APIListenAddress":"0.0.0.0",
	"APIListenPort":"6969",
	"APIUseTLS":true,
	"APITLSCert":"/etc/ssl/wgman.cert",
	"APITLSKey":"/etc/ssl/wgman.key",
	"APIAllowedIPSCIDR":["0.0.0.0/32"],
	"ClientDBPath":"/home/ubuntu/wgman/",
	"InstancesConfigPath":"/etc/wireguard/",
	"WGInstances":[
		"ClientInstanceDNSServers"			:["1.1.1.1","8.8.8.8"]
		"InstanceFireWallPostUP"			:""
		"InstanceFireWallPostDown"			:""
		"ClientKeepAlive"					:10
		"ClientAllowedIPsCIDR"				:["0.0.0.0/0"]
	],
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
	instanceName             string
	instanceServerIP         string
	instanceServerPort       string
	ClientInstanceDNSServers []string `json:"ClientInstanceDNS"`
	InstanceFireWallPostUP   string   `json:"InstanceFireWallPostUP"`
	InstanceFireWallPostDown string   `json:"InstanceFireWallPostDown"`
	ClientKeepAlive          uint64   `json:"ClientKeepAlive"`
	ClientAllowedIPsCIDR     []string `json:"ClientAllowedIPsCIDR"`
}

//ParseConfigFile Parse Config File by specified path
func (w *WGConfig) ParseConfigFile(configpath string) error {
	return nil
}

//ParseConfig Parse Config string
func (w *WGConfig) ParseConfig(configstring string) error {
	return nil
}
