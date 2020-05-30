package scanner

import "sync"

type PortScanPayload struct {
	AWSKey    string
	AWSSecret string
	AWSToken  string
	AWSRegion string
}

type PortScanResult struct {
	PublicIPs []string
	Ports     map[string]map[int]string
}

func (p *PortScanResult) Set(key string, value map[int]string) {
	var mu sync.Mutex
	mu.Lock()
	p.Ports[key] = value
	mu.Unlock()
}
