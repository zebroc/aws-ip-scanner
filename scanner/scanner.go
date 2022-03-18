package scanner

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func ScanPort(ip string, port int, timeout time.Duration) string {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)
	state := ""

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanPort(ip, port, timeout)
		} else {
			state = fmt.Sprintf("%d:closed (%v)", port, err)
		}

		state = fmt.Sprintf("%d:error (%v)", port, err)
		return state
	}

	defer conn.Close()
	state = fmt.Sprintf("%d:open", port)

	return state
}
