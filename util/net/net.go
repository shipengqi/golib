package net

import (
	"math"
	"net"
	"net/http"
	"strings"
)

// LocalIP returns the local IP address.
func LocalIP() (net.IP, error) {
	var (
		ip    net.IP
		err   error
		addrs []net.Addr
	)
	addrs, err = net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				ip = ipnet.IP
				break
			}
		}
	}

	return ip, err
}

// IsIPV4 reports whether the given IP is ipv4.
func IsIPV4(ip string) bool {
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false
	}
	return true
}

// ClientIP returns the IP address of client.
func ClientIP(req *http.Request) string {
	var (
		ip  string
		err error
	)
	// "X-Forwarded-For: 10.0.0.1, 10.0.0.2, 10.0.0.3"
	// returns the first global address
	ip = strings.TrimSpace(
		strings.Split(req.Header.Get("X-Forwarded-For"), ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(req.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err = net.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// IPString2Uint converts the given IP string to an uint value
func IPString2Uint(ip string) uint {
	b := net.ParseIP(ip).To4()
	if b == nil {
		return 0
	}

	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24
}

// Uint2IPString converts the given uint value to IP string
func Uint2IPString(i uint) string {
	if i > math.MaxUint32 {
		return ""
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip.String()
}
