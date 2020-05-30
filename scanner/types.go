package scanner

type PortScanPayload struct {
	AWSKey     string
	AWSSecret  string
	AWSRregion string
}

type PortScanResult struct {
	PublicIPs []string
	Ports     map[string]map[string]string
}
