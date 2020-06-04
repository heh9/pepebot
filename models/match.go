package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Match struct {
	ID               *primitive.ObjectID   `bson:"_id, omitempty" json:"id, omitempty"`

	Duration         int                   `bson:"duration, omitempty" json:"duration"`
	Win              bool                  `bson:"win, omitempty" json:"win"`
	MatchID          int                   `bson:"match_id, omitempty" json:"match_id"`

	MostHeroDamage   MatchStat             `bson:"most_hero_damage, omitempty" json:"most_hero_damage"`
	MostKills        MatchStat             `bson:"most_kills, omitempty" json:"most_kills"`
	MostLastHit      MatchStat             `bson:"most_last_hit, omitempty" json:"most_last_hit"`
	MostHeroHealing  MatchStat             `bson:"most_hero_healing, omitempty" json:"most_hero_healing"`
	MostDenies       MatchStat             `bson:"most_denies, omitempty" json:"most_denies"`
	MostTowerDamage  MatchStat             `bson:"most_tower_damage, omitempty" json:"most_tower_damage"`

	CreatedAt        time.Time             `bson:"created_at, omitempty" json:"created_at"`
}

type MatchStat struct {
	Value     int                   `bson:"value, omitempty" json:"value"`
	Hero      int                   `bson:"hero, omitempty" json:"hero"`
	PlayerID  *primitive.ObjectID   `bson:"player_id, omitempty" json:"player_id, omitempty"`
}