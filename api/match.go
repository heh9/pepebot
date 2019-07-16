package api

import (
	"log"
	"fmt"
	"time"
	"strconv"
	"net/http"
	"github.com/pkg/errors"
	"github.com/bwmarrin/discordgo"
	"github.com/iamalirezaj/go-opendota"
)

var client = opendota.NewClient(http.DefaultClient)

func GetMostHeroDamage(match opendota.Match) (int ,opendota.Hero) {

	var mostHeroDamage = 0
	var mostHeroDamagePlayer = opendota.MatchPlayer{}

	for i, player := range match.Players {
		if i == 0 || player.HeroDamage > mostHeroDamage {
			mostHeroDamage = player.HeroDamage
			mostHeroDamagePlayer = player
		}
	}

	hero, _ := GetHeroById(strconv.Itoa(mostHeroDamagePlayer.HeroID))
	return mostHeroDamage, hero
}

func GetMostKillsHero(match opendota.Match) (int, opendota.Hero) {

	var mostHeroKills = 0
	var mostHeroDamagePlayer = opendota.MatchPlayer{}

	for i, player := range match.Players {
		if i == 0 || player.Kills > mostHeroKills {
			mostHeroKills = player.Kills
			mostHeroDamagePlayer = player
		}
	}

	hero, _ := GetHeroById(strconv.Itoa(mostHeroDamagePlayer.HeroID))
	return mostHeroKills, hero
}

func GetMostHeroHealing(match opendota.Match) (int, opendota.Hero) {

	var mostHeroHealing = 0
	var mostHeroDamagePlayer = opendota.MatchPlayer{}

	for i, player := range match.Players {
		if i == 0 || player.HeroHealing > mostHeroHealing {
			mostHeroHealing = player.HeroHealing
			mostHeroDamagePlayer = player
		}
	}

	hero, _ := GetHeroById(strconv.Itoa(mostHeroDamagePlayer.HeroID))
	return mostHeroHealing, hero
}

func GetMostHeroDenies(match opendota.Match) (int, opendota.Hero) {

	var denies = 0
	var mostHeroDamagePlayer = opendota.MatchPlayer{}

	for i, player := range match.Players {
		if i == 0 || player.Denies > denies {
			denies = player.Denies
			mostHeroDamagePlayer = player
		}
	}

	hero, _ := GetHeroById(strconv.Itoa(mostHeroDamagePlayer.HeroID))
	return denies, hero
}

func GetMostCampsStackedHero(match opendota.Match) (int, opendota.Hero) {

	var camps = 0
	var mostHeroDamagePlayer = opendota.MatchPlayer{}

	for i, player := range match.Players {
		if i == 0 || player.CampsStacked > camps {
			camps = player.CampsStacked
			mostHeroDamagePlayer = player
		}
	}

	hero, _ := GetHeroById(strconv.Itoa(mostHeroDamagePlayer.HeroID))
	return camps, hero
}

func GetMostHeroTowerDamage(match opendota.Match) (int, opendota.Hero) {

	var damage = 0
	var mostHeroDamagePlayer = opendota.MatchPlayer{}

	for i, player := range match.Players {
		if i == 0 || player.TowerDamage > damage {
			damage = player.TowerDamage
			mostHeroDamagePlayer = player
		}
	}

	hero, _ := GetHeroById(strconv.Itoa(mostHeroDamagePlayer.HeroID))
	return damage, hero
}

func GetMatchHistory(
	matchID int64,
	showGameStatus bool,
	won bool,
	showRandomJokeCommentForStatus bool,
	emoji *discordgo.Emoji) (string, error) {

	message := ""

	match, _, err := client.MatchService.Match(matchID)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if match.StartTime == 0 {
		return "", errors.New("Could not find the match")
	}

	unixIntValue, err := strconv.ParseInt(strconv.Itoa(match.Duration), 10, 64)
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
		message += fmt.Sprintf("**Game ended with duration __%s__ as %s with match id** `[%d]`", clock, statusText, matchID)
	} else {
		message += fmt.Sprintf("**Game ended with duration __%s__ and with match id** `[%d]`", clock, matchID)
	}

	message += fmt.Sprintln() + fmt.Sprintln()

	damage, hero := GetMostHeroDamage(match)
	if damage != 0 {
		message += fmt.Sprintf("Most Hero Damage        **(%d)** `%s`", damage, hero.LocalizedName)
		message += fmt.Sprintln()
	}

	kills, hero := GetMostKillsHero(match)
	if kills != 0 {
		message += fmt.Sprintf("Most Kills                          **(%d)** `%s`", kills, hero.LocalizedName)
		message += fmt.Sprintln()
	}

	healing, hero := GetMostHeroHealing(match)
	if healing != 0 {
		message += fmt.Sprintf("Most Hero Healing         **(%d)** `%s`", healing, hero.LocalizedName)
		message += fmt.Sprintln()
	}

	denies, hero := GetMostHeroDenies(match)
	if denies != 0 {
		message += fmt.Sprintf("Most Denies                    **(%d)** `%s`", denies, hero.LocalizedName)
		message += fmt.Sprintln()
	}

	camps, hero := GetMostCampsStackedHero(match)
	if camps != 0 {
		message += fmt.Sprintf("Most Camps Stacked    **(%d)** `%s`", camps, hero.LocalizedName)
		message += fmt.Sprintln()
	}

	towerDamage, hero := GetMostHeroTowerDamage(match)
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
