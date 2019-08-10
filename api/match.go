package api

import (
	"os"
	"log"
	"fmt"
	"time"
	"strconv"
	"pepe.bot/api/dota2"
	"github.com/pkg/errors"
	"github.com/bwmarrin/discordgo"
)

func GetMatchHistory(
	matchID string,
	showGameStatus bool,
	won bool,
	showRandomJokeCommentForStatus bool,
	emoji *discordgo.Emoji) (string, error) {

	message := ""

	client := dota2.NewClient(os.Getenv("STEAM_WEBAPI_API_KEY"))

	match, err := client.Match(matchID)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if match.Result.StartTime == 0 {
		return "", errors.New("Could not find the match")
	}

	unixIntValue, err := strconv.ParseInt(strconv.Itoa(match.Result.Duration), 10, 64)
	if err != nil {
		log.Println(err)
		return "", err
	}

	timeStamp := time.Unix(unixIntValue, 0).UTC()
	hours, mins, secs := timeStamp.Clock()

	clock := fmt.Sprintf("%d:%d:%d", hours, mins, secs)
	if hours == 0 {
		clock = fmt.Sprintf("%d:%d", mins, secs)
	}

	if showGameStatus {
		statusText := "loss"
		if won {
			statusText = "win"
		}
		message += fmt.Sprintf("**Game ended with duration __%s__ as %s and match id** `[%s]`", clock, statusText, matchID)
	} else {
		message += fmt.Sprintf("**Game ended with duration __%s__ and match id** `[%s]`", clock, matchID)
	}

	message += fmt.Sprintln() + fmt.Sprintln()

	damage, hero := match.GetMostHeroDamage()
	if damage != 0 {
		message += fmt.Sprintf("Most Hero Damage        **(%d)** `%s`", damage, hero.LocalizedName)
		message += fmt.Sprintln()
	}

	kills, hero := match.GetMostHeroKills()
	if kills != 0 {
		message += fmt.Sprintf("Most Kills                          **(%d)** `%s`", kills, hero.LocalizedName)
		message += fmt.Sprintln()
	}

	healing, hero := match.GetMostHeroHealing()
	if healing != 0 {
		message += fmt.Sprintf("Most Hero Healing         **(%d)** `%s`", healing, hero.LocalizedName)
		message += fmt.Sprintln()
	}

	denies, hero := match.GetMostHeroDenies()
	if denies != 0 {
		message += fmt.Sprintf("Most Denies                    **(%d)** `%s`", denies, hero.LocalizedName)
		message += fmt.Sprintln()
	}

	towerDamage, hero := match.GetMostHeroTowerDamage()
	if towerDamage != 0 {
		message += fmt.Sprintf("Most Tower Damage    **(%d)** `%s`", towerDamage, hero.LocalizedName)
	}

	if showRandomJokeCommentForStatus {

		message += fmt.Sprintln() + fmt.Sprintln()

		comment := "Try a bit harder next time"
		if won {
			comment = "Youuuuuu areeeee the championssssssss my friendsss"
		}

		if emoji != nil {
			comment += " " + emoji.MessageFormat()
		}

		message += fmt.Sprint(comment)
	}

	return message, nil
}
