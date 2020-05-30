package request

import (
	"encoding/json"
	"github.com/RingierIMU/rsb-go-lib/v9/base"
	"github.com/gin-gonic/gin"
	"github.com/zebroc/aws-ip-scanner/scanner"
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload must include values for AWSKey, AWSSecret and AWSRegion (AWSToken is optional for STS)"})
		return
	}

	portScanResult := scanner.Scan(portScanPayload.AWSKey, portScanPayload.AWSSecret, portScanPayload.AWSToken, portScanPayload.AWSRegion)

	c.JSON(http.StatusOK, portScanResult)
}
