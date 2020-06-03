package api

import (
	"fmt"
	"github.com/MrJoshLab/pepe.bot/api/dota2"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
	"time"
)

func GetMatchHistory(
	matchID string,
	showGameStatus bool,
	won bool,
	showRandomJokeCommentForStatus bool,
	emoji *discordgo.Emoji,
	discord *discordgo.Session) (string, error) {

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

	if damage, hero, player := match.GetMostHeroDamage(); damage != 0 {

		if u, err := GetDiscordUserBySteamAccountID(discord, player.AccountID); err == nil {

			message += fmt.Sprintf("Most Hero Damage        **(%d)** `%s` %s", damage, hero.LocalizedName, u.Mention())
		} else {

			message += fmt.Sprintf("Most Hero Damage        **(%d)** `%s`", damage, hero.LocalizedName)
		}

		message += fmt.Sprintln()
	}


	if kills, hero, player := match.GetMostHeroKills(); kills != 0 {

		if u, err := GetDiscordUserBySteamAccountID(discord, player.AccountID); err == nil {

			message += fmt.Sprintf("Most Kills                          **(%d)** `%s` %s", kills, hero.LocalizedName, u.Mention())
		} else {

			message += fmt.Sprintf("Most Kills                          **(%d)** `%s`", kills, hero.LocalizedName)
		}

		message += fmt.Sprintln()
	}

	if lastHits, hero, player := match.GetMostHeroLastHits(); lastHits != 0 {

		if u, err := GetDiscordUserBySteamAccountID(discord, player.AccountID); err == nil {

			message += fmt.Sprintf("Most Last Hits                 **(%d)** `%s` %s", lastHits, hero.LocalizedName, u.Mention())
		} else {

			message += fmt.Sprintf("Most Last Hits                 **(%d)** `%s`", lastHits, hero.LocalizedName)
		}

		message += fmt.Sprintln()
	}


	if healing, hero, player := match.GetMostHeroHealing(); healing != 0 {

		if u, err := GetDiscordUserBySteamAccountID(discord, player.AccountID); err == nil {

			message += fmt.Sprintf("Most Hero Healing         **(%d)** `%s` %s", healing, hero.LocalizedName, u.Mention())
		} else {

			message += fmt.Sprintf("Most Hero Healing         **(%d)** `%s`", healing, hero.LocalizedName)
		}

		message += fmt.Sprintln()
	}

	if denies, hero, player := match.GetMostHeroDenies(); denies != 0 {

		if u, err := GetDiscordUserBySteamAccountID(discord, player.AccountID); err == nil {

			message += fmt.Sprintf("Most Denies                    **(%d)** `%s` %s", denies, hero.LocalizedName, u.Mention())
		} else {

			message += fmt.Sprintf("Most Denies                    **(%d)** `%s`", denies, hero.LocalizedName)
		}

		message += fmt.Sprintln()
	}

	if towerDamage, hero, player := match.GetMostHeroTowerDamage(); towerDamage != 0 {

		if u, err := GetDiscordUserBySteamAccountID(discord, player.AccountID); err == nil {

			message += fmt.Sprintf("Most Tower Damage    **(%d)** `%s` %s", towerDamage, hero.LocalizedName, u.Mention())
		} else {

			message += fmt.Sprintf("Most Tower Damage    **(%d)** `%s`", towerDamage, hero.LocalizedName)
		}
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
