package main

import (
	"github.com/MrJoshLab/pepe.bot/components"
	_ "github.com/joho/godotenv/autoload"
	cmap "github.com/orcaman/concurrent-map"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func main() {

	// Define a new application
	application := &Application{
		DiscordAuthToken:      os.Getenv("DISCORD_API_TOKEN"),
		GSIHttpPort:           os.Getenv("DOTA2_GSI_HTTP_PORT"),
		GsiChannel:            make(chan *components.GSIResponse),
		GameEndChannel:        make(chan *GameEndChannel),
		GuildLiveMatches:      cmap.New(),
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