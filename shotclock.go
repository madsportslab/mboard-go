package main

import (
	"log"
	"net/http"

	"github.com/eknkc/amber"
)


func shotclockHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodGet:
		
		compiler := amber.New()

		err := compiler.ParseFile("mboard-www/shotclock.amber")

		if err != nil {
			
			log.Printf("[%s][Error] %s", version(), err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		template, err2 := compiler.Compile()

		if err2 != nil {
			
			log.Printf("[%s][Error] %s", version(), err2)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		template.Execute(w, nil)

  case http.MethodPost:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // shotclockHandler
