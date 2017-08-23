package main

import (
	"database/sql"
	"flag"
	"fmt"
	"text/template"
  "log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const (
	APPNAME 					= "mboard-go v%s"
	TEST_ADDRESS 	    = "127.0.0.1:8000"
	CLOUD_ADDRESS     = "madsportslab.com"
	MBOARD            = "mboard"
	VERSION 					= "0.1"
)

const (
	MODE_WIFI				= 0
	MODE_HOTSPOT   	= 1
	MODE_WIRED      = 2
	MODE_CLOUD      = 3
	MODE_TEST   		= 4
)

const (
	INTERFACE_WIFI 		= "en"
	INTERFACE_HOTSPOT	= "wlan"
	INTERFACE_WIRED   = "eth"
	INTERFACE_CLOUD   = "cloud"
	INTERFACE_TEST		= "lo"
	INTERFACE_ERROR   = ""
)

var database 	= flag.String("database", "./data/mboard.db", "database address")
var port 			= flag.String("port", "8000", "service port")
var mode      = flag.Int("mode", MODE_WIFI, "configuration mode")
var ssl       = flag.Bool("ssl", false, "use SSL encryption")
var certFile  = flag.String("cert", "ssl.crt", "SSL certificate")
var keyFile   = flag.String("key", "ssl.key", "SSL private key")

var testTmpl = template.Must(template.ParseFiles("mboard-www/test.html"))

var data *sql.DB = nil

func version() string {
  return fmt.Sprintf(APPNAME, VERSION)
} // version

func testHandler(w http.ResponseWriter, r *http.Request) {
	testTmpl.Execute(w, nil)
} // testHandler

func initDatabase() {

  db, err := sql.Open("sqlite3", *database)

	if err != nil {
		log.Fatal("Database connection error: ", err)
	}

	data = db

} // initDatabase

func initRouter() *mux.Router {

  router := mux.NewRouter()

  router.PathPrefix("/mboard-www/").Handler(http.StripPrefix("/mboard-www/",
    http.FileServer(http.Dir("./mboard-www"))))

	router.HandleFunc("/api/games", gameHandler)
	router.HandleFunc("/api/games/{id:[0-9a-f]+}", gameHandler)
  router.HandleFunc("/api/scores", scoreHandler)
	router.HandleFunc("/api/scores/{id:[0-9a-f]+}", scoreHandler)
	router.HandleFunc("/api/scores/{id:[0-9a-f]+}/logs", logHandler)
	router.HandleFunc("/api/version", versionHandler)

	// management apis

	router.HandleFunc("/api/mgmt/details", detailsHandler)
	router.HandleFunc("/api/mgmt/machine", machineHandler)
	
	router.HandleFunc("/display", displayHandler)
	router.HandleFunc("/test", testHandler)
	router.HandleFunc("/setup", setupHandler)

	//router.HandleFunc("/ws/games/{id:[0-9a-f]+}", controlHandler)
  router.HandleFunc("/ws/game", controlHandler)

  return router

} // initRouter

func main() {

  flag.Parse()
	
	addr, err := getAddress()

	if err != nil {
		log.Fatal(err)
	}
	
  log.Printf("[%s] listening on port %s", version(), *port)

  initDatabase()

  router := initRouter()

	if *ssl {
		log.Fatal(http.ListenAndServeTLS(addr, *certFile, *keyFile, router))
	} else {
		log.Fatal(http.ListenAndServe(addr, router))
	}

} // main
