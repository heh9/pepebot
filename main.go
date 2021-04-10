package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mrjosh/pepebot/components"
	"github.com/mrjosh/pepebot/config"
	cmap "github.com/orcaman/concurrent-map"
)

var (
	httpPort *int
	httpHost *string
)

func main() {

	log.SetFlags(log.Lshortfile | log.Ltime)

	configFileName := flag.String("config-file", "config.hcl", "config.hcl file")
	httpHost = flag.String("host", "0.0.0.0", "Http Host")
	httpPort = flag.Int("port", 9001, "Http Port")
	flag.Parse()

	log.Printf("Loading ConfigMap from file: [%s]", *configFileName)

	if err := config.LoadFile(*configFileName); err != nil {
		log.Fatal(fmt.Errorf("could not load config: %v", err))
	}

	// Define a new application
	application := &Application{
		GsiChannel:       make(chan *components.GSIResponse),
		GameEndChannel:   make(chan *GameEndChannel),
		GuildLiveMatches: cmap.New(),
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
	application.ListenAndServeGSIHttpServer(*httpHost, *httpPort)
}
