package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Match struct {
	ID              uint32    `gorm:"primary_key" json:"id"`
	Duration        int       `bson:"duration" json:"duration"`
	Win             bool      `bson:"win" json:"win"`
	MatchID         int       `bson:"match_id" json:"match_id"`
	MostHeroDamage  MatchStat `bson:"most_hero_damage" json:"most_hero_damage"`
	MostKills       MatchStat `bson:"most_kills" json:"most_kills"`
	MostLastHit     MatchStat `bson:"most_last_hit" json:"most_last_hit"`
	MostHeroHealing MatchStat `bson:"most_hero_healing" json:"most_hero_healing"`
	MostDenies      MatchStat `bson:"most_denies" json:"most_denies"`
	MostTowerDamage MatchStat `bson:"most_tower_damage" json:"most_tower_damage"`
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
}

type MatchStat struct {
	Value    int    `bson:"value" json:"value"`
	Hero     int    `bson:"hero" json:"hero"`
	PlayerID uint32 `bson:"player_id" json:"player_id"`
}

func (m *Match) BeforeCreate(db *gorm.DB) error {
	m.ID = uuid.New().ID()
	m.CreatedAt = time.Now()
	return nil
}
