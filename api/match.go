package api

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mrjosh/pepebot/api/dota2"
	"github.com/mrjosh/pepebot/config"
)

func GetMatchHistory(mid string, sgs bool, w bool, rJ bool, d *discordgo.Session, g *discordgo.Guild) (string, error) {

	message := ""
	client := dota2.NewClient(config.Map.Steam.WebApiToken)

	match, err := client.Match(mid)
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

	if sgs {
		statusText := "loss"
		if w {
			statusText = "win"
		}
		message += fmt.Sprintf("**Game ended with duration __%s__ as %s and match id** `[%s]`", clock, statusText, mid)
	} else {
		message += fmt.Sprintf("**Game ended with duration __%s__ and match id** `[%s]`", clock, mid)
	}

	message += fmt.Sprintln() + fmt.Sprintln()

	if damage, hero, player := match.GetMostHeroDamage(); damage != 0 {

		if u, err := GetDiscordUserBySteamAccountID(d, g, player.AccountID); err == nil {

			message += fmt.Sprintf("Most Hero Damage        **(%d)** `%s` %s", damage, hero.LocalizedName, u.Mention())
		} else {

			message += fmt.Sprintf("Most Hero Damage        **(%d)** `%s`", damage, hero.LocalizedName)
		}

		message += fmt.Sprintln()
	}

	if kills, hero, player := match.GetMostHeroKills(); kills != 0 {

		if u, err := GetDiscordUserBySteamAccountID(d, g, player.AccountID); err == nil {

			message += fmt.Sprintf("Most Kills                          **(%d)** `%s` %s", kills, hero.LocalizedName, u.Mention())
		} else {

			message += fmt.Sprintf("Most Kills                          **(%d)** `%s`", kills, hero.LocalizedName)
		}

		message += fmt.Sprintln()
	}

	if lastHits, hero, player := match.GetMostHeroLastHits(); lastHits != 0 {

		if u, err := GetDiscordUserBySteamAccountID(d, g, player.AccountID); err == nil {

			message += fmt.Sprintf("Most Last Hits                 **(%d)** `%s` %s", lastHits, hero.LocalizedName, u.Mention())
		} else {

			message += fmt.Sprintf("Most Last Hits                 **(%d)** `%s`", lastHits, hero.LocalizedName)
		}

		message += fmt.Sprintln()
	}

	if healing, hero, player := match.GetMostHeroHealing(); healing != 0 {

		if u, err := GetDiscordUserBySteamAccountID(d, g, player.AccountID); err == nil {

			message += fmt.Sprintf("Most Hero Healing         **(%d)** `%s` %s", healing, hero.LocalizedName, u.Mention())
		} else {

			message += fmt.Sprintf("Most Hero Healing         **(%d)** `%s`", healing, hero.LocalizedName)
		}

		message += fmt.Sprintln()
	}

	if denies, hero, player := match.GetMostHeroDenies(); denies != 0 {

		if u, err := GetDiscordUserBySteamAccountID(d, g, player.AccountID); err == nil {

			message += fmt.Sprintf("Most Denies                    **(%d)** `%s` %s", denies, hero.LocalizedName, u.Mention())
		} else {

			message += fmt.Sprintf("Most Denies                    **(%d)** `%s`", denies, hero.LocalizedName)
		}

		message += fmt.Sprintln()
	}

	if towerDamage, hero, player := match.GetMostHeroTowerDamage(); towerDamage != 0 {

		if u, err := GetDiscordUserBySteamAccountID(d, g, player.AccountID); err == nil {

			message += fmt.Sprintf("Most Tower Damage    **(%d)** `%s` %s", towerDamage, hero.LocalizedName, u.Mention())
		} else {

			message += fmt.Sprintf("Most Tower Damage    **(%d)** `%s`", towerDamage, hero.LocalizedName)
		}
	}

	if rJ {

		message += fmt.Sprintln() + fmt.Sprintln()

		comment := "Try a bit harder next time"
		if w {
			comment = "Youuuuuu areeeee the championssssssss my friendsss"
		}

		message += fmt.Sprint(comment)
	}

	return message, nil
}
