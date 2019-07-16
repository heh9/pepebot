package api

import (
	"strconv"
	"net/http"
	"github.com/pkg/errors"
	"github.com/iamalirezaj/go-opendota"
)

var heros, _, _ = client.HeroService.Heroes()

func GetMostHeroPlayed(accountID int64) (opendota.Hero, error) {
	client := opendota.NewClient(http.DefaultClient)
	playerHeros, _, _ := client.PlayerService.Heroes(accountID, nil)
	return GetHeroById(playerHeros[0].HeroID)
}

func GetHeroById(heroID string) (opendota.Hero, error) {

	for _, hero := range heros {
		HeroIdInt, _ := strconv.Atoi(heroID)
		if hero.ID == HeroIdInt {
			return hero, nil
		}
	}

	return opendota.Hero{}, errors.New("Could not found the hero")
}