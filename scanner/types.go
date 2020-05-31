package scanner

import (
	"encoding/json"
	"sync"
)

type PortScanPayload struct {
	AWSKey          string `json:"aws_key"`
	AWSSecret       string `json:"aws_secret"`
	AWSSessionToken string `json:"aws_session_token"`
	AWSRegion       string `json:"aws_region"`
	SendCallback    bool   `json:"send_callback"`
	CallbackURL     string `json:"callback_url"`
}

type PortScanResult struct {
	mu        sync.Mutex
	PublicIPs []string
	Ports     map[string]map[int]string
}

func (p *PortScanResult) Set(key string, value map[int]string) {
	p.mu.Lock()
	p.Ports[key] = value
	p.mu.Unlock()
}

func (p *PortScanResult) JSON() []byte {
	res, _ := json.Marshal(p)
	return res
}
