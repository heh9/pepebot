package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Player struct {
	ID            uint32    `gorm:"primary_key" json:"id"`
	Name          string    `bson:"name" json:"name"`
	AccountID     string    `bson:"account_id" json:"account_id"`
	UserDiscordID string    `bson:"user_discord_id" json:"user_discord_id"`
	GuildID       string    `bson:"guild_id" json:"guild_id"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
}

func (p *Player) BeforeCreate(db *gorm.DB) error {
	p.ID = uuid.New().ID()
	p.CreatedAt = time.Now()
	return nil
}
