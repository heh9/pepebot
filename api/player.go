package api

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/mrjoshlab/pepe.bot/db"
	"github.com/mrjoshlab/pepe.bot/models"
)

type SteamPlayer struct {
	Name      string `json:"name"`
	SteamID   string `json:"steam_id"`
	AccountID int64  `json:"account_id"`
	DiscordID string `json:"discord_id"`
}

func GetDiscordUserBySteamAccountID(d *discordgo.Session, g *discordgo.Guild, accountID int64) (*discordgo.User, error) {

	player, err := GetPlayerByAccountID(g, accountID)
	if err != nil {
		return nil, err
	}

	user, err := d.User(player.UserDiscordID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetPlayerByAccountID(guild *discordgo.Guild, accountID int64) (*models.Player, error) {
	var (
		player = new(models.Player)
		result = db.Connection.Where("account_id =?", strconv.Itoa(int(accountID))).
			Where("guild_id =?", guild.ID).
			First(&player)
	)
	if err := result.Error; err != nil {
		return nil, err
	}
	return player, nil
}
