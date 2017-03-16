package main

import (
	"flag"
	"fmt"
  "log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	APPNAME = "madpi v%s"
	VERSION = "0.1"
)

var server = flag.String("server", "127.0.0.1:8000", "http server address")

func version() string {
  return fmt.Sprintf(APPNAME, VERSION)
} // version

func initRouter() *mux.Router {

  router := mux.NewRouter()

	router.HandleFunc("/api/games", gameHandler)
	router.HandleFunc("/api/games/{id:[0-9a-f]+}", gameHandler)
	
	router.HandleFunc("/ws/games/{id:[0-9a-f]+}", socketHandler)

  return router

} // initRouter

func main() {

  flag.Parse()

  log.Printf("[%s] listening on address %s", version(), *server)

  router := initRouter()

	log.Fatal(http.ListenAndServe(*server, router))

} // main
