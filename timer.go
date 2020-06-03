package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"time"
)

type Timer struct {
	ChannelID,
	Seconds,
	Minutes,
	Hours             string

	Value             time.Duration

	MessageReaction   *discordgo.Message
	MessageCreate     *discordgo.MessageCreate

	DoneChannel       chan bool
	ReactedUser       *discordgo.User

	Done              bool

	Client            *discordgo.Session
	TimeRemaining     time.Duration

	StartedTime       time.Time
}

func (timer *Timer) UpdateTimerMessage(reactedUserId string)  {

	clock := time.Since(timer.StartedTime).Round(time.Second)

	tickTime := time.Duration(clock.Seconds()) * time.Second

	remainingDuration := (timer.Value - tickTime).Round(time.Second)

	reactedUser, _ := timer.Client.User(reactedUserId)

	remainingHours := strconv.FormatFloat(remainingDuration.Hours(), 'f', 0, 64)
	remainingMinutes := strconv.FormatFloat(remainingDuration.Minutes(), 'f', 0, 64)
	remainingSeconds := strconv.FormatFloat(remainingDuration.Seconds(), 'f', 0, 64)

	if len(remainingHours) == 1 {
		remainingHours = "0" + remainingHours
	}

	if len(remainingMinutes) == 1 {
		remainingMinutes = "0" + remainingMinutes
	}

	if len(remainingSeconds) == 1 {
		remainingSeconds = "0" + remainingSeconds
	}

	if !timer.Done {

		_, _ = timer.Client.ChannelMessageEdit(timer.ChannelID, timer.MessageReaction.ID,
			fmt.Sprintf(
				"%s <:alarm_clock:603595913058975758> Timer started with duration : `%s:%s:%s` \n" +
					"You can use the 🛑 button to stop the timer or see the remaining time with 👀 button! \n\n" +
					"⏲ `%s:%s:%s` Time Remaining ...",
				timer.MessageCreate.Author.Mention(),
				timer.Hours,
				timer.Minutes,
				timer.Seconds,
				remainingHours,
				remainingMinutes,
				remainingSeconds))

		_ = timer.Client.MessageReactionRemove(timer.ChannelID, timer.MessageReaction.ID, "👀", reactedUser.ID)
	}
}

func (timer *Timer) Stop(user string)  {
	timer.ReactedUser, _ = timer.Client.User(user)
	timer.DoneChannel <- true
}

func (timer *Timer) sendDoneMessage(DoneByUser bool)  {

	_, _ = timer.Client.ChannelMessageEdit(timer.ChannelID, timer.MessageReaction.ID,
		fmt.Sprintf(
			"%s <:alarm_clock:603595913058975758> Timer started with duration : `%s:%s:%s` \n" +
				"💥 Timer Expired! \n\n",
			timer.MessageCreate.Author.Mention(),
			timer.Hours,
			timer.Minutes,
			timer.Seconds))

	_ = timer.Client.MessageReactionsRemoveAll(timer.ChannelID, timer.MessageReaction.ID)

	if DoneByUser {

		_, _ = timer.Client.ChannelMessageSend(timer.ChannelID,
			"@here <:alarm_clock:603595913058975758> Timer stopped by: " + timer.ReactedUser.Mention())

	} else {

		_, _ = timer.Client.ChannelMessageSend(timer.ChannelID,
			fmt.Sprintf("@here <:alarm_clock:603595913058975758> Timer ended with duration: `%s:%s:%s`",
				timer.Hours, timer.Minutes, timer.Seconds))
	}

	timer.Done = true
}

func (timer *Timer) Start()  {

	timer.StartedTime = time.Now()

	if len(timer.Hours) == 1 {
		timer.Hours = "0" + timer.Hours
	}

	if len(timer.Minutes) == 1 {
		timer.Minutes = "0" + timer.Minutes
	}

	if len(timer.Seconds) == 1 {
		timer.Seconds = "0" + timer.Seconds
	}

	msg, _ := timer.Client.ChannelMessageSend(timer.ChannelID,
		fmt.Sprintf(
			"%s <:alarm_clock:603595913058975758> Timer started with duration : `%s:%s:%s` \n" +
				"You can use the 🛑 button to stop the timer or see the remaining time with 👀 button!",
			timer.MessageCreate.Author.Mention(), timer.Hours, timer.Minutes, timer.Seconds))

	timer.MessageReaction = msg

	err := timer.Client.MessageReactionAdd(timer.ChannelID, msg.ID, "🛑")
	if err != nil {
		log.Println(err)
	}

	err = timer.Client.MessageReactionAdd(timer.ChannelID, msg.ID, "👀")
	if err != nil {
		log.Println(err)
	}

	select {
	case <-timer.DoneChannel:
		if !timer.Done {
			timer.sendDoneMessage(true)
		}
		break
	case <-time.After(timer.Value):
		timer.sendDoneMessage(false)
		break
	}
}