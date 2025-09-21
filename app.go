package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

type App struct {
	Router *gin.Engine
	Geoip  *geoip2.Reader
	ASN    *geoip2.Reader
}

func (a *App) Initialize() {
	a.Router = gin.New()
	a.Router.Use(gin.Logger())
	a.Router.Use(gin.Recovery())
	// Trust all proxies so ClientIP relies on X-Forwarded-For/X-Real-IP
	_ = a.Router.SetTrustedProxies([]string{"0.0.0.0/0", "::/0"})
	a.InitGeoIP()
	a.initializeRoutes()
}

func (a *App) InitGeoIP() {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	a.Geoip = db

	asnDb, err := geoip2.Open("GeoLite2-ASN.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	a.ASN = asnDb
}

func (a *App) Run(port string) {
	// gin runs the HTTP server internally
	_ = a.Router.Run(port)
}

func (a *App) initializeRoutes() {
	a.Router.GET("/info", a.getClientInfo)
}
