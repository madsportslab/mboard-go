package main

import (
	"encoding/json"
	"fmt"
	"log"
  "net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	WS_CLOCK_START			= "CLOCK_START"
	WS_CLOCK_STOP				= "CLOCK_STOP"
	WS_CLOCK_RESET      = "CLOCK_RESET"
	WS_CLOCK_STEP       = "CLOCK_STEP"
	WS_SHOT_RESET       = "SHOT_RESET"
	WS_SHOT_STEP        = "SHOT_STEP"
	WS_PERIOD_UP        = "PERIOD_UP"
	WS_PERIOD_DOWN      = "PERIOD_DOWN"
	WS_POSSESSION_HOME  = "POSSESSION_HOME"
	WS_POSSESSION_AWAY  = "POSSESSION_AWAY"
	WS_FINAL            = "FINAL"
	WS_ABORT      			= "ABORT"
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

const (
	WS_RET_POSSESSION_HOME    = "POSSESSION_HOME"
	WS_RET_POSSESSION_AWAY    = "POSSESSION_AWAY"
	WS_RET_HOME_SCORE    			= "HOME_SCORE"
	WS_RET_AWAY_SCORE    			= "AWAY_SCORE"
	WS_RET_CLOCK              = "CLOCK"
	WS_RET_PERIOD             = "PERIOD"
	WS_RET_HOME_FOUL          = "HOME_FOUL"
	WS_RET_AWAY_FOUL          = "AWAY_FOUL"
	WS_RET_HOME_TIMEOUT       = "HOME_TIMEOUT"
	WS_RET_AWAY_TIMEOUT       = "AWAY_TIMEOUT"
	WS_RET_SHOT_VIOLATION     = "SHOT_VIOLATION"
	WS_RET_END_PERIOD         = "END_PERIOD"
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

type Notification struct {
  Key 			string				`json:"key"`
	Val				string				`json:"val"`
}

type Req struct {
	Cmd					string 			`json:"cmd"`
	Step				int					`json:"step"`        
}

var periodNames = []string{"1st", "2nd", "3rd", "4th"}

//Games[id], has sockets, each socket has a mutex, and record

//var connections = make(map[string]map[*websocket.Conn]*sync.Mutex)

func calcTotalScore(home bool) int {

  total := 0

	if home {

    for _, v := range game.GameData.Home.Points {
			total = total + v
		}

	} else {

    for _, v := range game.GameData.Away.Points {
			total = total + v
		}

	}

  return total

} // calcTotalScore

func notify(key string, val string) {

	log.Println(key, val)

	n := Notification{
		Key: key,
		Val: val,
	}

	j, jsonErr := json.Marshal(n)

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	for c, mu := range game.Conns {

		mu.Lock()
		c.WriteMessage(websocket.TextMessage, j)
		mu.Unlock()

	}


} // notify

func incrementPoints(name string, val int) {

  if game == nil {
		return
	}

  if name == HOME {

		total := game.GameData.Home.Points[game.GameData.Period]
		
		if (total + val) < 0 {
			return
		}

		game.GameData.Home.Points[game.GameData.Period] = total +
			val

		notify(WS_RET_HOME_SCORE, fmt.Sprintf("%d", calcTotalScore(true)))
		
	} else if name == AWAY {

		total := game.GameData.Away.Points[game.GameData.Period]
		
		if (total + val) < 0 {
			return
		}

		game.GameData.Away.Points[game.GameData.Period] = total +
			val

		notify(WS_RET_AWAY_SCORE, fmt.Sprintf("%d", calcTotalScore(false)))
		
	}

} // incrementPoints

func incrementFoul(name string, val int) {

  if game == nil {
		return
	}

  if name == HOME {

		if game.GameData.Home.Fouls + val < 0 {
			return
		}

		game.GameData.Home.Fouls = game.GameData.Home.Fouls + val
		
		notify(WS_RET_HOME_FOUL, fmt.Sprintf("%d", game.GameData.Home.Fouls))

	} else if name == AWAY {

		if game.GameData.Away.Fouls + val < 0  {
			return
		}

		game.GameData.Away.Fouls = game.GameData.Away.Fouls + val

		notify(WS_RET_AWAY_FOUL, fmt.Sprintf("%d", game.GameData.Away.Fouls))

	}

} // incrementFoul

func incrementTimeout(name string, val int) {

  if game == nil {
		return
	}

  if name == HOME {

		if (game.Settings.Timeouts - (game.GameData.Home.Timeouts + val) < 0) ||
		  (game.GameData.Home.Timeouts + val < 0) {
			return
		}

		game.GameData.Home.Timeouts = game.GameData.Home.Timeouts + val

		notify(WS_RET_HOME_TIMEOUT, fmt.Sprintf("%d", game.GameData.Home.Timeouts))
		
	} else if name == AWAY {

		if (game.Settings.Timeouts - (game.GameData.Away.Timeouts + val) < 0) ||
		  (game.GameData.Away.Timeouts + val < 0) {
			return
		}

		game.GameData.Away.Timeouts = game.GameData.Away.Timeouts + val

		notify(WS_RET_AWAY_TIMEOUT, fmt.Sprintf("%d", game.GameData.Away.Timeouts))

	}

} // incrementTimeout

func incrementPeriod(val int) {

	if (game.GameData.Period + val) < 0 {
		return
	}

  game.GameData.Period = game.GameData.Period + val

	game.GameData.Clk.GameClockReset()

	if game.GameData.Period < 5 {
		notify(WS_RET_PERIOD, periodNames[game.GameData.Period])
	} else {
		notify(WS_RET_PERIOD, fmt.Sprintf("OT%d",  game.GameData.Period - 3))
	}



} // incrementPeriod

func setPossession(name string) {

  if name == HOME {
 	 	
		game.GameData.Possession = true

		notify(WS_RET_POSSESSION_HOME, "")

	} else {
		
		game.GameData.Possession = false

		notify(WS_RET_POSSESSION_AWAY, "")

	}

} // setPossession

func firehose(game *GameInfo) {

  for {

		select {
		case <-game.GameData.Clk.ShotViolationChan:
		
		  game.GameData.Clk.Ticker.Stop()
			notify(WS_RET_SHOT_VIOLATION, "1")
		
		case <-game.GameData.Clk.FinalChan:

		  game.GameData.Clk.Ticker.Stop()
			notify(WS_RET_END_PERIOD, "1")
		
		case s := <-game.GameData.Clk.OutChan:
		  notify(WS_RET_CLOCK, string(s))
		}

	}

} // firehose

func controlHandler(w http.ResponseWriter, r *http.Request) {

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

	game.Conns[c] = &sync.Mutex{}

	go firehose(game)

	log.Println("connection #:", len(game.Conns))

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
			go game.GameData.Clk.Start()

		case WS_CLOCK_STOP:
		  game.GameData.Clk.Stop()

		case WS_CLOCK_RESET:
		  game.GameData.Clk.GameClockReset()

		case WS_SHOT_RESET:
		  game.GameData.Clk.ShotClockReset()

    case WS_SHOT_STEP:
		  game.GameData.Clk.StepShotClock(req.Step)

		case WS_CLOCK_STEP:
		  game.GameData.Clk.StepGameClock(req.Step)

		case WS_PERIOD_UP:
			incrementPeriod(1)
		
		case WS_PERIOD_DOWN:
		  incrementPeriod(-1)

		case WS_POSSESSION_HOME:
			setPossession(HOME)

		case WS_POSSESSION_AWAY:
		  setPossession(AWAY)

		case WS_FINAL:
			game.Final = true

		case WS_ABORT:

			game.GameData.Clk.Stop()

			// close connections and channels
		
		case WS_SCORE_HOME:

		  log.Println(req.Step)

      incrementPoints(HOME, req.Step)

    case WS_SCORE_AWAY:

		  log.Println(req.Step)

      incrementPoints(AWAY, req.Step)

		case WS_FOUL_HOME_UP:
		  incrementFoul(HOME, 1)

		case WS_FOUL_HOME_DOWN:
			incrementFoul(HOME, -1)
		
		case WS_FOUL_AWAY_UP:
		  incrementFoul(AWAY, 1)

		case WS_FOUL_AWAY_DOWN:
			incrementFoul(AWAY, -1)

		case WS_TIMEOUT_HOME_UP:
			incrementTimeout(HOME, -1)

		case WS_TIMEOUT_HOME_DOWN:
			incrementTimeout(HOME, 1)

		case WS_TIMEOUT_AWAY_UP:
			incrementTimeout(AWAY, -1)

		case WS_TIMEOUT_AWAY_DOWN:
			incrementTimeout(AWAY, 1)
		
		default:
		  log.Printf("[%s][Error] unsupported command: %s", version(), string(msg))
		}

	}

} // controlHandler
