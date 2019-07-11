package main

import (
	"os"
	"log"
	"net/http"
	"encoding/json"
	_ "github.com/joho/godotenv/autoload"
)

type EndStruct struct {
	Won        bool
	MatchId    string
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func main()  {

	// create the application
	application := &Application{
		MainVoiceChannelId:    os.Getenv("MAIN_VOICE_CHANNEL_ID"),
		MainTextChannelId:     os.Getenv("MAIN_TEXT_CHANNEL_ID"),
		GameEnded:             false,
		DiscordAuthToken:      os.Getenv("DISCORD_API_TOKEN"),
		VoiceChannel:          nil,
		Runes:                 NewRunesType(),
	}

	// register the discord bot and run
	application.RegisterAndServeBot()

	// Check clock time for runes
	go application.CheckRunes()

	// check game end status for stopping the bot
	go application.CheckGameEndStatus()

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")

		response := &GSIResponse{}
		json.NewDecoder(r.Body).Decode(response)

		if response.CheckAuthToken(os.Getenv("DOTA2_GSI_AUTH_TOKEN")) {
			application.GsiChannel <- response
			json.NewEncoder(w).Encode(map[string] string {
				"status": "Ok",
			})
		} else {
			json.NewEncoder(w).Encode(map[string] string {
				"status": "Failed",
			})
		}
	})

	// Listen and serve the gsi application
	http.ListenAndServe(":9001", nil)
}