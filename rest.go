package main

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
		"::/128",
		"ff00::/8",
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err == nil {
			privateIPBlocks = append(privateIPBlocks, block)
		}
	}
}

func isPublicIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsUnspecified() {
		return false
	}
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return false
		}
	}
	return true
}

func getExternalIP(c *gin.Context) string {
	// Prefer original client from X-Forwarded-For (left-most public IP)
	if xff := strings.TrimSpace(c.GetHeader("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		for _, p := range parts {
			candidate := net.ParseIP(strings.TrimSpace(p))
			if isPublicIP(candidate) {
				return candidate.String()
			}
		}
	}
	// Next, X-Real-Ip if public
	if xri := strings.TrimSpace(c.GetHeader("X-Real-Ip")); xri != "" {
		if ip := net.ParseIP(xri); isPublicIP(ip) {
			return ip.String()
		}
	}
	// Fallback to ClientIP (may be proxy if headers absent)
	if ip := net.ParseIP(c.ClientIP()); ip != nil {
		if isPublicIP(ip) {
			return ip.String()
		}
		return ip.String()
	}
	return c.ClientIP()
}

type info struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Code    string `json:"code"`
	City    string `json:"city"`
	ISP     string `json:"isp"`
	ISPCode uint   `json:"isp_code"`
}

func (a *App) getClientInfo(c *gin.Context) {
	var i info

	// optional query param ?ip=1.2.3.4; fallback to client's IP
	ipParam := strings.TrimSpace(c.Query("ip"))
	if ipParam != "" {
		if net.ParseIP(ipParam) == nil {
			c.JSON(400, gin.H{"status": "invalid ip"})
			return
		}
		i.IP = ipParam
	} else {
		i.IP = getExternalIP(c)
	}
	record, err := a.Geoip.City(net.ParseIP(i.IP))
	if err != nil {
		c.JSON(404, gin.H{"status": "GeoIP not found"})
		return
	}

	i.Country = record.Country.Names["en"]
	i.Code = record.Country.IsoCode
	i.City = record.City.Names["en"]

	asnRecord, err := a.ASN.ASN(net.ParseIP(i.IP))
	if err == nil {
		i.ISP = asnRecord.AutonomousSystemOrganization
		i.ISPCode = asnRecord.AutonomousSystemNumber
	}

	c.JSON(200, i)
}
