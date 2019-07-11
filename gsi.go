package main

const (
	PreGame               = "DOTA_GAMERULES_STATE_PRE_GAME"
	PostGame              = "DOTA_GAMERULES_STATE_POST_GAME"
	StrategyTime          = "DOTA_GAMERULES_STATE_STRATEGY_TIME"
	HeroSelection         = "DOTA_GAMERULES_STATE_HERO_SELECTION"
	InProgress            = "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS"
	WaitForMapToLoad      = "DOTA_GAMERULES_STATE_WAIT_FOR_MAP_TO_LOAD"
	WaitForPlayersToLoad  = "DOTA_GAMERULES_STATE_WAIT_FOR_PLAYERS_TO_LOAD"
)

type GSIResponse struct {

	Auth                           map[string] string      `json:"auth"`

	Buildings struct {
		Radiant                    map[string] interface{} `json:"radiant"`
		Dire                       map[string] interface{} `json:"dire"`
	}

	Provider struct {
		Name                       string `json:"name"`
		Appid                      int    `json:"appid"`
		Version                    int    `json:"version"`
		Timestamp                  int    `json:"timestamp"`
	}

	Map struct {
		Name                       string `json:"name"`
		Matchid                    string `json:"matchid"`
		GameTime                   int    `json:"game_time"`
		ClockTime                  int    `json:"clock_time"`
		DayTime                    bool   `json:"daytime"`
		NightStalkerNight          bool   `json:"nightstalker_night"`
		GameState                  string `json:"game_state"`
		Paused                     bool   `json:"paused"`
		WinTeam                    string `json:"win_team"`
		CustomGameName             string `json:"customgamename"`
		WardPurchaseCooldown       int    `json:"ward_purchase_cooldown"`
	}

	Player struct {
		SteamId                    string `json:"steamid"`
		Name                       string `json:"name"`
		Activity                   string `json:"activity"`
		Kills                      int `json:"kills"`
		Deaths                     int `json:"deaths"`
		Assists                    int `json:"assists"`
		LastHits                   int `json:"last_hits"`
		Denies                     int `json:"denies"`
		KillStreak                 int `json:"kill_streak"`
		CommandsIssued             int `json:"commands_issued"`
		KillList                   map[string] interface{} `json:"kill_list"`
		TeamName                   string `json:"team_name"`
		Gold                       int `json:"gold"`
		GoldReliable               int `json:"gold_reliable"`
		GoldUnreliable             int `json:"gold_unreliable"`
		GoldFromHeroKills          int `json:"gold_from_hero_kills"`
		GoldFromCreepKills         int `json:"gold_from_creep_kills"`
		GoldFromIncome             int `json:"gold_from_income"`
		GoldFromShared             int `json:"gold_from_shared"`
		GPM                        int `json:"gpm"`
		XPM                        int `json:"xpm"`
	}

	Hero struct {
		Xpos                       interface{} `json:"xpos"`
		Ypos                       interface{} `json:"ypos"`
		Id                         int `json:"id"`
		Name                       string `json:"name"`
		Level                      int `json:"level"`
		Alive                      bool `json:"alive"`
		RespawnSeconds             int `json:"respawn_seconds"`
		BuybackCost                int `json:"buyback_cost"`
		BuybackCooldown            int `json:"buyback_cooldown"`
		Health                     int `json:"health"`
		MaxHealth                  int `json:"max_health"`
		HealthPercent              int `json:"health_percent"`
		Mana                       int `json:"mana"`
		MaxMana                    int `json:"max_mana"`
		ManaPercent                int `json:"mana_percent"`
		Silenced                   bool `json:"silenced"`
		Stunned                    bool `json:"stunned"`
		Disarmed                   bool `json:"disarmed"`
		MagicImmune                bool `json:"magicimmune"`
		Hexed                      bool `json:"hexed"`
		Muted                      bool `json:"muted"`
		Break                      bool `json:"break"`
		Smoked                     bool `json:"smoked"`
		HasDebuff                  bool `json:"has_debuff"`
		Talent1                    bool `json:"talent_1"`
		Talent2                    bool `json:"talent_2"`
		Talent3                    bool `json:"talent_3"`
		Talent4                    bool `json:"talent_4"`
		Talent5                    bool `json:"talent_5"`
		Talent6                    bool `json:"talent_6"`
		Talent7                    bool `json:"talent_7"`
		Talent8                    bool `json:"talent_8"`
	}
}

func (r *GSIResponse) getAuthToken() string {
	return r.Auth["token"]
}

func (r *GSIResponse) CheckAuthToken(token string) bool {
	return r.getAuthToken() == token
}
