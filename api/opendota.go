package api

import (
	"github.com/MrJoshLab/pepe.bot/api/dota2/responses"
	"github.com/iamalirezaj/go-opendota"
	"net/http"
	"strconv"
)

func GetMostHeroPlayed(accountID int64) (responses.Hero, error) {
	client := opendota.NewClient(http.DefaultClient)
	playerHeros, _, _ := client.PlayerService.Heroes(accountID, nil)
	heroID, _ := strconv.Atoi(playerHeros[0].HeroID)
	return responses.GetHeroByID(heroID)
}