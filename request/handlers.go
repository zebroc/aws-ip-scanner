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

	var portScanPayload scanner.PortScanPayload
	errUnmarshal := json.Unmarshal(event.Payload, &portScanPayload)
	if errUnmarshal != nil {

	}

	portScanResult := scanner.Scan(portScanPayload.AWSKey, portScanPayload.AWSSecret, "", portScanPayload.AWSRregion)

	c.JSON(http.StatusOK, portScanResult)
}
