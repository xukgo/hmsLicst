package netUtil

import (
	"net"
)

const IPV4_LOCALHOST = "127.0.0.1"

func GetNetIpList() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var ipv4List = make([]net.IP, 0, len(addrs))
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() != nil {
			ipv4List = append(ipv4List, ipnet.IP.To4())
		}
	}
	return ipv4List, nil
}

//理论上的localhost地址
func CheckIsLocalhost(ip net.IP) bool {
	return ip.IsLoopback()
}

//理论上的私网地址
func CheckIsPrivateIP(ip net.IP) bool {
	blen := len(ip)
	if ip[blen-4] == 10 {
		return true
	}
	if ip[blen-4] == 172 && ip[blen-3] >= 16 && ip[blen-3] <= 31 {
		return true
	}
	if ip[blen-4] == 192 && ip[blen-3] == 168 {
		return true
	}
	return false
}

func CheckIsPublicIP(ip net.IP) bool {
	match := ip.IsGlobalUnicast()
	if !match {
		return false
	}
	match = CheckIsPrivateIP(ip)
	if match {
		return false
	}
	return true
}

func GetIPList() ([]net.IP, error) {
	var resIpArr = make([]net.IP, 0, 2)

	ipv4List, err := GetNetIpList()
	if err != nil {
		return nil, err
	}

	for i := range ipv4List {
		if CheckIsPrivateIP(ipv4List[i]) || CheckIsPublicIP(ipv4List[i]) {
			resIpArr = append(resIpArr, ipv4List[i])
		}
	}
	return resIpArr, nil
}

func GetPrivateIPList() ([]net.IP, error) {
	var resIpArr = make([]net.IP, 0, 2)

	ipv4List, err := GetNetIpList()
	if err != nil {
		return nil, err
	}

	for i := range ipv4List {
		if CheckIsPrivateIP(ipv4List[i]) {
			resIpArr = append(resIpArr, ipv4List[i])
		}
	}
	return resIpArr, nil
}

func GetPublicIPList() ([]net.IP, error) {
	var resIpArr = make([]net.IP, 0, 1)

	ipv4List, err := GetNetIpList()
	if err != nil {
		return nil, err
	}

	for i := range ipv4List {
		if CheckIsPublicIP(ipv4List[i]) {
			resIpArr = append(resIpArr, ipv4List[i])
		}
	}
	return resIpArr, nil
}
