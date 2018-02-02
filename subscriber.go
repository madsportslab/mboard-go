package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type SubscriberResponse struct {
	Page 			string			`json:"page"`
}

var subscribers map[*websocket.Conn] bool

func relay(msg string) {

	r := SubscriberResponse{
		Page: msg,
	}

	j, err := json.Marshal(r)

	if err != nil {
		log.Println(err)
		return
	}

	for s, _ := range subscribers {
		s.WriteMessage(websocket.TextMessage, j)
	}

} // relay

func subscriberHandler(w http.ResponseWriter, r *http.Request) {

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

	if subscribers == nil {
		subscribers = make(map[*websocket.Conn]bool)
	}

	subscribers[c] = true

	for {
   
		_, msg, err := c.ReadMessage()

		if err != nil {

			log.Println("[Error] ", err)

			if websocket.IsUnexpectedCloseError(err) {

				//log.Println("Removing connection: ", c)
				delete(subscribers, c)

			}

			break

		}

    if msg == nil {
			log.Println(msg)
			break
		}		

	}
	
} // subscriberHandler
