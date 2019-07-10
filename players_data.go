package main

import (
	"log"
	"strconv"
	"net/http"
	"github.com/iamalirezaj/go-opendota"
)

type Player struct {
	Name              string
	RankTier          int
	RankName          string
	Avatar            string
	MostHeroPlayed    []MostHeroPlayed
}

type MostHeroPlayed struct {
	GamesCount   int
	Hero         opendota.Hero
}

var (
	err      error
	heros    []opendota.Hero
)

func GetHeroById(heroId string) opendota.Hero {
	for _, hero := range heros {
		if id, _ := strconv.Atoi(heroId); hero.ID == id {
			return hero
		}
	}

	return opendota.Hero{}
}

func GetPlayerDatasByMatchID(matchId int) (map[string] []Player, error) {

	var players  = make(map[string] []Player)

	// OpenDota client
	client := opendota.NewClient(http.DefaultClient)

	heros, _, err = client.HeroService.Heroes()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Get Player Data
	match, _, err := client.MatchService.Match(int64(matchId))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, player := range match.Players {

		teamName := "dire"
		if player.IsRadiant {
			teamName = "radiant"
		}

		accountID := int64(player.AccountID)

		player, _, err := client.PlayerService.Player(accountID)
		if err != nil {
			log.Println(err)
			break
		}

		if accountID != 0 {
			playerHeros, _, err :=  client.PlayerService.Heroes(accountID, &opendota.PlayerParam{
				Significant: 1,
			})
			if err != nil {
				log.Println(err)
				break
			}

			var mostHeroPlayed []MostHeroPlayed

			for index, h := range playerHeros {

				if index == 2 { break }

				hero := GetHeroById(h.HeroID)
				mostHeroPlayed = append(mostHeroPlayed, MostHeroPlayed{
					GamesCount:   h.Games,
					Hero:         hero,
				})
			}

			players[teamName] = append(players[teamName], Player{
				Name:             player.Profile.Personaname,
				RankName:         getMedalText(player.RankTier),
				Avatar:           player.Profile.AvatarFull,
				MostHeroPlayed:   mostHeroPlayed,
			})
		}
	}

	return players, nil
}