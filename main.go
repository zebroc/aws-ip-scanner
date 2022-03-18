package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zebroc/aws-ip-scanner/aws"
	"github.com/zebroc/aws-ip-scanner/scanner"
)

var (
	ports = []int{21, 22, 80, 443, 3306, 6379, 8080, 11211}
)

func init() {
	if p := os.Getenv("AWS_IP_SCANNER_PORTS"); p != "" {
		ports = []int{}
		for _, portStr := range strings.Split(p, ",") {
			port, _ := strconv.Atoi(portStr)
			ports = append(ports, port)
		}
	}
}

func main() {
	ips, err := aws.GetIPs()
	if err != nil {
		fmt.Printf("Unable to fetch IPs: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(ips))
	scan(&wg, ips)
	wg.Wait()
}

func scan(wg *sync.WaitGroup, ips []string) {
	for _, ip := range ips {
		go func(ip string, wg *sync.WaitGroup) {
			defer wg.Done()
			hostState := fmt.Sprintf("%s\t", ip)
			openPorts := make(map[int]string)
			for _, port := range ports {
				state := scanner.ScanPort(ip, port, time.Second)
				if strings.Contains(state, "open") {
					openPorts[port] = state
					hostState = hostState + fmt.Sprintf("%d", port) + ", "
				}
			}

			if _, ok := openPorts[80]; ok {
				body, code, err := getURL(fmt.Sprintf("http://%s", ip))
				if err != nil {
					hostState = strings.TrimSuffix(hostState, ", ")
					hostState = fmt.Sprintf("%s\tHTTP (%d): Error (%s)", hostState, code, err.Error())
				} else {
					hostState = strings.TrimSuffix(hostState, ", ")
					hostState = fmt.Sprintf("%s\tHTTP (%d): %s", hostState, code, strings.Replace(body, "\n", "", 999))
				}
			}

			if _, ok := openPorts[443]; ok {
				body, code, err := getURL(fmt.Sprintf("https://%s", ip))
				if err != nil {
					hostState = strings.TrimSuffix(hostState, ", ")
					hostState = fmt.Sprintf("%s\tHTTPS (%d): Error (%s)", hostState, code, err.Error())
				} else {
					hostState = strings.TrimSuffix(hostState, ", ")
					hostState = fmt.Sprintf("%s\tHTTPS (%d): %s", hostState, code, strings.Replace(body, "\n", "", 999))
				}
			}

			if _, ok := openPorts[8080]; ok {
				body, code, err := getURL(fmt.Sprintf("http://%s", ip))
				if err != nil {
					hostState = strings.TrimSuffix(hostState, ", ")
					hostState = fmt.Sprintf("%s\tHTTP_ALT (%d): Error (%s)", hostState, code, err.Error())
				} else {
					hostState = strings.TrimSuffix(hostState, ", ")
					hostState = fmt.Sprintf("%s\tHTTP-ALT (%d): %s", hostState, code, strings.Replace(body, "\n", "", 999))
				}
			}

			if hostState != fmt.Sprintf("%s\t", ip) {
				hostState = strings.TrimSuffix(hostState, ", ")
				fmt.Printf(hostState + "\n")
			}
		}(ip, wg)
	}
}

func getURL(url string) (string, int, error) {
	success := false

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	for i := 1; !success && i <= 3; i++ {
		req, _ := http.NewRequest("GET", url, bytes.NewBuffer(nil))

		resp, _ := httpClient.Do(req)

		if resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			success = true
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			return buf.String(), resp.StatusCode, nil
		} else if resp != nil && resp.StatusCode == http.StatusNotFound {
			return "", http.StatusNotFound, nil
		} else {
			time.Sleep(time.Duration(i) * time.Second)
		}

		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}

	return "", http.StatusRequestTimeout, errors.New("timeout")
}
