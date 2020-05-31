package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RingierIMU/rsb-go-lib/v9/base"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"github.com/zebroc/aws-ip-scanner/scanner"
	"net/http"
	"time"
)

func eventAction(c *gin.Context) {
	var event base.RsbEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if event.Event != "AWSAccountPubScan" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only AWSAccountPubScan events accepted here"})
		return
	}

	var portScanPayload scanner.PortScanPayload
	errUnmarshal := json.Unmarshal(event.Payload, &portScanPayload)
	if errUnmarshal != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect payload format: " + errUnmarshal.Error()})
		return
	}

	if portScanPayload.AWSKey == "" || portScanPayload.AWSSecret == "" || portScanPayload.AWSRegion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload must include values for AWSKey, AWSSecret and AWSRegion (AWSSessionToken is optional for STS)"})
		return
	}

	if portScanPayload.SendCallback {
		go func() {
			portScanResult := scanner.Scan(portScanPayload)
			event := base.RsbEvent{
				Event:            "AWSAccountPubScanResult",
				VentureConfigID:  "",
				VentureReference: "",
				CreatedAt:        time.Now().Format(time.RFC3339),
				Culture:          "en_EN",
				Payload:          portScanResult.JSON(),
			}
			post(portScanPayload.CallbackURL, string(event.JSON()))
			log.Debugf("Callback sent to %s", portScanPayload.CallbackURL)
		}()
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Callback with resulst will be sent to %s", portScanPayload.CallbackURL)})
	} else {
		portScanResult := scanner.Scan(portScanPayload)
		c.JSON(http.StatusOK, portScanResult)
	}
}

func post(url string, jsonData string) {
	var jsonStr = []byte(jsonData)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	defer resp.Body.Close()
}
