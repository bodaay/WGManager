package resource

type WgAllocateClientRequest struct {
	Clientuuid   string `json:"Clientuuid"`
	Instancename string `json:"Instancename"`
}
type WgRevokeClientRequest struct {
	Clientuuid   string `json:"Clientuuid"`
	Instancename string `json:"Instancename"`
}
