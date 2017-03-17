package main

import (
	"flag"
	"fmt"
	"text/template"
  "log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	APPNAME = "madpi v%s"
	VERSION = "0.1"
)

var server = flag.String("server", "127.0.0.1:8000", "http server address")

var testTmpl = template.Must(template.ParseFiles("www/test.html"))

func version() string {
  return fmt.Sprintf(APPNAME, VERSION)
} // version

func testHandler(w http.ResponseWriter, r *http.Request) {
	testTmpl.Execute(w, nil)
} // testHandler

func initRouter() *mux.Router {

  router := mux.NewRouter()

  router.PathPrefix("/www/").Handler(http.StripPrefix("/www/",
    http.FileServer(http.Dir("./www"))))

	router.HandleFunc("/api/games", gameHandler)
	router.HandleFunc("/api/games/{id:[0-9a-f]+}", gameHandler)

	router.HandleFunc("/display", displayHandler)
	router.HandleFunc("/test", testHandler)

	router.HandleFunc("/ws/games/{id:[0-9a-f]+}", controlHandler)

  return router

} // initRouter

func main() {

  flag.Parse()

  log.Printf("[%s] listening on address %s", version(), *server)

  router := initRouter()

	log.Fatal(http.ListenAndServe(*server, router))

} // main
