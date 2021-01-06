package resource

type WgAddInstanceRequest struct {
	IP           string   `json:"IP"`
	Port         uint16   `json:"Port"`
	DNS          []string `json:"DNS"`
	EndpointHost string   `json:"EndpointHost"`
	UseNAT       bool     `json:"UseNAT"`
	EthernetName string   `json:"EthernetName"`
	MaxClient    uint16   `json:"MaxClient"`
}
type WgRemoveInstanceRequest struct {
	Instancename string `json:"name"`
}
type WgDeploynstanceRequest struct {
	Instancename string `json:"name"`
}
