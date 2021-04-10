package api

import (
	"net/http"
	"strconv"

	"github.com/mrjosh/pepebot/api/dota2/responses"
	"github.com/iamalirezaj/go-opendota"
)

func GetMostHeroPlayed(accountID int64) (responses.Hero, error) {
	client := opendota.NewClient(http.DefaultClient)
	playerHeros, _, _ := client.PlayerService.Heroes(accountID, nil)
	heroID, _ := strconv.Atoi(playerHeros[0].HeroID)
	return responses.GetHeroByID(heroID)
}
