package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/oschwald/geoip2-golang"
)

type App struct {
	Router *mux.Router
	Geoip  *geoip2.Reader
	ASN    *geoip2.Reader
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
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
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Content-Length", "Accept-Encoding", "Content-Range", "Content-Disposition", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "OPTIONS"})

	http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/info", a.getClientInfo).Methods("GET")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"status": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
