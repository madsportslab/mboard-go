package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
  "log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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
	Conns 				map[*websocket.Conn]*sync.Mutex
	Final         bool
	Created       string
	ID            int64
	Active        bool
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
	Msg 	string 	`json:"msg"`
}

type GameTbl struct {
	ID			string `json:"id"`
  Data    sql.NullString `json:"data"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

var game = &GameInfo{}

func parseConfig(r *http.Request) *Config {

    config := Config{
			Periods: 		4,
			Minutes:		1,
			Shot:				24,
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
				config.Home = val
			} else if f == AWAY {
			  config.Away = val
			} else {

				if val == "" {
					continue
				}

				i, err := strconv.ParseInt(val, 0, 8)

				if err != nil {
					log.Println(err)
				} else {

					// TODO: fouls and shot can equal 0
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

			err := rows.Scan(&g.ID, &g.Data, &g.Created, &g.Updated)

			if err == sql.ErrNoRows || err != nil {
				log.Printf("[%s][Error] %s", version(), err)
				return nil
			}

			gt = append(gt, g)

	}

	return gt

} // getGames


func gameHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
  case http.MethodPost:

		log.Printf("[%s] POST /games", version())

    config := parseConfig(r)

		log.Println(config)
		
		if game.Active {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		id := addGame()

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
			Conns: make(map[*websocket.Conn]*sync.Mutex),
			Final: false,
			Created: time.Now().String(),
			ID: id,
			Active: true,
		}

    game = &gi

	case http.MethodGet:

		log.Println(game)

		if game.Active {

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

		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	case http.MethodPut:
	  // save game

	case http.MethodDelete:
	default:
		log.Printf("[%s][Error] unsupported command", version())
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // gameHandler
