package main

import (
	//"fmt"
	"log"
  "net/http"

	"github.com/eknkc/amber"
	"github.com/gorilla/mux"
	"github.com/madsportslab/glbs"
)

func videoHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
  case http.MethodGet:

		data := make(map[string]string)

		vars := mux.Vars(r)

		id := vars["id"]

		if id == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {

			glbs.SetNamespace("blobs")
			data["video"] = "/" + *glbs.GetPath(id)

			compiler := amber.New()
	
			err := compiler.ParseFile("mboard-www/video.amber")
	
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
	
			template.Execute(w, data)

		}

  case http.MethodPost:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // videoHandler
