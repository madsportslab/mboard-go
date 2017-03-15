package main

import (
	//"encoding/json"
	"net/http"
	"net/url"
	"strings"
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

		ws.Close()
	
	}

} // TestConnect

func TestInit(t *testing.T) {

  form := url.Values{}
	form.Add("periods", "19")
	form.Add("minutes", "12")

  r, postErr := http.Post("http://127.0.0.1:8000/api/games",
	  "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	
	if postErr != nil {
		t.Fatal(postErr)
	}

	t.Log(r)

	
} // TestInit
