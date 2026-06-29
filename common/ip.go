package common

import (
	"net"
	"strings"
)

func IsIP(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil
}

func ParseIP(s string) net.IP {
	return net.ParseIP(s)
}

func IsPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	private := []net.IPNet{
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
	}

	for _, privateNet := range private {
		if privateNet.Contains(ip) {
			return true
		}
	}
	return false
}

// NormalizeIP 规范化 IP 地址
// - 去掉 IPv6 zone identifier (如 %eth0)
// - 将 IPv4-mapped IPv6 (::ffff:x.x.x.x) 转换为纯 IPv4
func NormalizeIP(ipStr string) string {
	if ipStr == "" {
		return ""
	}
	// 去掉 zone identifier
	if idx := strings.Index(ipStr, "%"); idx != -1 {
		ipStr = ipStr[:idx]
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ipStr // 解析失败则原样返回
	}
	// IPv4-mapped IPv6 转换为 IPv4
	if ipv4 := ip.To4(); ipv4 != nil {
		return ipv4.String()
	}
	return ip.String()
}

func IsIpInCIDRList(ip net.IP, cidrList []string) bool {
	for _, cidr := range cidrList {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			// 尝试作为单个IP处理
			if whitelistIP := net.ParseIP(cidr); whitelistIP != nil {
				if ip.Equal(whitelistIP) {
					return true
				}
			}
			continue
		}

		if network.Contains(ip) {
			return true
		}
	}
	return false
}
