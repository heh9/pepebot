package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Guild struct {
	ID                 uint32    `gorm:"primary_key" json:"id"`
	Name               string    `bson:"name" json:"name"`
	DiscordID          string    `bson:"discord_id" json:"discord_id"`
	UserID             string    `bson:"user_id" json:"user_id"`
	Deleted            bool      `bson:"deleted" json:"deleted"`
	MainTextChannelID  string    `bson:"main_text_channel_id" json:"main_text_channel_id"`
	MainVoiceChannelID string    `bson:"main_voice_channel_id" json:"main_voice_channel_id"`
	Token              string    `bson:"token" json:"token"`
	CreatedAt          time.Time `bson:"created_at" json:"created_at"`
	DeletedAt          time.Time `bson:"deleted_at" json:"deleted_at"`
}

func (g *Guild) BeforeCreate(db *gorm.DB) error {
	g.ID = uuid.New().ID()
	g.CreatedAt = time.Now()
	g.DeletedAt = time.Now()
	return nil
}
