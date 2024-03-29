package dota2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mrjosh/pepebot/api/dota2/responses"
)

type Client struct {
	Key     string
	BaseUri string
}

func NewClient(key string) *Client {
	return &Client{
		Key:     key,
		BaseUri: "https://api.steampowered.com",
	}
}

func (c *Client) Match(matchID string) (*responses.Match, error) {

	url := fmt.Sprintf("%s/IDOTA2Match_570/GetMatchDetails/v1/?match_id=%s&key=%s", c.BaseUri, matchID, c.Key)

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

	response := &responses.Match{}
	if jsonErr := json.Unmarshal(body, response); jsonErr != nil {
		return nil, jsonErr
	}

	return response, nil
}
