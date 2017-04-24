package main

import (
	"database/sql"
	"log"
)

type Log struct {
	ID          string  `json:"id"`
	Data				string	`json:"data"`
	Created			string	`json:"created"`
	Updated			string	`json:"updated"`
}

func add(gid string, msg string) {

  _, err := data.Exec(
		LogCreate, gid, msg,
	)

	if err != nil {
		
		log.Printf("[%s][Error] %s", version(), err)
		return

	}

} // add

func get(gid string) []Log {

	rows, err := data.Query(
		LogGet, gid,
	)

	if err != nil {

		log.Printf("[%s][Error] %s", version(), err)
		return nil

	}

	defer rows.Close()

	logs := []Log{}

	for rows.Next() {

		l := Log{}

		err := rows.Scan(&l.ID, &l.Data, &l.Created, &l.Updated)

		if err == sql.ErrNoRows || err != nil {
			
			log.Printf("[%s][Error] %s", version(), err)
			return nil

		}

		logs = append(logs, l)

	}

	return logs

} // get

