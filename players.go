package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/pkg/errors"
)

type SteamPlayer struct {
	Name       string   `json:"name"`
	SteamID    string   `json:"steam_id"`
	DiscordID  string   `json:"discord_id"`
}

func FindSteamPlayer(discordID string) (*SteamPlayer, error) {

	jsonFile, err := os.Open("players.json")
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var players []*SteamPlayer
	json.Unmarshal(byteValue, &players)

	for _, player := range players {
		if player.DiscordID == discordID {
			return player, nil
		}
	}

	return nil, errors.New("Could not found you in database")
}