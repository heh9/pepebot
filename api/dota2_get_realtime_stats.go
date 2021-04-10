package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iamalirezaj/go-opendota"
	"github.com/mrjosh/pepebot/config"
)

type Response struct {
	Match struct {
		ServerSteamId int64 `json:"server_steam_id"`
		MatchID       int64 `json:"matchid"`
	} `json:"match"`

	Teams []Team `json:"teams"`
}

type Team struct {
	Players []Player `json:"players"`
}

type Player struct {
	AccountID int64  `json:"accountid"`
	Name      string `json:"name"`
}

type OpenDotaPlayer struct {
	AccountID int64
	Name      string
	Rank      string
}

func GetPlayerOpenDotaProfile(accountID int64) (*OpenDotaPlayer, error) {

	client := opendota.NewClient(http.DefaultClient)

	player, _, err := client.PlayerService.Player(accountID)
	if err != nil {
		return nil, err
	}

	return &OpenDotaPlayer{
		AccountID: int64(player.Profile.AccountID),
		Name:      player.Profile.Personaname,
		Rank:      GetPlayerMedalString(player.RankTier),
	}, nil
}

func GetRealTimeStats(serverSteamID string) (*Response, error) {

	emptyResponse := new(Response)

	client := http.Client{}
	uri := "https://api.steampowered.com/IDOTA2MatchStats_570/GetRealtimeStats/v001/?server_steam_id=" +
		serverSteamID + "&key=" + config.Map.Steam.WebApiToken

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return emptyResponse, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return emptyResponse, err
	}

	response := &Response{}

	jsonErr := json.NewDecoder(resp.Body).Decode(response)
	if jsonErr != nil {
		return emptyResponse, jsonErr
	}

	if resp.StatusCode != http.StatusOK {
		return response, errors.New(resp.Status)
	}

	return response, nil
}
