package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	WS_LOGIN						=	"LOGIN"
	WS_LOGO             = "LOGO"
	WS_CLOCKONLY        = "CLOCKONLY"
	WS_SHOTCLOCK        = "SHOTCLOCK"
  WS_SETUP            = "SETUP"
	WS_SCOREBOARD				= "SCOREBOARD"
	WS_THEME            = "THEME"
	WS_THEME_CURRENT    = "THEME_CURRENT"
	WS_ADVERTISEMENT		= "ADVERTISEMENT"
	WS_VIDEO_PLAY       = "VIDEO_PLAY"
	WS_VIDEO_STOP				= "VIDEO_STOP"
	WS_PHOTO_PLAY				= "PHOTO_PLAY"
	WS_PHOTO_STOP       = "PHOTO_STOP"
	WS_AUDIO_PLAY       = "AUDIO_PLAY"
	WS_AUDIO_STOP       = "AUDIO_STOP"
)

const (
	THEME_DEFAULT				= "DEFAULT"
	THEME_ORANGE        = "ORANGE"
	THEME_TEAL					= "TEAL"
)

const (
	KEY_CURRENT_THEME   = "current-theme"
	KEY_THEME           = "theme"
)

var manager *websocket.Conn

var currentTheme string = THEME_DEFAULT

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

		// TODO: check for error unmarshalling json

		json.Unmarshal(msg, &mc)

		switch mc.Cmd {
		case WS_LOGO:
			pushMap(WS_LOGO, mc.Options)

		case WS_LOGIN:
			log.Println(WS_LOGIN)

		case WS_SCOREBOARD:

			log.Println("fag")
			log.Println(mc)

			if game != nil && game.Active {
				pushMap(WS_SCOREBOARD, mc.Options)
			} else {
				pushMap(WS_SETUP, mc.Options)
			}

		case WS_ADVERTISEMENT:
			pushMap(WS_ADVERTISEMENT, mc.Options)

		case WS_VIDEO_PLAY:

			go videoPlay(mc.Options)

			//pushMap(WS_VIDEO_PLAY, mc.Options)

		case WS_VIDEO_STOP:

			err := videoStop()

			if err != nil {
				log.Println(err)
			}

			//pushMap(WS_VIDEO_STOP, mc.Options)

		case WS_AUDIO_PLAY:
			pushMap(WS_AUDIO_PLAY, mc.Options)

		case WS_AUDIO_STOP:
			pushMap(WS_AUDIO_STOP, mc.Options)

		case WS_PHOTO_PLAY:
			pushMap(WS_PHOTO_PLAY, mc.Options)

		case WS_PHOTO_STOP:
			pushMap(WS_PHOTO_STOP, mc.Options)

		case WS_THEME:

			if mc.Options[KEY_THEME] != "" {

				newTheme := strings.ToUpper(mc.Options[KEY_THEME])

				if newTheme == THEME_ORANGE {
					currentTheme = THEME_ORANGE
				} else if newTheme == THEME_TEAL {
					currentTheme = THEME_TEAL
				} else {
					currentTheme = THEME_DEFAULT
				}
	
				pushMap(WS_THEME, mc.Options)
	
			} else {
				log.Println("mboard-go manager(): theme was empty")
			}

		case WS_THEME_CURRENT:

			m := make(map[string]string)

			m[KEY_CURRENT_THEME] = currentTheme

			pushMap(WS_THEME_CURRENT, m)

		default:
			log.Println("unknown")
		}

	}

} // managerHandler
