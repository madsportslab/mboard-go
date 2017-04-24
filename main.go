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
	APPNAME = "madpi v%s"
	VERSION = "0.1"
)

var database 	= flag.String("database", "./db/med.db", "database address")
//var server 		= flag.String("server", "127.0.0.1:8000", "http server address")
var port 			= flag.String("port", ":8000", "service port")

var testTmpl = template.Must(template.ParseFiles("www/test.html"))

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

  router.PathPrefix("/www/").Handler(http.StripPrefix("/www/",
    http.FileServer(http.Dir("./www"))))

	router.HandleFunc("/api/games", gameHandler)
	router.HandleFunc("/api/games/{id:[0-9a-f]+}", gameHandler)
  router.HandleFunc("/api/scores", scoreHandler)
	router.HandleFunc("/api/version", versionHandler)

	router.HandleFunc("/display", displayHandler)
	router.HandleFunc("/test", testHandler)
	router.HandleFunc("/setup", setupHandler)

	//router.HandleFunc("/ws/games/{id:[0-9a-f]+}", controlHandler)
  router.HandleFunc("/ws/game", controlHandler)

  return router

} // initRouter

func main() {

  flag.Parse()

	addr := getAddress("en0")

  log.Printf("[%s] listening on address %s", version(), addr)

  initDatabase()

  router := initRouter()

	log.Fatal(http.ListenAndServe(addr, router))

} // main
