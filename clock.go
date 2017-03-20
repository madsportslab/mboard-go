package main

import (
	"encoding/json"
	"log"
	"time"
)

type Clock struct {
	Tenths			int		`json:"tenths"`
	Seconds 		int		`json:"seconds"`
	//Minutes     int		`json:"minutes"`
}

type GameClocks struct {
	Ticker            *time.Ticker
	ShotViolationChan chan bool
	FinalChan         chan bool
	OutChan           chan []byte
	PlayClock         *Clock
	ShotClock         *Clock
}

type ReadableClock struct {
	GameClock		*Clock		`json:"game"`
	ShotClock   *Clock		`json:"shot"`
}

func (gc *GameClocks) Run(settings *Config) {

	for _ = range gc.Ticker.C {

		//log.Println(t)

		gc.PlayClock.Seconds++
		gc.ShotClock.Seconds++
	
		if gc.PlayClock.Tenths == 9 {
			gc.PlayClock.Tenths = 0
		} else {
			gc.PlayClock.Tenths++
		}

		if gc.ShotClock.Tenths == 9 {
			gc.ShotClock.Tenths = 0
		} else {
			gc.ShotClock.Tenths++
		}

		if gc.PlayClock.Seconds == settings.Minutes * 60 {
			gc.Ticker.Stop()
			gc.FinalChan <- true
		}

		if gc.ShotClock.Seconds == settings.Shot {
			gc.ShotViolationChan <- true
		}

		rc := ReadableClock{
			ShotClock: gc.ShotClock,
			GameClock: gc.PlayClock,
		}

		j, jsonErr := json.Marshal(rc)

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

	gc.PlayClock.Seconds 	= 0
	gc.PlayClock.Tenths 	= 0

} // GameClockReset

func (gc *GameClocks) stepGameClock(settings *Config, ticks int) {

	if gc.Ticker != nil {
		gc.Ticker.Stop()
	}

	total := gc.PlayClock.Seconds + ticks

	if total >= 0 && total < settings.Minutes * 60 {
		gc.PlayClock.Seconds = total
	}

	if total == settings.Minutes * 60 {
		gc.FinalChan <- true
	}

} // stepGameClock

func (gc *GameClocks) stepShotClock(settings *Config, ticks int) {

	if gc.Ticker != nil {
		gc.Ticker.Stop()
	}

  total := gc.ShotClock.Seconds + ticks

	if total >= 0 && total < settings.Shot {
		gc.ShotClock.Seconds = total
	}

	if gc.ShotClock.Seconds == settings.Shot {
		gc.ShotViolationChan <- true
	}

} // stepShotClock
