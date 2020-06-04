package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Guild struct {
	ID         *primitive.ObjectID   `bson:"_id, omitempty" json:"id, omitempty"`

	Name       string                `bson:"name, omitempty" json:"name"`
	DiscordID  string                `bson:"discord_id, omitempty" json:"discord_id"`
	UserID     string                `bson:"user_id, omitempty" json:"user_id"`
	Deleted    bool                  `bson:"deleted, omitempty" json:"deleted"`

	MainTextChannelID   string       `bson:"main_text_channel_id, omitempty" json:"main_text_channel_id"`
	MainVoiceChannelID  string       `bson:"main_voice_channel_id, omitempty" json:"main_voice_channel_id"`

	Token      string                `bson:"token, omitempty" json:"token"`

	CreatedAt  time.Time             `bson:"created_at, omitempty" json:"created_at,omitempty"`
	DeletedAt  time.Time             `bson:"deleted_at, omitempty" json:"deleted_at,omitempty"`
}