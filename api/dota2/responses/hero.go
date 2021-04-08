package responses

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mrjoshlab/pepe.bot/config"
)

type Hero struct {
	ID            int    `json:"id"`
	LocalizedName string `json:"localized_name"`
	Name          string `json:"name"`
	Icon          string `json:"icon"`
}

type HeroesResponse struct {
	Result struct {
		Heroes []Hero `json:"heroes"`
		Status uint32 `json:"status"`
		Count  uint32 `json:"count"`
	} `json:"result"`
}

var (
	err    error
	heroes []Hero
)

func FeatchHeroes() error {

	url := fmt.Sprintf(
		"https://api.steampowered.com/IEconDOTA2_570/GetHeroes/v0001/?key=%s&language=en_us&format=JSON",
		config.Map.Steam.WebApiToken,
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bufferBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	response := &HeroesResponse{}
	if jsonErr := json.Unmarshal(bufferBytes, response); jsonErr != nil {
		return jsonErr
	}

	heroes = response.Result.Heroes
	return nil
}

func GetHeroByID(id int) (Hero, error) {

	if len(heroes) == 0 {
		return Hero{}, errors.New("Could not find any hero in database!")
	}

	for _, hero := range heroes {
		if hero.ID == id {
			return hero, nil
		}
	}

	return Hero{}, errors.New("Could not found the hero")
}
