package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"net/http"
	"time"
)

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
