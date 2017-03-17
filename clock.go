package main

import (
	"encoding/json"
	"log"
	"time"
)

type Clock struct {
	Tenths			int		`json:"tenths"`
	Seconds 		int		`json:"seconds"`
	Minutes     int		`json:"minutes"`
}

type GameClocks struct {
	Ticker            *time.Ticker
	ShotViolationChan chan bool
	FinalChan         chan bool
	OutChan           chan []byte
	PlayClock         *Clock
	ShotClock         *Clock
}

func (gc *GameClocks) Run(settings *Config) {

	for t := range gc.Ticker.C {

		if gc.ShotClock.Tenths == 9 {
		
			gc.ShotClock.Tenths = 0
			gc.ShotClock.Seconds++
		
		} else {
			gc.ShotClock.Tenths++
		}

		if gc.PlayClock.Tenths == 9 {

			if gc.PlayClock.Seconds == 59 {
			
				gc.PlayClock.Minutes++
				gc.PlayClock.Seconds = 0
			
			} else {
				gc.PlayClock.Seconds++
			}

			gc.PlayClock.Tenths = 0

		} else {
			gc.PlayClock.Tenths++
		}

		log.Println(gc.PlayClock.Minutes, gc.PlayClock.Seconds,
			gc.PlayClock.Tenths, t)

		log.Println(gc.ShotClock.Seconds,
			gc.ShotClock.Tenths, t)

		if gc.PlayClock.Minutes == settings.Minutes {
			gc.Ticker.Stop()
			gc.FinalChan <- true
		}

		if gc.ShotClock.Seconds == settings.Shot {
			gc.ShotViolationChan <- true
		}

		j, jsonErr := json.Marshal(gc.PlayClock)

		if jsonErr != nil {
			log.Println("[Error]", jsonErr)
		}

		gc.OutChan <- j

	}

} // Run

func (gc *GameClocks) Start(settings *Config) {

	//TODO: prevent multiple starts
	if gc.Ticker != nil {
		gc.Ticker.Stop()
	}

	gc.Ticker = time.NewTicker(time.Millisecond * 100)

	go gc.Run(settings)

	for {
		select {
		case <-gc.ShotViolationChan:
			gc.ShotClockReset()
		case <-gc.FinalChan:
			return
		}
	}

} // Start

func (gc *GameClocks) Stop() {

	if gc.Ticker != nil {
		gc.Ticker.Stop()
	}

} // Stop

func (gc *GameClocks) ShotClockReset() {

	if gc.Ticker != nil {

		gc.ShotClock.Seconds 	= 0
		gc.ShotClock.Tenths 	= 0
	
	}

} // ShotClockReset

func (gc *GameClocks) GameClockReset() {

	if gc.Ticker != nil {
		gc.Ticker.Stop()
	}

	gc.PlayClock.Minutes 	= 0
	gc.PlayClock.Seconds 	= 0
	gc.PlayClock.Tenths 	= 0

} // GameClockReset

func (gc *GameClocks) Rew(ticks int) {

	if ticks < 0 {
		return
	}

	if gc.Ticker != nil {
		gc.Ticker.Stop()
	}

	if gc.PlayClock.Seconds == 0 && gc.PlayClock.Minutes == 0 {
		return
	}

	var delta = gc.PlayClock.Seconds - ticks

	if delta >= 0 {
		gc.PlayClock.Seconds = delta
	} else {
		
		gc.PlayClock.Seconds = 60 + delta
		
		if gc.PlayClock.Minutes > 0 {
			gc.PlayClock.Minutes--
		}

	}

	var sdelta = gc.ShotClock.Seconds - ticks

	if sdelta >= 0 && sdelta < 24 {
		gc.ShotClock.Seconds = sdelta
	} else {
		gc.ShotViolationChan <- true
	}

} // Rew

func (gc *GameClocks) Fwd(ticks int) {

	if ticks < 0 {
		return
	}

	if gc.Ticker != nil {
		gc.Ticker.Stop()
	}

	var sum = gc.PlayClock.Seconds + ticks

	if sum > 59 {

		gc.PlayClock.Minutes++
		gc.PlayClock.Seconds = sum - 60

	} else {
		gc.PlayClock.Seconds = gc.PlayClock.Seconds + ticks
	}

	var ssum = gc.ShotClock.Seconds + ticks

	if ssum >= 0 && ssum < 24 {
		gc.ShotClock.Seconds = ssum
	} else {
		gc.ShotViolationChan <- true
	}

} // Fwd
