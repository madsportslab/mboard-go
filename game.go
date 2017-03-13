package main

import (
	"encoding/json"
	"log"
  "net/http"

	"github.com/gorilla/websocket"
)

const (
	WS_INIT							= "INIT"
	WS_GAME_START				= "GAME_START"
	WS_GAME_STOP				= "GAME_STOP"
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

type Req struct {
	Cmd					string 	`json:"cmd"`
	Periods			int 		`json:"periods"`
  Minutes			int 		`json:"minutes"`
	Shot			  int 		`json:"shot"`
	Timeouts		int 		`json:"timeouts"`
	Fouls				int 		`json:"fouls"`
}

func gameHandler(w http.ResponseWriter, r *http.Request) {

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
		case WS_INIT:

		  log.Println(req)

		case WS_GAME_START:
		case WS_GAME_STOP:
		default:
		  log.Println("[Error] unsupported command:", string(msg))
		}

	}

} // gameHandler
