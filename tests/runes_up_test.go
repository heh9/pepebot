package tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/mrjosh/pepebot/components"
	"github.com/stretchr/testify/assert"
)

type TestTime struct {
	Value         int
	ExpectedClock string
	ExpectedUp    bool
	NoError       bool
}

func TestRunesAreUpEvery3Minutes(t *testing.T) {

	testTimes := []TestTime{
		{60, "1:0", false, false},
		{165, "2:45", true, false},
		{345, "5:45", true, false},
		{525, "8:45", true, false},
		{524, "8:44", false, false},
		{705, "11:45", true, false},
		{885, "14:45", true, false},
		{1065, "17:45", true, false},
		{1245, "20:45", true, false},
	}

	for _, testTime := range testTimes {
		t.Run(fmt.Sprintf("TestTime[%s]", testTime.ExpectedClock), func(t *testing.T) {
			gm := &components.GuildMatch{
				Runes: &components.Runes{
					ClockTime: strconv.Itoa(testTime.Value),
					RuneTimes: make([]string, 0),
				},
			}
			ok, clock, err := gm.Runes.Up()
			if testTime.NoError {
				assert.NoError(t, err)
			}
			assert.Equal(t, clock, testTime.ExpectedClock)
			if testTime.ExpectedUp {
				assert.True(t, ok)
			} else {
				assert.False(t, ok)
			}
		})
	}

}
