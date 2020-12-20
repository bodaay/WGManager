package wg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/skip2/go-qrcode"
)

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

/*
[Interface]
Address = 10.0.22.28/32
PrivateKey =
DNS = 192.168.200.10
[Peer]
PublicKey =
Endpoint = 10.10.20.2:11222
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 10

*/

func (wg *WGClient) createClientConfigString(serverAddress string, serverPubKey string, DNSServers []string, AllowedIPs []string, Endpoint string, KeepAlive uint16, instancePort uint16) (string, error) {
	if !wg.IsAllocated {
		return "", errors.New("Client is not allocated, sorry")
	}
	var sb strings.Builder
	//server config is actually the client interface itself
	sb.WriteString("[interface]\n")
	sb.WriteString(fmt.Sprintf("Address = %s\n", wg.ClientIPCIDR))
	sb.WriteString(fmt.Sprintf("PrivateKey = %s\n", wg.ClientPriKey))
	tempDNSLine := ""
	if len(DNSServers) > 0 {
		for _, d := range DNSServers {
			tempDNSLine += d
			tempDNSLine += ","
		}
		tempDNSLine = tempDNSLine[:len(tempDNSLine)-1]
		sb.WriteString(fmt.Sprintf("DNS = %s\n", tempDNSLine))
	}
	//peer config is the wg instance
	sb.WriteString("[Peer]\n")
	sb.WriteString(fmt.Sprintf("PublicKey = %s\n", serverPubKey))
	sb.WriteString(fmt.Sprintf("Endpoint = %s:%d\n", Endpoint, instancePort))
	tempAIPSLine := ""
	if len(AllowedIPs) > 0 {
		for _, d := range AllowedIPs {
			tempAIPSLine += d
			tempAIPSLine += ","
		}
		tempAIPSLine = tempAIPSLine[:len(tempAIPSLine)-1]
		sb.WriteString(fmt.Sprintf("AllowedIPs = %s\n", tempAIPSLine))
	}
	sb.WriteString(fmt.Sprintf("PersistentKeepalive = %d\n", KeepAlive))
	return sb.String(), nil
}

func (wg *WGClient) createClientConfigQRCodePicture(content string, filepath string) error {
	data, err := wg.createClientConfigQRCode(content)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
func (wg *WGClient) createClientConfigQRCode(content string) ([]byte, error) {
	imgbytes, err := qrcode.Encode(content, qrcode.High, 256)
	if err != nil {
		return nil, err
	}
	return imgbytes, nil
}
