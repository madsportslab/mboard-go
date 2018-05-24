package main

import (
	"fmt"
	"log"
  "net/http"

	"github.com/eknkc/amber"
	"github.com/gorilla/mux"
)

const (
	VIEW_CLOCK        = "clock"
	VIEW_SCORE        = "score"
	VIEW_SCOREBOARD		= "scoreboard"
	VIEW_SHOTCLOCK    = "shotclock"
)

func checkView(view string) bool {

	if view == VIEW_CLOCK || view == VIEW_SCORE || view == VIEW_SCOREBOARD ||
	  view == VIEW_SHOTCLOCK {
			return true
		} else {
			return false
		}


} // checkView

func themeHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
  case http.MethodGet:

	  compiler := amber.New()

		vars := mux.Vars(r)

		theme	:= vars["theme"]
		view 	:= vars["view"]

		if checkView(view) {

			err := compiler.ParseFile(fmt.Sprintf(
				"mboard-www/theme.%s.%s.amber", theme, view))
	
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

		} else {
			w.WriteHeader(http.StatusNotFound)
		}


	
  case http.MethodPost:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // themeHandler
