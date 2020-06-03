package main

import (
	"encoding/json"
	"errors"
	"github.com/MrJoshLab/arrays"
	"github.com/MrJoshLab/pepe.bot/api"
	"github.com/MrJoshLab/pepe.bot/disc"
	"github.com/bwmarrin/discordgo"
	"github.com/polds/imgbase64"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type GameEndChannel struct {
	Won        bool
	MatchId    string
}

type Application struct {

	// Discord client
	Client               *discordgo.Session

	// Discord guild created
	Guild                *discordgo.Guild

	// Discord server emojis
	Emojis               []*discordgo.Emoji

	// gsi chan
	GsiChannel           chan *GSIResponse

	// game ended chan
	GameEndChannel       chan *GameEndChannel

	// game end status
	GameEnded            bool

	// discord target voice channel
	VoiceChannel         *discordgo.VoiceConnection

	TimerChannel         chan *Timer
	Timers               []*Timer

	DiscordAuthToken,
	MainTextChannelId,
	MainVoiceChannelId,
	GSIAuthToken,
	GSIHttpPort          string

	// Runes struct
	Runes                *Runes
}

func (a *Application) CheckRunes()  {

	for data := range a.GsiChannel {

		switch data.Map.GameState {
		case PreGame:
			if a.VoiceChannel == nil {
				channel := disc.Channel{ ID: a.MainVoiceChannelId, Client: a.Client }
				_, _, _, a.VoiceChannel = channel.Join()
			}
			break
		case StrategyTime: break
		case HeroSelection: break
		case WaitForMapToLoad: break
		case WaitForPlayersToLoad: break
		case InProgress:
			if a.GameEnded {
				a.GameEnded = false
			}
			a.Runes.ClockTime = strconv.Itoa(data.Map.ClockTime)
			if ok, clock := a.Runes.Up(); ok {
				if coll := collection.New(a.Runes.RuneTimes); !coll.Has(clock) {
					a.Runes.RuneTimes = append(a.Runes.RuneTimes, clock)
					a.PlaySound(a.Runes.GetRandomVoiceFileName())
				}
			}
			break
		case PostGame:
			if !a.GameEnded {
				if data.Map.WinTeam != "none" && data.Map.WinTeam != "" {

					endStruct := &GameEndChannel{
						MatchId: data.Map.Matchid,
						Won:     false,
					}

					if data.Player.TeamName != data.Map.WinTeam {
						a.PlaySound("loss")
					} else {
						endStruct.Won = true
						a.PlaySound(a.getRandomWinSound())
					}

					a.Runes.RuneTimes = []string {}
					a.Runes.ClockTime = ""
					a.GameEnded = true
					a.GameEndChannel <- endStruct
				}
			}
			break
		}
	}

}

func (a *Application) getRandomWinSound() string {
	return []string{"win", "gta"}[RNG(0, 1)]
}

func (a *Application) CheckGameEndStatus() {

	for game := range a.GameEndChannel {

		if a.VoiceChannel != nil {
			_ = a.VoiceChannel.Disconnect()
			a.VoiceChannel = nil
		}

		msg, _ := api.GetMatchHistory(game.MatchId, true, game.Won,
			true, a.GetEmoji("peepoblush"), a.Client)

		_, _ = a.Client.ChannelMessageSend(a.MainTextChannelId, msg)
	}

}

func (a *Application) FindTimerWithMessageID(messageID string) *Timer {
	for _, timer := range a.Timers {
		if timer.MessageReaction.ID == messageID {
			return timer
		}
	}
	return nil
}

func (a *Application) RegisterAndServeBot() {

	discord, err := discordgo.New("Bot " + a.DiscordAuthToken)
	if err != nil {
		log.Println(err)
	}

	a.Client = discord

	a.Client.AddHandler(func(s *discordgo.Session, event *discordgo.Disconnect) {
		if a.VoiceChannel != nil {
			_ = a.VoiceChannel.Disconnect()
			a.VoiceChannel = nil
		}

		log.Println("disconnected from discord!")
	})

	a.Client.AddHandler(func (s *discordgo.Session, event *discordgo.Ready) {
		log.Println("Logged in as !" + event.User.ID)
	})

	a.Client.AddHandler(func (s *discordgo.Session, event *discordgo.Connect) {
		_ = s.UpdateStatus(0, "Dota 2")
	})

	a.Client.AddHandler(func (s *discordgo.Session, guild *discordgo.GuildCreate) {
		if guild.ID == "415566764697583628" {
			a.Emojis = guild.Emojis
		}
	})

	a.Client.AddHandler(func (s *discordgo.Session, m *discordgo.MessageCreate) {

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
				_, _ = a.Client.ChannelMessageSend(channel.ID,
					m.Author.Mention() + " This bot is a runes reminder bot for dota 2 games that works with" +
					" Dota 2 GSI API. \n" +
					"Isn't that cool ? " + a.GetEmoji("peepoblush").MessageFormat())
				break
			case "leave":
				if a.VoiceChannel != nil {
					_ = a.VoiceChannel.Disconnect()
					a.VoiceChannel = nil
				}
				break
			case "mh":

				if len(command) == 2 {

					var matchID = command[1]

					if _, err := strconv.ParseInt(matchID,10,64); err != nil {
						_, _ = a.Client.ChannelMessageSend(channel.ID, m.Author.Mention() +
							" The match id should be a number " +
							a.GetEmoji("peepoblush").MessageFormat())
						return
					}

					_ = s.ChannelTyping(channel.ID)

					message, _ := a.Client.ChannelMessageSend(channel.ID, m.Author.Mention() + " Looking for a match ...")

					msg, err := api.GetMatchHistory(
						matchID,
						false,
						false,
						false, a.GetEmoji("peepoblush"), a.Client)
					if err != nil {
						_, _ = a.Client.ChannelMessageEdit(channel.ID, message.ID, m.Author.Mention() +
							" " + err.Error() +
							" " + a.GetEmoji("peepoblush").MessageFormat())
						return
					}

					_, _ = a.Client.ChannelMessageEdit(channel.ID, message.ID, m.Author.Mention() + " " + msg)
					return

				} else {

					_, _ = a.Client.ChannelMessageSend(
						channel.ID,
						m.Author.Mention() + " No match_id argument found " + a.GetEmoji("peepoblush").MessageFormat() + "\n" +
							"***Try with an argument like***  `-mh [match_id]`")
					return
				}
			case "join":

				if a.VoiceChannel != nil {

					_, _ = a.Client.ChannelMessageSend(
						channel.ID,
						m.Author.Mention() + " The session already connected to a voice channel!")
					return
				}

				a.JoinChannel(channel, m, g)
			}
		}
	})

	// Open the websocket and begin listening.
	if err = a.Client.Open(); err != nil {
		log.Println("Error opening Discord session: ", err)
	}

	log.Println("Pepe.bot is now running!")
}

