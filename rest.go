package main

import (
	"net"
	"net/http"
	"strings"
)

type info struct {
	IP string                 `json:"ip"`
	Country string 			  `json:"country"`
	Code string 			  `json:"code"`
	City string 			  `json:"city"`
}


func (a *App) getClientInfo(w http.ResponseWriter, r *http.Request) {

	var i info

	i.IP = getRealIP(r)
	record, err := a.Geoip.City(net.ParseIP(i.IP))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "GeoIP not found")
		return
	}

	i.Country = record.Country.Names["en"]
	i.Code = record.Country.IsoCode
	i.City = record.City.Names["en"]

	respondWithJSON(w, http.StatusOK, i)
}

func getRealIP(r *http.Request) string {

	remoteIP := ""
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
		// parse X-Real-Ip header
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}

	return remoteIP
}
