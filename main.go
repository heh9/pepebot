package main

import (
	"os"
	"log"
	"strings"
	"./disc"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/polds/imgbase64"
	"github.com/MrJoshLab/arrays"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

var (
	gameEnded = false
	runes = NewRunesType()
	pepeEmoji *discordgo.Emoji
	gsiChannel = make(chan *Response)
	voiceChannel *discordgo.VoiceConnection
	gameEndedChannel = make(chan *EndStruct)
	token = os.Getenv("DISCORD_API_TOKEN")
	MainTextChannelId = os.Getenv("MAIN_TEXT_CHANNEL_ID")
	MainVoiceChannelId = os.Getenv("MAIN_VOICE_CHANNEL_ID")
)

type EndStruct struct {
	Won        bool
	MatchId    string
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func main()  {

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")

		response := &Response{}
		json.NewDecoder(r.Body).Decode(response)

		if response.CheckAuthToken(os.Getenv("DOTA2_GSI_AUTH_TOKEN")) {
			gsiChannel <- response
			json.NewEncoder(w).Encode(map[string] string {
				"status": "Ok",
			})
		} else {
			json.NewEncoder(w).Encode(map[string] string {
				"status": "Failed",
			})
		}
	})

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println(err)
	}

	discord.AddHandler(func (s *discordgo.Session, event *discordgo.Ready) {
		s.UpdateStatus(0, "Dota 2 [-help]")
		log.Println("Logged in as !" + event.User.ID)
	})

	discord.AddHandler(func (s *discordgo.Session, guild *discordgo.GuildCreate) {

		var exists = false

		for _, guildEmoji := range guild.Emojis {
			if guildEmoji.Name == "peepoblush" {
				pepeEmoji = guildEmoji
				exists = true
				break
			}
		}

		if !exists {
			img, err := imgbase64.FromLocal("./emojies/peepoblush.png")
			if err != nil {
				log.Println(err)
				return
			}

			pepeEmoji, err = discord.GuildEmojiCreate(
				guild.ID, "peepoblush", img, []string {})
			if err != nil {
				log.Println(err)
				return
			}
		}

	})

	discord.AddHandler(func (s *discordgo.Session, m *discordgo.MessageCreate) {

		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, "-") {

			channel, err := s.State.Channel(m.ChannelID)
			if err != nil {
				// Could not find channel.
				return
			}

			// Find the guild for that channel.
			g, err := s.State.Guild(channel.GuildID)
			if err != nil {
				// Could not find guild.
				return
			}

			prefixCommand := strings.Split(m.Content, "-")

			command := strings.Split(prefixCommand[1], " ")

			switch command[0] {
			case "help":
				discord.ChannelMessageSend(channel.ID, "This bot is a runes reminder bot for dota 2 games that works with" +
					" Dota 2 GSI API. \n" +
					"Isn't that cool ? <:" + pepeEmoji.Name + ":" + pepeEmoji.ID + ">")
			break
			case "leave":
				if voiceChannel != nil {
					voiceChannel.Disconnect()
					voiceChannel = nil
				}
				break
			case "join":

				if voiceChannel != nil {
					voiceChannel.Disconnect()
					voiceChannel = nil
				}

				for _, vs := range g.VoiceStates {

					if vs.UserID == m.Author.ID {

						var ch *discordgo.Channel

						voicech := disc.Channel{ID: vs.ChannelID, Client: discord}
						_, _, ch, voiceChannel = voicech.Join()

						discord.ChannelMessageSend(channel.ID, ":white_check_mark: Bot successfully connected to " + ch.Name)
						return
					}
				}

				if voiceChannel == nil {

					discord.ChannelMessageSend(
						channel.ID,
						"You must be connected to a voice channel! \n" +
						"Connect to a voice channel and try -join!")
					return
				}
				break
			case "play":

				if len(command) == 2 {

					if coll := collection.New(runes.Sounds); !coll.Has(command[1]) {
						discord.ChannelMessageSend(
							channel.ID,
							"Could not found sound " + command[1] + " <:" + pepeEmoji.Name + ":" + pepeEmoji.ID + ">")
						return
					}

				} else {

					discord.ChannelMessageSend(
						channel.ID,
						"No sound argument found <:" + pepeEmoji.Name + ":" + pepeEmoji.ID + ">\n" +
							"***Try with an argument like***  `-play [sound_name]`")
					return
				}

				if !playSound(runes.GetRandomVoiceFileName()) {

					discord.ChannelMessageSend(
						channel.ID,
						"Bot is not connected to a voice channel. \n" +
						"Try -join to join a voice channel that you're in.")
				}
				break
			}
		}
	})

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		log.Println("Error opening Discord session: ", err)
	}

	log.Println("Pepe.bot is now running.  Press CTRL-C to exit.")

	go func() {

		for data := range gsiChannel {

			switch data.Map.GameState {
			case PreGame:
				if voiceChannel == nil {
					channel := disc.Channel{ID: MainVoiceChannelId, Client: discord}
					_, _, _, voiceChannel = channel.Join()
				}
				break
			case StrategyTime: break
			case HeroSelection: break
			case WaitForMapToLoad: break
			case WaitForPlayersToLoad: break
			case InProgress:
				gameEnded = false
				runes.ClockTime = strconv.Itoa(data.Map.ClockTime)
				if ok, clock := runes.Up(); ok {
					if coll := collection.New(runes.RuneTimes); !coll.Has(clock) {
						runes.RuneTimes = append(runes.RuneTimes, clock)
						playSound(runes.GetRandomVoiceFileName())
					}
				}
				break
			case PostGame:
				if !gameEnded {
					if data.Map.WinTeam != "none" && data.Map.WinTeam != "" {

						endStruct := &EndStruct{
							MatchId: data.Map.Matchid,
						}

						if data.Player.TeamName != data.Map.WinTeam {
							endStruct.Won = false
							playSound("loss")
						} else {
							endStruct.Won = true
							playSound("win")
						}

						gameEnded = true
						gameEndedChannel <- endStruct
					}
				}
				break
			}
		}

	}()

	go func() {

		for game := range gameEndedChannel {

			if voiceChannel != nil {
				voiceChannel.Disconnect()
				voiceChannel = nil
			}

			var wonText = "lost"
			var StatusText = "Try a bit harder next time <:" + pepeEmoji.Name + ":" + pepeEmoji.ID + ">"

			if game.Won {
				wonText = "win"
				StatusText = "Weeeeee Areeeee the championssssssss my friendsss <:" + pepeEmoji.Name + ":" + pepeEmoji.ID + ">"
			}

			discord.ChannelMessageSend(MainTextChannelId,
				"```css\n" +
					"Game ended as " + wonText + " with match id [" + game.MatchId + "]" +
					"```" + StatusText)

		}

	}()

	// Cleanly close down the Discord session.
	defer discord.Close()

	http.ListenAndServe(":9001", nil)
}