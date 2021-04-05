package components

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mrjoshlab/pepe.bot/models"
)

const (
	PreGame              = "DOTA_GAMERULES_STATE_PRE_GAME"
	PostGame             = "DOTA_GAMERULES_STATE_POST_GAME"
	StrategyTime         = "DOTA_GAMERULES_STATE_STRATEGY_TIME"
	HeroSelection        = "DOTA_GAMERULES_STATE_HERO_SELECTION"
	InProgress           = "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS"
	WaitForMapToLoad     = "DOTA_GAMERULES_STATE_WAIT_FOR_MAP_TO_LOAD"
	WaitForPlayersToLoad = "DOTA_GAMERULES_STATE_WAIT_FOR_PLAYERS_TO_LOAD"
)

type GSIResponse struct {
	Auth map[string]string `json:"auth"`

	DiscordGuild *discordgo.Guild
	Guild        *models.Guild

	Provider struct {
		Name      string `json:"name"`
		Appid     int    `json:"appid"`
		Version   int    `json:"version"`
		Timestamp int    `json:"timestamp"`
	}

	Map struct {
		Name                 string `json:"name"`
		Matchid              string `json:"matchid"`
		GameTime             int    `json:"game_time"`
		ClockTime            int    `json:"clock_time"`
		DayTime              bool   `json:"daytime"`
		NightStalkerNight    bool   `json:"nightstalker_night"`
		GameState            string `json:"game_state"`
		Paused               bool   `json:"paused"`
		WinTeam              string `json:"win_team"`
		CustomGameName       string `json:"customgamename"`
		WardPurchaseCooldown int    `json:"ward_purchase_cooldown"`
	}

	Player struct {
		SteamId            string                 `json:"steamid"`
		Name               string                 `json:"name"`
		Activity           string                 `json:"activity"`
		Kills              int                    `json:"kills"`
		Deaths             int                    `json:"deaths"`
		Assists            int                    `json:"assists"`
		LastHits           int                    `json:"last_hits"`
		Denies             int                    `json:"denies"`
		KillStreak         int                    `json:"kill_streak"`
		CommandsIssued     int                    `json:"commands_issued"`
		KillList           map[string]interface{} `json:"kill_list"`
		TeamName           string                 `json:"team_name"`
		Gold               int                    `json:"gold"`
		GoldReliable       int                    `json:"gold_reliable"`
		GoldUnreliable     int                    `json:"gold_unreliable"`
		GoldFromHeroKills  int                    `json:"gold_from_hero_kills"`
		GoldFromCreepKills int                    `json:"gold_from_creep_kills"`
		GoldFromIncome     int                    `json:"gold_from_income"`
		GoldFromShared     int                    `json:"gold_from_shared"`
		GPM                int                    `json:"gpm"`
		XPM                int                    `json:"xpm"`
	}
}

func (r *GSIResponse) GetAuthToken() string {
	return r.Auth["token"]
}
