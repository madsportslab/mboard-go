package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
  "log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const (
	HOME			= "home"
	AWAY			= "away"
	PERIODS   = "periods"
	MINUTES   = "minutes"
	FOULS     = "fouls"
	TIMEOUTS  = "timeouts"
	SHOT      = "shot"
)

type Config struct {
  Periods			int 		`json:"periods"`
  Minutes			int 		`json:"minutes"`
	Shot			  int 		`json:"shot"`
	Timeouts		int 		`json:"timeouts"`
	Fouls				int 		`json:"fouls"`
	Home        string  `json:"home"`
	Away        string  `json:"away"`
}

type GameInfo struct {
  Settings			*Config
	Game			    *Game
}

type PostRes struct {
	GameId 	string 	`json:"gameId"`
}

var games = make(map[string]*GameInfo)

func parseConfig(r *http.Request) *Config {

    config := Config{
			Periods: 		4,
			Minutes:		12,
			Shot:				30,
			Timeouts:		3,
			Fouls:			10,
			Home:				HOME,
			Away:       AWAY,
		}

		fields := []string{HOME, AWAY, PERIODS, MINUTES, FOULS, TIMEOUTS, SHOT}

		for _, f := range fields {

			val := r.FormValue(f)

			if f == HOME || f == AWAY {
				// string value

			} else {

				if val == "" {
					continue
				}

				i, err := strconv.ParseInt(val, 0, 8)

				if err != nil {
					log.Println(err)
				} else {

          if i < 1 || i > 30 {
						continue
					}

					switch f {
					case PERIODS:

					  if i == 2 || i == 4 {
					    config.Periods = int(i)
						}

					case MINUTES:
					  config.Minutes = int(i)
					case FOULS:
					  config.Fouls = int(i)
					case TIMEOUTS:
					  config.Timeouts = int(i)
					case SHOT:
						config.Shot = int(i)
					}
					
				}

			}

		}

		return &config

		
} // parseConfig

func generateId(config *Config, length int) string {

  now := time.Now().String()

  digest := hmac.New(sha256.New, []byte("ABC"))

	digest.Write([]byte(fmt.Sprintf("%s%s", now, config)))

	hash := hex.EncodeToString(digest.Sum(nil))

	return hash[:length]
 
} // generateId

func gameHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
  case http.MethodPost:

		log.Printf("[%s] POST /games", version())

    config := parseConfig(r)

		log.Println(config)

		gid := generateId(config, 10)

		log.Println(gid)

    ps := PostRes{
			GameId: gid,
		}
		
		j, err := json.Marshal(ps)

		if err != nil {
			log.Println(err)
		}

    gi := GameInfo{
			Settings:	config,
			Game: nil,
		}

    games[gid] = &gi

		log.Println(games)

		w.Write(j)

	case http.MethodGet:

		vars := mux.Vars(r)

		id := vars["id"]

		log.Printf("[%s] GET /games/%s", version(), id)

		gameInfo := games[id]

		if gameInfo == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {

			j, jsonErr := json.Marshal(gameInfo.Settings)

			if jsonErr != nil {
				log.Println(jsonErr)
			}

			w.Write(j)

		}

	  
	case http.MethodPut:
	case http.MethodDelete:
	default:
		log.Printf("[%s][Error] unsupported command", version())
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // gameHandler
