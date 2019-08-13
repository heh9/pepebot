package responses

type Match struct {
	Result struct {
		Players               []Player    `json:"players"`
		RadiantWin            bool        `json:"radiant_win"`
		Duration              int         `json:"duration"`
		PreGameDuration       int         `json:"pre_game_duration"`
		StartTime             int         `json:"start_time"`
		MatchID               int64       `json:"match_id"`
		MatchSeqNum           int64       `json:"match_seq_num"`
		TowerStatusRadiant    int         `json:"tower_status_radiant"`
		TowerStatusDire       int         `json:"tower_status_dire"`
		BarracksStatusRadiant int         `json:"barracks_status_radiant"`
		BarracksStatusDire    int         `json:"barracks_status_dire"`
		Cluster               int         `json:"cluster"`
		FirstBloodTime        int         `json:"first_blood_time"`
		LobbyType             int         `json:"lobby_type"`
		HumanPlayers          int         `json:"human_players"`
		Leagueid              int         `json:"leagueid"`
		PositiveVotes         int         `json:"positive_votes"`
		NegativeVotes         int         `json:"negative_votes"`
		GameMode              int         `json:"game_mode"`
		Flags                 int         `json:"flags"`
		Engine                int         `json:"engine"`
		RadiantScore          int         `json:"radiant_score"`
		DireScore             int         `json:"dire_score"`
		PicksBans             []struct {
			IsPick bool `json:"is_pick"`
			HeroID int  `json:"hero_id"`
			Team   int  `json:"team"`
			Order  int  `json:"order"`
		} `json:"picks_bans"`
	} `json:"result"`
}

func (m *Match) GetMostHeroDamage() (int, Hero, Player) {

	var mostHeroDamage = 0
	var mostHeroDamagePlayer = Player{}

	for i, player := range m.Result.Players {
		if i == 0 || player.HeroDamage > mostHeroDamage {
			mostHeroDamage = player.HeroDamage
			mostHeroDamagePlayer = player
		}
	}

	hero, _ := GetHeroByID(mostHeroDamagePlayer.HeroID)

	return mostHeroDamage, hero, mostHeroDamagePlayer
}

func (m *Match) GetMostHeroKills() (int, Hero, Player) {

	var mostHeroKills = 0
	var mostHeroDamagePlayer = Player{}

	for i, player := range m.Result.Players {
		if i == 0 || player.Kills > mostHeroKills {
			mostHeroKills = player.Kills
			mostHeroDamagePlayer = player
		}
	}

	hero, _ := GetHeroByID(mostHeroDamagePlayer.HeroID)

	return mostHeroKills, hero, mostHeroDamagePlayer
}

func (m *Match) GetMostHeroLastHits() (int, Hero, Player) {

	var MostHeroLastHits = 0
	var mostHeroDamagePlayer = Player{}

	for i, player := range m.Result.Players {
		if i == 0 || player.LastHits > MostHeroLastHits {
			MostHeroLastHits = player.LastHits
			mostHeroDamagePlayer = player
		}
	}

	hero, _ := GetHeroByID(mostHeroDamagePlayer.HeroID)

	return MostHeroLastHits, hero, mostHeroDamagePlayer
}

func (m *Match) GetMostHeroHealing() (int, Hero, Player) {

	var mostHeroHealing = 0
	var mostHeroDamagePlayer = Player{}

	for i, player := range m.Result.Players {
		if i == 0 || player.HeroHealing > mostHeroHealing {
			mostHeroHealing = player.HeroHealing
			mostHeroDamagePlayer = player
		}
	}

	hero, _ := GetHeroByID(mostHeroDamagePlayer.HeroID)

	return mostHeroHealing, hero, mostHeroDamagePlayer
}

func (m *Match) GetMostHeroDenies() (int, Hero, Player) {

	var denies = 0
	var mostHeroDamagePlayer = Player{}

	for i, player := range m.Result.Players {
		if i == 0 || player.Denies > denies {
			denies = player.Denies
			mostHeroDamagePlayer = player
		}
	}

	hero, _ := GetHeroByID(mostHeroDamagePlayer.HeroID)

	return denies, hero, mostHeroDamagePlayer
}

func (m *Match) GetMostCampsStackedHero() (int, Hero) {

	return 0, Hero{}
}

func (m *Match) GetMostHeroTowerDamage() (int, Hero, Player) {

	var damage = 0
	var mostHeroDamagePlayer = Player{}

	for i, player := range m.Result.Players {
		if i == 0 || player.TowerDamage > damage {
			damage = player.TowerDamage
			mostHeroDamagePlayer = player
		}
	}

	hero, _ := GetHeroByID(mostHeroDamagePlayer.HeroID)

	return damage, hero, mostHeroDamagePlayer
}