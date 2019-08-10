package responses

import (
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/pkg/errors"
)

type Hero struct {
	ID             int      `json:"id"`
	LocalizedName  string   `json:"localized_name"`
	Name           string   `json:"name"`
}

var heros []Hero

func init() {

	jsonFile, err := os.Open("./api/dota2/heros.json")
	if err != nil {
		log.Println(err)
		return
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &heros)
}

func GetHeroByID(id int) (Hero, error) {

	if len(heros) == 0 {
		return Hero{}, errors.New("Could not found any hero in database!")
	}

	for _, hero := range heros {
		if hero.ID == id {
			return hero, nil
		}
	}

	return Hero{}, errors.New("Could not found the hero")
}