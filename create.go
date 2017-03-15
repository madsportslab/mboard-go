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

					switch f {
					case PERIODS:
					  config.Periods = int(i)
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

    x := struct {GameId string}{GameId: gid,}
		
		j, err := json.Marshal(x)

		if err != nil {
			log.Println(err)
		}

		w.Write(j)

	case http.MethodGet:
	case http.MethodPut:
	case http.MethodDelete:
	default:
		log.Printf("[%s][Error] unsupported command", version())
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // gameHandler