func (a *Application) JoinChannel(channel *discordgo.Channel, m *discordgo.MessageCreate, g *discordgo.Guild)  {

	for _, vs := range g.VoiceStates {

		if vs.UserID == m.Author.ID {

			var ch *discordgo.Channel

			voicech := disc.Channel{ID: vs.ChannelID, Client: a.Client}
			_, _, ch, a.VoiceChannel = voicech.Join()

			_, _ = a.Client.ChannelMessageSend(channel.ID, m.Author.Mention() + " :white_check_mark: Bot successfully connected to " + ch.Name)
			return
		}
	}

	if a.VoiceChannel == nil {

		_, _ = a.Client.ChannelMessageSend(
			channel.ID,
			m.Author.Mention() + " You must be connected to a voice channel! \n" +
				"Connect to a voice channel and try -join!")
		return
	}
}

func (a *Application) GetEmoji(name string) *discordgo.Emoji {

	for _, emoji := range a.Emojis {
		if emoji.Name == name {
			return emoji
		}
	}

	return nil
}

func (a *Application) ListenAndServeGSIHttpServer()  {

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")

		response := &GSIResponse{}
		_ = json.NewDecoder(r.Body).Decode(response)

		if response.CheckAuthToken(a.GSIAuthToken) {
			a.GsiChannel <- response
			_ = json.NewEncoder(w).Encode(map[string] interface{} {
				"code":     200,
				"status":   "success",
			})
		} else {
			_ = json.NewEncoder(w).Encode(map[string] interface{} {
				"code":     0,
				"status":   "failed",
			})
		}
	})

	if a.GSIHttpPort == "" {
		a.GSIHttpPort = "9001"
	}

	log.Println("Dota 2 GSI Http server running! on :" + a.GSIHttpPort)

	// Listen and serve the gsi application
	log.Println(http.ListenAndServe(":" + a.GSIHttpPort, nil))
}

func (a *Application) CreateEmojiIfNotExists(emojiName, filename string) (emoji *discordgo.Emoji, err error) {

	if a.GetEmoji(emojiName) == nil {
		err = errors.New("Emoji {" + emojiName + "} exist!")
		return
	}

	img, err := imgbase64.FromLocal(filename)
	if err != nil {
		return
	}

	emoji, err = a.Client.GuildEmojiCreate(a.Guild.ID, emojiName, img, []string {})
	return
}