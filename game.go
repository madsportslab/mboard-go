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

	//"github.com/gorilla/mux"
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
	GameData			*Game
}

type GameState struct {
	Settings      *Config   `json:"settings"`
	Period        int				`json:"period"`
	Possession    bool			`json:"possession"`
	Home          *Team			`json:"home"`
	Away          *Team			`json:"away"`
	GameClock     *Clock    `json:"game"`
	ShotClock     *Clock    `json:"shot"`
}

type GameRes struct {
	GameId 	string 	`json:"gameId"`
}

var game *GameInfo

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

func initTeam(name string, timeouts int) *Team {

  team := Team{
		Name: name,
		Points: make(map[int]int),
		Timeouts: timeouts,
	}

  return &team

} // initTeam

func initGameClocks() *GameClocks {

  gc := GameClocks{
		ShotViolationChan: make(chan bool),
		FinalChan: make(chan bool),
		OutChan: make(chan []byte),
		PlayClock: &Clock{Tenths: 0, Seconds: 0},
		ShotClock: &Clock{Tenths: 0, Seconds: 0},
	}

	return &gc

} // initGameClocks

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
/*
		gid := generateId(config, 10)

		log.Println(gid)

    gr := GameRes{
			GameId: gid,
		}
		
		j, err := json.Marshal(gr)

		if err != nil {
			log.Println(err)
		}
*/
		h := initTeam(config.Home, config.Timeouts)
		a := initTeam(config.Away, config.Timeouts)

		c := initGameClocks()

    gi := GameInfo{
			Settings:	config,
			GameData: &Game{
				Home: h,
				Away: a,
				Clk: c,
			},
		}

    game = &gi

		w.Write([]byte("SUCCESS"))

	case http.MethodGet:

		if game != nil {

			gs := GameState{
				Settings: game.Settings,
				Period: game.GameData.Period,
				Possession: game.GameData.Possession,
				Home: game.GameData.Home,
				Away: game.GameData.Away,
				GameClock: game.GameData.Clk.PlayClock,
				ShotClock: game.GameData.Clk.ShotClock,
			}

			j, jsonErr := json.Marshal(gs)

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
