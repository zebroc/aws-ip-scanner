package main

import (
	"github.com/zebroc/aws-ip-scanner/request"
)

// @title AWS Scanner
// @version 1.0
// @description A scanner
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /
func main() {
	request.StartWebServer()
}
