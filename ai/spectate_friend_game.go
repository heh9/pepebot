package ai

import (
	"net/http"
	"encoding/json"
	"github.com/pkg/errors"
)

type SpectateFriendGameResponse struct {
	Code               int         `json:"code"`
	Status             string      `json:"status"`
	Message            string      `json:"message"`
	Result             struct {
		ServerSteamId  int64       `json:"server_steam_id"`
	}                              `json:"result"`
}

func SpectateFriendGame(steamID string) (*SpectateFriendGameResponse, error) {

	emptyResponse := new(SpectateFriendGameResponse)

	client := http.Client{}
	uri := "http://192.168.1.7:9002/api/v1/spectate_friend_game/" + steamID

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return emptyResponse, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return emptyResponse, err
	}

	response := &SpectateFriendGameResponse{}

	jsonErr := json.NewDecoder(resp.Body).Decode(response)
	if jsonErr != nil {
		return emptyResponse, jsonErr
	}

	if resp.StatusCode != http.StatusOK && response.Status != "success" {
		return response, errors.New(response.Message)
	}

	return response, nil
}