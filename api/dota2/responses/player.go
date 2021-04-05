package responses

type Player struct {
	AccountID         int64 `json:"account_id"`
	PlayerSlot        int   `json:"player_slot"`
	HeroID            int   `json:"hero_id"`
	Item0             int   `json:"item_0"`
	Item1             int   `json:"item_1"`
	Item2             int   `json:"item_2"`
	Item3             int   `json:"item_3"`
	Item4             int   `json:"item_4"`
	Item5             int   `json:"item_5"`
	Backpack0         int   `json:"backpack_0"`
	Backpack1         int   `json:"backpack_1"`
	Backpack2         int   `json:"backpack_2"`
	Kills             int   `json:"kills"`
	Deaths            int   `json:"deaths"`
	Assists           int   `json:"assists"`
	LeaverStatus      int   `json:"leaver_status"`
	LastHits          int   `json:"last_hits"`
	Denies            int   `json:"denies"`
	GoldPerMin        int   `json:"gold_per_min"`
	XpPerMin          int   `json:"xp_per_min"`
	Level             int   `json:"level"`
	HeroDamage        int   `json:"hero_damage"`
	TowerDamage       int   `json:"tower_damage"`
	HeroHealing       int   `json:"hero_healing"`
	Gold              int   `json:"gold"`
	GoldSpent         int   `json:"gold_spent"`
	ScaledHeroDamage  int   `json:"scaled_hero_damage"`
	ScaledTowerDamage int   `json:"scaled_tower_damage"`
	ScaledHeroHealing int   `json:"scaled_hero_healing"`
	AbilityUpgrades   []struct {
		Ability int `json:"ability"`
		Time    int `json:"time"`
		Level   int `json:"level"`
	} `json:"ability_upgrades"`
}
