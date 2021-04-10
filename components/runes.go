package components

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/mrjoshlab/pepe.bot/config"
)

const (
	BountyRunesEveryMinutes = 3
)

type Runes struct {
	ClockTime string
	RuneTimes []string
	Sounds    []string
}

func NewRunes() *Runes {
	return &Runes{
		Sounds: config.Map.Sounds.Runes,
	}
}

func (r *Runes) GetRandomVoiceFileName() string {
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)
	return r.Sounds[random.Intn(len(r.Sounds))]
}

func (r *Runes) Up() (bool, string, error) {

	unixIntValue, err := strconv.ParseInt(r.ClockTime, 10, 64)
	if err != nil {
		return false, "", err
	}

	var (
		timeStamp           = time.Unix(unixIntValue, 0).UTC()
		_, minutes, seconds = timeStamp.Clock()
		clock               = fmt.Sprintf("%d:%d", minutes, seconds)
	)

	// If minutes if 0, basiclly the game just started and we dont
	// Want to remind anything about runes yet.
	if minutes == 0 {
		return false, clock, nil
	}

	// Because we want to calculate every 3 minutes but
	// We only want to remind the runes 15 seconds before the runes
	// spawn, we have to add 15 seconds to the current time and if its 60 seconds
	// We add 1 minute to the minutes value and we get the remaiting of the minutes by BountyRunesEveryMinutes
	//
	// If its 0 == Runes are about to swap
	// Else -> Nothing
	if seconds+15 == 60 {
		if (minutes+1)%BountyRunesEveryMinutes == 0 {
			return true, clock, nil
		}
	}

	return false, clock, nil
}
