package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Log struct {
	ID          string  `json:"id"`
	Msg					string	`json:"msg"`
	Created			string	`json:"created"`
	Updated			string	`json:"updated"`
}

const (

	LogCreate = "INSERT into logs" +
	  "(game_id, msg) " + 
		"VALUES ($1, $2)"
	
	LogGet = "SELECT " +
	  "id, msg, created, updated " +
		"FROM logs " + 
		"WHERE game_id=? ORDER BY created DESC"

  LogDelete = "DELETE from logs WHERE id=?"

)

func put(game_id string, req Req) {

  j, errJson := json.Marshal(req)

	if errJson != nil {
		log.Println(errJson)
		return
	}

  _, err := data.Exec(
		LogCreate, game_id, j,
	)

	if err != nil {
		
		log.Printf("[%s][Error] %s", version(), err)
		return

	}

} // put

func get(game_id string) []Log {

	rows, err := data.Query(
		LogGet, game_id,
	)

	if err != nil {

		log.Printf("[%s][Error] %s", version(), err)
		return nil

	}

	defer rows.Close()

	logs := []Log{}

	for rows.Next() {

		l := Log{}

		err := rows.Scan(&l.ID, &l.Msg, &l.Created, &l.Updated)

		if err == sql.ErrNoRows || err != nil {
			
			log.Printf("[%s][Error] %s", version(), err)
			return nil

		}

		logs = append(logs, l)

	}

	return logs

} // get

func delete(log_id string) bool {

  _, err := data.Exec(
		LogDelete, log_id,
	)

	if err != nil {
		
		log.Printf("[%s][Error] %s", version(), err)
		return false

	}

	return true

} // delete


func logHandler(w http.ResponseWriter, r *http.Request) {

	mux := mux.Vars(r)

  switch r.Method {
  case http.MethodGet:

	  id := mux["id"]

		logs := get(id)

		j, err := json.Marshal(logs)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Write(j)
		}

  case http.MethodPost:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // logHandler