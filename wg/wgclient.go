package wg

type WGClient struct {
	ClientIPCIDR       string `json:"client_ipcidr"`
	ClientPubKey       string `json:"client_pub_key"`
	ClientPriKey       string `json:"client_pri_key"`
	IsAllocated        bool   `json:"is_allocated"`
	ClientUUID         string `json:"client_uuid"`
	GeneratedTimestamp string `json:"generated_timestamp"`
	AllocatedTimestamp string `json:"allocated_timestamp"`
	RevokedTimestamp   string `json:"revoked_timestamp"`
}
