package main

import (
	"encoding/json"
	"testing"

	"github.com/gorilla/websocket"
)

var invalidURL = []string{
	"127.0.0.1:8000/games",
	"ws://127.0.0.1:8000/games",
	"ws://127.0.0.1:8000/games/",
	"ws://127.0.0.1:8000/api/games/",
	"ws://127.0.0.1:8000/api/games/a",
}

var validURL = []string {
	"ws://127.0.0.1:8000/api/games",
	"ws://127.0.0.1:8000/api/games/0",
}

func TestInvalidConnect(t *testing.T) {

  for _, u := range invalidURL {

    _, _, err := websocket.DefaultDialer.Dial(u, nil)

		if err == nil {
			t.Fatal("Should not connect successfully.")
		}
		
	}


} // TestInvalidConnect

func TestConnect(t *testing.T) {

  for _, u := range validURL {

    ws, _, err := websocket.DefaultDialer.Dial(u, nil)

		if err != nil {
			t.Fatal(err)
		}

		defer ws.Close()
	
	}

} // TestConnect

func TestInit(t *testing.T) {

  ws, _, err := websocket.DefaultDialer.Dial(
		"ws://127.0.0.1:8000/api/games", nil)
	
	if err != nil {
		t.Fatal(err)
	}

	defer ws.Close()

  req := Req{
		Cmd: WS_INIT,
		Periods: 4,
	}

  j, marshalErr := json.Marshal(req)

  if marshalErr != nil {
		t.Error(marshalErr)
	}

	writeErr := ws.WriteMessage(websocket.TextMessage, j)

  if writeErr != nil {
		t.Fatal(writeErr)
	}

} // TestInit
