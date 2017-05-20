package main

import (
	"log"
  "net/http"

	"github.com/eknkc/amber"
)

func displayHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
  case http.MethodGet:

		data := make(map[string]string)

		data["base"] = getAddress("en0")

    log.Println(data)
	  compiler := amber.New()

		err := compiler.ParseFile("www/display.amber")

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

		log.Println(data)

		template.Execute(w, data)

  case http.MethodPost:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // displayHandler
