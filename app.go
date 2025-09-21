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
	// Minimal CORS: allow all origins and handle preflight
	a.Router.Use(corsMiddleware())
	// Trust all proxies so ClientIP relies on X-Forwarded-For/X-Real-IP
	_ = a.Router.SetTrustedProxies([]string{"0.0.0.0/0", "::/0"})
	a.InitGeoIP()
	a.initializeRoutes()
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Content-Length, Accept-Encoding, Content-Range, Content-Disposition, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (a *App) InitGeoIP() {
	db, err := geoip2.Open("/opt/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	a.Geoip = db

	asnDb, err := geoip2.Open("/opt/GeoLite2-ASN.mmdb")
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
