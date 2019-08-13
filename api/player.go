package api

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/bwmarrin/discordgo"
)

type SteamPlayer struct {
	Name       string   `json:"name"`
	SteamID    string   `json:"steam_id"`
	AccountID  int64    `json:"account_id"`
	DiscordID  string   `json:"discord_id"`
}

func GetDiscordUserBySteamAccountID(discord *discordgo.Session, accountID int64) (*discordgo.User, error) {

	player , err := GetPlayerByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	user, err := discord.User(player.DiscordID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetPlayerByAccountID(accountID int64) (*SteamPlayer, error)  {

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
		if player.AccountID == accountID {
			return player, nil
		}
	}

	return nil, errors.New("Could not found you in database")
}
