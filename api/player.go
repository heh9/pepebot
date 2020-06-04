package api

import (
	"context"
	"github.com/MrJoshLab/pepe.bot/db"
	"github.com/MrJoshLab/pepe.bot/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
)

type SteamPlayer struct {
	Name       string   `json:"name"`
	SteamID    string   `json:"steam_id"`
	AccountID  int64    `json:"account_id"`
	DiscordID  string   `json:"discord_id"`
}

func GetDiscordUserBySteamAccountID(d *discordgo.Session, g *discordgo.Guild, accountID int64) (*discordgo.User, error) {

	player , err := GetPlayerByAccountID(g, accountID)
	if err != nil {
		return nil, err
	}

	user, err := d.User(player.UserDiscordID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetPlayerByAccountID(guild *discordgo.Guild, accountID int64) (*models.Player, error)  {

	var (
		player     = new(models.Player)
		collection = db.Connection.Collection("players")
		result     = collection.FindOne(context.Background(), bson.M{
			"account_id": strconv.Itoa(int(accountID)),
			"guild_id": guild.ID,
		})
	)

	if err := result.Err(); err != nil {
		return nil, err
	}

	if err := result.Decode(player); err != nil {
		return nil, err
	}

	return player, nil
}
