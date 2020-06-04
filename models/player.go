package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Player struct {
	ID             *primitive.ObjectID   `bson:"_id, omitempty" json:"id, omitempty"`
	Name           string                `bson:"name, omitempty" json:"name"`
	AccountID      string                `bson:"account_id, omitempty" json:"account_id"`
	UserDiscordID  string                `bson:"user_discord_id, omitempty" json:"user_discord_id"`
	GuildID        string                `bson:"guild_id, omitempty" json:"guild_id"`
	CreatedAt      time.Time             `bson:"created_at, omitempty" json:"created_at,omitempty"`
}