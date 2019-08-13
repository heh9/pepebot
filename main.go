package main

import (
	"os"
	"log"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func main()  {

	// Define a new application
	application := &Application{

		// Discord config
		MainVoiceChannelId:    os.Getenv("MAIN_VOICE_CHANNEL_ID"),
		MainTextChannelId:     os.Getenv("MAIN_TEXT_CHANNEL_ID"),
		DiscordAuthToken:      os.Getenv("DISCORD_API_TOKEN"),
		VoiceChannel:          nil,
		GameEnded:             false,
		GsiChannel:            make(chan *GSIResponse),
		GameEndChannel:        make(chan *GameEndChannel),

		TimerChannel:          make(chan *Timer),

		Runes:                 NewRunes(),

		// GSI Config
		GSIAuthToken:          os.Getenv("DOTA2_GSI_AUTH_TOKEN"),
		GSIHttpPort:           os.Getenv("DOTA2_GSI_HTTP_PORT"),
	}

	// register the discord bot and run
	application.RegisterAndServeBot()

	// Close discord client when program will close
	defer application.Client.Close()

	// Check clock time for runes
	go application.CheckRunes()

	// check game end status for stopping the bot
	go application.CheckGameEndStatus()

	// Listen and serve gsi http server
	application.ListenAndServeGSIHttpServer()
}