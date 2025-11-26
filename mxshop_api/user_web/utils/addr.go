package utils

import (
	"net"
)

// GetOutBoundIP 获取对外通信的本机 IP 
// 作用是为了获取Consul的动态IP
func GetOutBoundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}