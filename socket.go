package main

import (
	"encoding/json"
	"log"
  "net/http"

	"github.com/gorilla/websocket"
)

const (
	WS_CLOCK_START			= "CLOCK_START"
	WS_CLOCK_STOP				= "CLOCK_STOP"
	WS_SHOT_START				= "SHOT_START"
	WS_SHOT_STOP				= "SHOT_STOP"
	WS_SCORE_UP         = "SCORE_UP"
	WS_SCORE_DOWN       = "SCORE_DOWN"
	WS_FOUL_UP         	= "FOUL_UP"
	WS_FOUL_DOWN       	= "FOUL_DOWN"
	WS_TIMEOUT_UP       = "TIMEOUT_UP"
	WS_TIMEOUT_DOWN     = "TIMEOUT_DOWN"
	WS_PERIOD_UP        = "PERIOD_UP"
	WS_PERIOD_DOWN      = "PERIOD_DOWN"
	WS_FINAL            = "FINAL" 
)

type Point struct {
  Total			int			`json:"total"`
	Period		int			`json:"period"`
}

type Team struct {
	Name      string    `json:"name"`
	Logo      string    `json:"logo"`
	Fouls			int				`json:"fouls"`
	Timeouts  int     	`json:"timeouts"`
	Points    []Point   `json:"points"`
}

type Game struct {
	Home				*Team		`json:"home"`
	Away      	*Team		`json:"away"`
	Period    	int			`json:"period"`
	Clock				string	`json:"clock"`
	Possesion 	bool		`json:"possesion"`
}

type ReqData struct {
	Value				string
	Step				int
}

type Req struct {
	Cmd					string 		`json:"cmd"`
	Data        ReqData  	`json:"d"`        
}

func socketHandler(w http.ResponseWriter, r *http.Request) {

  upgrader := websocket.Upgrader {
		ReadBufferSize:		1024,
		WriteBufferSize: 	1024,
		CheckOrigin:		func(r *http.Request) bool { return true },
	}

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		
		log.Println(err)
		return

	}

	for {

		_, msg, err := c.ReadMessage()

		if err != nil || msg == nil {

			log.Println("[Error] ", err)
			return

		}

    defer c.Close()

    req := Req{}

		json.Unmarshal(msg, &req)

		switch req.Cmd {
		case WS_CLOCK_START:
		case WS_CLOCK_STOP:
		case WS_SHOT_START:
		case WS_SHOT_STOP:
		case WS_SCORE_UP:
		case WS_SCORE_DOWN:
		case WS_FOUL_UP:
		case WS_FOUL_DOWN:
		case WS_TIMEOUT_UP:
		case WS_TIMEOUT_DOWN:
		case WS_PERIOD_UP:
		case WS_PERIOD_DOWN:
		case WS_FINAL:
		default:
		  log.Printf("[%s][Error] unsupported command: %s", version(), string(msg))
		}

	}

} // socketHandler
