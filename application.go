package main

import (
	"log"
	"errors"
	"strings"
	"strconv"
	"net/http"
	"reflect"
	"pepe.bot/ai"
	"pepe.bot/disc"
	"encoding/json"
	"github.com/polds/imgbase64"
	"github.com/MrJoshLab/arrays"
	"github.com/bwmarrin/discordgo"
	"github.com/iamalirezaj/go-opendota"
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
					}

					if data.Player.TeamName != data.Map.WinTeam {
						endStruct.Won = false
						a.PlaySound("loss")
					} else {
						endStruct.Won = true
						a.PlaySound("win")
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

func (a *Application) CheckGameEndStatus() {

	for game := range a.GameEndChannel {

		if a.VoiceChannel != nil {
			a.VoiceChannel.Disconnect()
			a.VoiceChannel = nil
		}

		if game.Won {

			a.Client.ChannelMessageSend(a.MainTextChannelId,
				"```bash\n" +
					`"` + "Game ended as win with match id [" + game.MatchId + `]"` +
					"``` Weeeeee areeeee the championssssssss my friendsss " + a.GetEmoji("peepoblush").MessageFormat())
		} else {

			a.Client.ChannelMessageSend(a.MainTextChannelId,
				"```diff\n" +
					"- Game ended as loss with match id [" + game.MatchId + "]" +
					"``` - Try a bit harder next time " + a.GetEmoji("peepoblush").MessageFormat())
		}

	}

}

func (a *Application) RegisterAndServeBot() {

	discord, err := discordgo.New("Bot " + a.DiscordAuthToken)
	if err != nil {
		log.Println(err)
	}

	a.Client = discord

	a.Client.AddHandler(func (s *discordgo.Session, event *discordgo.Ready) {
		s.UpdateStatus(0, "Dota 2 [-help]")
		log.Println("Logged in as !" + event.User.ID)
	})

	a.Client.AddHandler(func (s *discordgo.Session, guild *discordgo.GuildCreate) {
		a.Emojis = guild.Emojis
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
				a.Client.ChannelMessageSend(channel.ID,
					m.Author.Mention() + " This bot is a runes reminder bot for dota 2 games that works with" +
					" Dota 2 GSI API. \n" +
					"Isn't that cool ? " + a.GetEmoji("peepoblush").MessageFormat())
				break
			case "loss_pattern":
				a.Client.ChannelMessageSend(channel.ID,
					"```diff\n" +
						"- Game ended as loss with match id [match_id]" +
						"``` - Try a bit harder next time " + a.GetEmoji("peepoblush").MessageFormat())
				break
			case "win_pattern":
				a.Client.ChannelMessageSend(channel.ID,
					"```bash\n" +
						`"` + "Game ended as win with match id [match_id]" + `"` +
						"``` Weeeeee areeeee the championssssssss my friendsss " + a.GetEmoji("peepoblush").MessageFormat())
				break
			case "leave":
				if a.VoiceChannel != nil {
					a.VoiceChannel.Disconnect()
					a.VoiceChannel = nil
				}
				break
			case "gm":

				s.ChannelTyping(channel.ID)

				player, err := FindSteamPlayer(m.Author.ID)
				if err != nil {
					a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " " + err.Error() + " " + a.GetEmoji("peepoblush").MessageFormat())
					return
				}

				message, _ := a.Client.ChannelMessageSend(channel.ID, m.Author.Mention() + " Looking for a live match ...")

				var response *ai.SpectateFriendGameResponse
				response, err = ai.SpectateFriendGame(player.SteamID)
				if err != nil || response.Code == 0 {

					// try again
					response, err = ai.SpectateFriendGame("76561198372283905")
					if err != nil || response.Code == 0 {
						a.Client.ChannelMessageEdit(channel.ID, message.ID,
							m.Author.Mention() + " You're not in a live match! " + a.GetEmoji("peepoblush").MessageFormat())
						return
					}
				}

				steamServerID := strconv.Itoa(int(response.Result.ServerSteamId))

				resp, err := ai.GetRealTimeStats(steamServerID)
				if err != nil {
					resp, err = ai.GetRealTimeStats(steamServerID)
					if err != nil {
						a.Client.ChannelMessageEdit(channel.ID, message.ID,
								m.Author.Mention() + " You're not in a live match! " + a.GetEmoji("peepoblush").MessageFormat())
							return
					}
				}

				radiant := resp.Teams[0]
				dire := resp.Teams[1]

				players := map[string] []ai.Player {
					"radiant": radiant.Players,
					"dire": dire.Players,
				}

				matchID := strconv.Itoa(int(resp.Match.MatchID))

				msg := " **Found a game with match id** `" + matchID + "`\n"
				message, _ = discord.ChannelMessageEdit(channel.ID, message.ID, msg + "Loading ... please wait")

				s.ChannelTyping(channel.ID)

				msg += "\n***Suggest to BAN :octagonal_sign:*** \n"

				var heros = map[int] opendota.Hero{}

				for _, radPl := range players["radiant"] {
					hero, err := ai.GetMostHeroPlayed(radPl.AccountID)
					if err == nil {
						reflect.ValueOf(heros).SetMapIndex(reflect.ValueOf(hero.ID), reflect.ValueOf(hero))
					}
				}

				for _, hero := range heros {
					msg += "`" + hero.LocalizedName + "` \n"
				}

				if players["dire"] == nil && players["radiant"] == nil {
					discord.ChannelMessageEdit(channel.ID, message.ID, m.Author.Mention() +
						"No players found " + a.GetEmoji("peepoblush").MessageFormat())
					break
				}

				msg += "\n***Dire Players :video_game: *** \n"
				for _, direPlayer := range players["dire"] {
					pl, err := ai.GetPlayerOpenDotaProfile(direPlayer.AccountID)
					if err != nil || pl.Rank == "" {
						msg += "`" + direPlayer.Name + " | Unknown` \n"
					} else {
						msg += "`" + direPlayer.Name + " | " + pl.Rank + "` \n"
					}
				}

				msg += "\n***Radiant Players :video_game: *** \n"
				for _, radiantPlayer := range players["radiant"] {
					pl, err := ai.GetPlayerOpenDotaProfile(radiantPlayer.AccountID)
					if err != nil || pl.Rank == "" {
						msg += "`" + radiantPlayer.Name + " | Unknown` \n"
					} else {
						msg += "`" + radiantPlayer.Name + " | " + pl.Rank + "` \n"
					}
				}

				discord.ChannelMessageEdit(channel.ID, message.ID, msg)
				break
			case "join":

				if a.VoiceChannel != nil {
					a.VoiceChannel.Disconnect()
					a.VoiceChannel = nil
				}

				for _, vs := range g.VoiceStates {

					if vs.UserID == m.Author.ID {

						var ch *discordgo.Channel

						voicech := disc.Channel{ID: vs.ChannelID, Client: a.Client}
						_, _, ch, a.VoiceChannel = voicech.Join()

						a.Client.ChannelMessageSend(channel.ID, m.Author.Mention() + " :white_check_mark: Bot successfully connected to " + ch.Name)
						return
					}
				}

				if a.VoiceChannel == nil {

					a.Client.ChannelMessageSend(
						channel.ID,
						m.Author.Mention() + " You must be connected to a voice channel! \n" +
							"Connect to a voice channel and try -join!")
					return
				}
				break
			case "play":

				if len(command) == 2 {

					if coll := collection.New(a.Runes.Sounds); !coll.Has(command[1]) {
						a.Client.ChannelMessageSend(
							channel.ID,
							m.Author.Mention() + " Could not found sound " + command[1] + " " + a.GetEmoji("peepoblush").MessageFormat())
						return
					}

				} else {

					a.Client.ChannelMessageSend(
						channel.ID,
						m.Author.Mention() + " No sound argument found " + a.GetEmoji("peepoblush").MessageFormat() + "\n" +
							"***Try with an argument like***  `-play [sound_name]`")
					return
				}

				if !a.PlaySound(a.Runes.GetRandomVoiceFileName()) {

					a.Client.ChannelMessageSend(
						channel.ID,
						m.Author.Mention() + " Bot is not connected to a voice channel. \n" +
							"Try -join to join a voice channel that you're in.")
				}
				break
			}
		}
	})

	// Open the websocket and begin listening.
	err = a.Client.Open()
	if err != nil {
		log.Println("Error opening Discord session: ", err)
	}

	log.Println("Pepe.bot is now running!")
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
		json.NewDecoder(r.Body).Decode(response)

		if response.CheckAuthToken(a.GSIAuthToken) {
			a.GsiChannel <- response
			json.NewEncoder(w).Encode(map[string] interface{} {
				"code":     200,
				"status":   "success",
			})
		} else {
			json.NewEncoder(w).Encode(map[string] interface{} {
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
	http.ListenAndServe(":" + a.GSIHttpPort, nil)
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