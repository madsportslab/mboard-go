package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	"ws://127.0.0.1:8000/ws",
	"http://127.0.0.1:8000/ws/games",
	"ws://127.0.0.1:8000/ws/games/g",
}

var validURL = []string {
	"ws://127.0.0.1:8000/ws/games/a",
	"ws://127.0.0.1:8000/ws/games/0",
}

func parseBody(body io.ReadCloser, s interface {}) {

  j, readErr := ioutil.ReadAll(body)

  if readErr != nil {
		log.Println(readErr)
	}

  json.Unmarshal(j, &s)

} // parseBody

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

func TestNewGame(t *testing.T) {

  form := url.Values{}
	form.Add("periods", "19")
	form.Add("minutes", "12")

  r, postErr := http.Post("http://127.0.0.1:8000/api/games",
	  "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	
	if postErr != nil {
		t.Fatal(postErr)
	}

	defer r.Body.Close()

  b := PostRes{}

  j, readErr := ioutil.ReadAll(r.Body)

  if readErr != nil {
		t.Fatal(readErr)
	}

  json.Unmarshal(j, &b)
	
	if b.GameId == "" {
		t.Fatal("No GameId returned.")
	}

  url := fmt.Sprintf("%s%s",
	  "http://127.0.0.1:8000/api/games/", b.GameId)
	
	r2, getErr := http.Get(url)

	if getErr != nil {
		t.Fatal(getErr)
	}

  config := Config{}

	parseBody(r2.Body, &config)

  if config.Periods != 4 {
		t.Fatal("Returned incorrect periods")
	}

	if config.Minutes != 12 {
		t.Fatal("Returned incorrect minutes")
	}
	
} // TestNewGame
