package controller

import (
	"net/http"
	"strings"
)

func GetIPAddr(r *http.Request) string {
	var ipAddress string

	ipAddress = r.RemoteAddr

	if ipAddress != "" {
		ipAddress = strings.Split(ipAddress, ":")[0]
	}

	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			if ip != "" {
				ipAddress = ip
			}
		}
	}
	return ipAddress
}
