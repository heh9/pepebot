package main

import (
	"fmt"
	"time"
	"strconv"
)

type Runes struct {
	Data   string
	Saied  bool
}

func (r *Runes) AreRunesUp() bool {

	unixIntValue, err := strconv.ParseInt(r.Data, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	timeStamp := time.Unix(unixIntValue, 0).UTC()

	_, mins, secs := timeStamp.Clock()

	seconds := strconv.Itoa(secs)
	minutes := strconv.Itoa(mins)

	fmt.Println("Minutes : " + minutes, " Seconds : " + seconds)

	if minutes == "0" {
		return false
	}

	if len(minutes) == 2 {
		if (minutes[1:2] == "4" && seconds == "45") ||
			(minutes[1:2] == "9" && seconds == "45") {
			return true
		}
	} else {
		if (minutes == "4" && seconds == "45") ||
			(minutes == "9" && seconds == "45") {
			return true
		}
	}

	return false
}