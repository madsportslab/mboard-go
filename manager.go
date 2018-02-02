package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	WS_LOGIN						=	"LOGIN"
	WS_SCOREBOARD				= "SCOREBOARD"
	WS_ADVERTISEMENT		= "ADVERTISEMENT"
	WS_VIDEO						= "VIDEO"
	WS_PHOTO						= "PHOTO"
	WS_LOGO             = "LOGO"
)

var manager *websocket.Conn

type ManagerCommand struct {
  Cmd 			string								`json:"cmd"`
	Options		map[string]string			`json:"options"`
}

func managerHandler(w http.ResponseWriter, r *http.Request) {

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

	manager = c

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

    mc := ManagerCommand{}

		json.Unmarshal(msg, &mc)
		
		log.Println(mc)

		switch mc.Cmd {
		case WS_LOGO:
			relay(WS_LOGO)

		case WS_LOGIN:
			log.Println(WS_LOGIN)

		case WS_SCOREBOARD:
			relay(WS_SCOREBOARD)

		case WS_ADVERTISEMENT:
			relay(WS_ADVERTISEMENT)

		case WS_VIDEO:
			relay(WS_VIDEO)

		case WS_PHOTO:
			relay(WS_PHOTO)

		default:
			log.Println("unknown")
		}

	}

} // managerHandler
