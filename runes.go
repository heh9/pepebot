package main

import (
	"log"
	"time"
	"strconv"
	"math/rand"
)

type Runes struct {
	ClockTime   string
	RuneTimes   []string
	Sounds      []string
}

func NewRunesType() *Runes {
	return &Runes{
		Sounds: []string {
			"runes",
			"runes_hello",
			"runes_milad",
			"runes_mohammad",
		},
	}
}

func (r *Runes) RNG(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func (r *Runes) GetRandomVoiceFileName() string {
	return r.Sounds[r.RNG(0, len(r.Sounds))]
}

func (r *Runes) Up() (bool, string) {

	unixIntValue, err := strconv.ParseInt(r.ClockTime, 10, 64)
	if err != nil {
		log.Println(err)
	}

	timeStamp := time.Unix(unixIntValue, 0).UTC()

	_, mins, secs := timeStamp.Clock()

	seconds := strconv.Itoa(secs)
	minutes := strconv.Itoa(mins)

	clock := minutes + ":" + seconds

	if minutes == "0" {
		return false, clock
	}

	if len(minutes) == 2 {
		if (minutes[1:2] == "4" && seconds == "45") ||
			(minutes[1:2] == "9" && seconds == "45") {
			return true, clock
		}
	} else {
		if (minutes == "4" && seconds == "45") ||
			(minutes == "9" && seconds == "45") {
			return true, clock
		}
	}

	return false, clock
}