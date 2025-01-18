package handlers

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func getIPfromRequest(r *http.Request) (net.IP, error) {
	// first try to get IPs from headers made by proxy server:
	ipStr := r.Header.Get("X-Real-IP")
	ip := net.ParseIP(ipStr)
	if ip != nil {
		return ip, nil
	}
	// then check for X-Forwarded-For header
	ips := r.Header.Get("X-Forwarded-For")
	ipStrs := strings.Split(ips, ",")
	// parse first ip from list
	ip = net.ParseIP(ipStrs[0])
	if ip != nil {
		return ip, nil
	}
	// then get ip from request.remoteAddr
	addr := r.RemoteAddr
	ipStr, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	ip = net.ParseIP(ipStr)
	if ip != nil {
		return ip, nil
	}
	return nil, fmt.Errorf("failed to parse ip from request.remoteAddr")
}
