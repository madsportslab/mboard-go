package main

import (
	"encoding/json"
	"log"
  "net/http"

  "github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	WS_CLOCK_START			= "CLOCK_START"
	WS_CLOCK_STOP				= "CLOCK_STOP"
	WS_CLOCK_RESET      = "CLOCK_RESET"
	WS_SHOT_RESET       = "SHOT_RESET"
	WS_PERIOD_UP        = "PERIOD_UP"
	WS_PERIOD_DOWN      = "PERIOD_DOWN"
	WS_POSSESSION_HOME  = "POSSESSION_HOME"
	WS_POSSESSION_AWAY  = "POSSESSION_AWAY"
	WS_FINAL            = "FINAL" 
)

const (
	WS_SCORE_HOME       			= "SCORE_HOME"
	WS_SCORE_AWAY       			= "SCORE_AWAY"
	WS_FOUL_HOME_UP     			= "FOUL_HOME_UP"
	WS_FOUL_HOME_DOWN   			=	"FOUL_HOME_DOWN"
	WS_FOUL_AWAY_UP     			= "FOUL_AWAY_UP"
	WS_FOUL_AWAY_DOWN   			= "FOUL_AWAY_DOWN"
	WS_TIMEOUT_HOME_UP  			= "TIMEOUT_HOME_UP"
	WS_TIMEOUT_HOME_DOWN     	= "TIMEOUT_HOME_DOWN"
	WS_TIMEOUT_AWAY_UP       	= "TIMEOUT_AWAY_UP"
	WS_TIMEOUT_AWAY_DOWN     	= "TIMEOUT_AWAY_DOWN"
)

type Team struct {
	Name      string    		`json:"name"`
	Logo      string    		`json:"logo"`
	Fouls			int						`json:"fouls"`
	Timeouts  int     			`json:"timeouts"`
	Points    map[int]int   `json:"points"`
}

type Game struct {
	Home				*Team					`json:"home"`
	Away      	*Team					`json:"away"`
	Period    	int						`json:"period"`
	Clk					*GameClocks		`json:"clk"`
	Possession 	bool					`json:"possesion"`
}

type Req struct {
	Cmd					string 			`json:"cmd"`
	Step				int					`json:"step"`        
}

var control *websocket.Conn

func incrementPoints(id string, name string, val int) {

  if games[id] == nil {
		return
	}

  if name == HOME {

		total := games[id].GameData.Home.Points[games[id].GameData.Period]
		
		games[id].GameData.Home.Points[games[id].GameData.Period] = total +
			val
		
	} else if name == AWAY {

		total := games[id].GameData.Away.Points[games[id].GameData.Period]

		games[id].GameData.Away.Points[games[id].GameData.Period] = total
		
		games[id].GameData.Away.Points[games[id].GameData.Period] = total +
			val
		
	}

} // incrementPoints

func incrementFoul(id string, name string, val int) {

  if games[id] == nil {
		return
	}

  if name == HOME {

		games[id].GameData.Home.Fouls = games[id].GameData.Home.Fouls + val
		
	} else if name == AWAY {

		games[id].GameData.Away.Fouls = games[id].GameData.Away.Fouls + val

	}

} // incrementFoul

func incrementTimeout(id string, name string, val int) {

  if games[id] == nil {
		return
	}

  if name == HOME {

		games[id].GameData.Home.Fouls = games[id].GameData.Home.Timeouts + val
		
	} else if name == AWAY {

		games[id].GameData.Away.Fouls = games[id].GameData.Away.Timeouts + val

	}

} // incrementTimeout


func incrementPeriod(id string, val int) {

  games[id].GameData.Period = games[id].GameData.Period + val

} // incrementPeriod

func setPossession(id string, name string) {

  if name == HOME {
 	 	games[id].GameData.Possession = true
	} else {
		games[id].GameData.Possession = false
	}

} // setPossession

func hose(id string, c *websocket.Conn) {

  for {

		select {
		case s := <-games[id].GameData.Clk.OutChan:
		  log.Println(string(s))

		}

	}

} // hose

func controlHandler(w http.ResponseWriter, r *http.Request) {

  vars := mux.Vars(r)

	id := vars["id"]

  log.Println("fuck:", id)

  upgrader := websocket.Upgrader {
		ReadBufferSize:		1024,
		WriteBufferSize: 	1024,
		CheckOrigin:		func(r *http.Request) bool { return true },
	}

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		
		log.Println("[Error]", err)
		return

	}

	defer c.Close()

	for {
   
		_, msg, err := c.ReadMessage()

		if err != nil {

			log.Println("[Error] ", err)
			break

		}

    if msg == nil {
			log.Println(msg)
			break
		}

    req := Req{}

		json.Unmarshal(msg, &req)

		log.Println(req.Cmd)

		switch req.Cmd {
		case WS_CLOCK_START:
			go games[id].GameData.Clk.Start(games[id].Settings)

			go hose(id, c)

		case WS_CLOCK_STOP:
		  games[id].GameData.Clk.Stop()

		case WS_CLOCK_RESET:
		case WS_SHOT_RESET:
		case WS_PERIOD_UP:
			incrementPeriod(id, 1)
		
		case WS_PERIOD_DOWN:
		  incrementPeriod(id, -1)

		case WS_POSSESSION_HOME:
			setPossession(id, HOME)

		case WS_POSSESSION_AWAY:
		  setPossession(id, AWAY)

		case WS_FINAL:
		
		case WS_SCORE_HOME:

		  log.Println(req.Step)

      incrementPoints(id, HOME, req.Step)

    case WS_SCORE_AWAY:

		  log.Println(req.Step)

      incrementPoints(id, AWAY, req.Step)

		case WS_FOUL_HOME_UP:
		  incrementFoul(id, HOME, 1)

		case WS_FOUL_HOME_DOWN:
			incrementFoul(id, HOME, -1)
		
		case WS_FOUL_AWAY_UP:
		  incrementFoul(id, AWAY, 1)

		case WS_FOUL_AWAY_DOWN:
			incrementFoul(id, AWAY, -1)

		case WS_TIMEOUT_HOME_UP:
			incrementTimeout(id, HOME, 1)

		case WS_TIMEOUT_HOME_DOWN:
			incrementTimeout(id, HOME, -1)

		case WS_TIMEOUT_AWAY_UP:
			incrementTimeout(id, AWAY, 1)

		case WS_TIMEOUT_AWAY_DOWN:
			incrementTimeout(id, AWAY, -1)
		
		default:
		  log.Printf("[%s][Error] unsupported command: %s", version(), string(msg))
		}

	}

} // controlHandler
