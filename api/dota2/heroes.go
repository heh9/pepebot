package dota2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mrjosh/pepebot/api/dota2/responses"
	"github.com/mrjosh/pepebot/config"
)

func GetHeroes() (*responses.HeroesResponse, error) {

	url := fmt.Sprintf(
		"https://api.steampowered.com/IEconDOTA2_570/GetHeroes/v0001/?key=%s&language=en_us&format=JSON",
		config.Map.Steam.WebApiToken,
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	response := &responses.HeroesResponse{}
	if jsonErr := json.Unmarshal(body, response); jsonErr != nil {
		return nil, jsonErr
	}

	return response, nil
}
