package main

import (
	"database/sql"
	"encoding/json"
  "log"
	"net/http"

)

type GameRecord struct {
	Home		*Team			`json:"home"`
	Away		*Team			`json:"away"`
}

type GameTbl struct {
	ID			string 	`json:"id"`
  Data    sql.NullString `json:"data"`
	Created string 	`json:"created"`
	Updated string 	`json:"updated"`
	Status  int 		`json:"status"`
}


func addGame() int64 {

	res, err := data.Exec(
		GameCreate,
	)

	if err != nil {
		log.Println(err)
		return -1
	}

	id, err := res.LastInsertId()

	if err != nil {
		
		log.Println(err)
		return -1

	}

	return id

} // addGame

func updateGame(id int64, val string) {

	_, err := data.Exec(
		GameUpdate, val, 1, id,
	)

	if err != nil {
		log.Println(err)
	}

} // updateGame

func getGames() []GameTbl {

  rows, err := data.Query(
		GamesGet,
	)

	if err != nil {
		log.Printf("[%s][Error][DB] %s", version(), err)
		return nil
	}

	defer rows.Close()

	gt := []GameTbl{}

	for rows.Next() {

			g := GameTbl{}

			err := rows.Scan(&g.ID, &g.Data, &g.Status, &g.Created, &g.Updated)

			if err == sql.ErrNoRows || err != nil {
				log.Printf("[%s][Error] %s", version(), err)
				return nil
			}

			gt = append(gt, g)

	}

	return gt

} // getGames

func scoreHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
  case http.MethodPost:
	case http.MethodGet:

		gts := getGames()

		j, jsonErr := json.Marshal(gts)

		if jsonErr != nil {
			log.Printf("[%s] %s", version(), jsonErr)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(j)
		}

	case http.MethodPut:
	case http.MethodDelete:
	default:
		log.Printf("[%s][Error] unsupported command", version())
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // scoreHandler
