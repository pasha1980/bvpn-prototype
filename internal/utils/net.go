package utils

import (
	"net"
)

func GetIp(url string) string {
	ips, _ := net.LookupIP(url)
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}

	return ""
}

func MyIP() string {
	out, _ := Exec("curl ifconfig.me")
	return out
}
