package main

import (
	"context"
	"encoding/json"
	collection "github.com/MrJoshLab/arrays"
	"github.com/MrJoshLab/pepe.bot/api"
	"github.com/MrJoshLab/pepe.bot/components"
	"github.com/MrJoshLab/pepe.bot/db"
	"github.com/MrJoshLab/pepe.bot/disc"
	"github.com/MrJoshLab/pepe.bot/models"
	"github.com/bwmarrin/discordgo"
	"github.com/polds/imgbase64"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GameEndChannel struct {
	Won        bool
	MatchId    string
	GuildId    string
}

type GuildMatch struct {
	VoiceConnection *discordgo.VoiceConnection
	Guild           *discordgo.Guild
	GameEnded       bool
	Runes           *Runes
}

type Application struct {

	// Discord client
	Client               *discordgo.Session

	// gsi chan
	GsiChannel           chan *GSIResponse

	// game ended chan
	GameEndChannel       chan *GameEndChannel

	// discord target voice channel
	GuildLiveMatches     map[string] *GuildMatch

	DiscordAuthToken,
	GSIHttpPort          string
}

func (a *Application) CheckRunes()  {

	for data := range a.GsiChannel {

		gm := a.GuildLiveMatches[data.GuildID]

		switch data.Map.GameState {
		case PreGame:

			if gm == nil {

				guild, err := a.Client.Guild(data.GuildID)
				if err == nil {

					channel := disc.Channel{ ID: "", Client: a.Client }
					_, _, _, voiceConnection := channel.Join()

					a.GuildLiveMatches[data.GuildID] = &GuildMatch{
						Runes: NewRunes(),
						Guild: guild,
						VoiceConnection: voiceConnection,
					}

				}

			}

			break
		case StrategyTime: break
		case HeroSelection: break
		case WaitForMapToLoad: break
		case WaitForPlayersToLoad: break
		case InProgress:

			if gm != nil {
				if gm.GameEnded {
					gm.GameEnded = false
				}
				gm.Runes.ClockTime = strconv.Itoa(data.Map.ClockTime)
				if ok, clock := gm.Runes.Up(); ok {
					if coll := collection.New(gm.Runes.RuneTimes); !coll.Has(clock) {
						gm.Runes.RuneTimes = append(gm.Runes.RuneTimes, clock)

						// This should call from guild
						a.PlaySound(gm.Runes.GetRandomVoiceFileName())
					}
				}
			}

			break
		case PostGame:

			if gm != nil {
				if !gm.GameEnded {
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

						gm.Runes.RuneTimes = []string {}
						gm.Runes.ClockTime = ""
						gm.GameEnded = true
						a.GameEndChannel <- endStruct
					}
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

		gm := a.GuildLiveMatches[game.GuildId]

		if gm.VoiceConnection != nil {
			_ = gm.VoiceConnection.Disconnect()
			gm.VoiceConnection = nil
		}

		msg, _ := api.GetMatchHistory(game.MatchId, true, game.Won,
			true, a.Client, gm.Guild)

		_, _ = a.Client.ChannelMessageSend("MainTextChannelId", msg)
	}

}

func (a *Application) RegisterAndServeBot() {

	discord, err := discordgo.New("Bot " + a.DiscordAuthToken)
	if err != nil {
		log.Println(err)
	}

	a.Client = discord

	a.Client.AddHandler(func(s *discordgo.Session, event *discordgo.Disconnect) {
		log.Println("disconnected from discord!")
	})

	a.Client.AddHandler(func (s *discordgo.Session, event *discordgo.Ready) {
		log.Println("Logged in as !" + event.User.ID)
	})

	a.Client.AddHandler(func (s *discordgo.Session, event *discordgo.Connect) {
		_ = s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Game: &discordgo.Game{
				Name: "Dota 2",
				Type: discordgo.GameTypeGame,
			},
		})
	})

	a.Client.AddHandler(func (s *discordgo.Session, event *discordgo.GuildCreate) {

		guildModel := new(models.Guild)
		coll := db.Connection.Collection("guilds")
		mCtx, _ := context.WithTimeout(context.Background(), 10 * time.Second)

		result := coll.FindOne(mCtx, bson.M{ "discord_id": event.ID })

		if err := result.Decode(guildModel); err != nil {
			guild := bson.M{
				"name": event.Name,
				"discord_id": event.ID,
				"user_id": event.OwnerID,
				"deleted": false,
				"token": components.Random(25),
				"created_at": time.Now(),
				"deleted_at": time.Now(),
			}
			_, _ = coll.InsertOne(mCtx, guild)
		}

	})

	a.Client.AddHandler(func (s *discordgo.Session, m *discordgo.MessageCreate) {

		mCtx, _ := context.WithTimeout(context.Background(), 10 * time.Second)

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
			guild, err := s.State.Guild(channel.GuildID)
			if err != nil {
				// Could not find guild.
				return
			}

			prefixCommand := strings.Split(strings.TrimSpace(m.Content), "-")

			command := strings.Split(prefixCommand[1], " ")

			switch command[0] {
			case "help":

				_, _ = a.Client.ChannelMessageSend(channel.ID,
					m.Author.Mention() + " This bot is a runes reminder bot for dota 2 games that works with" +
					" Dota 2 GSI API. \n" +
					"Isn't that cool ? ")
				break

			case "pr":

				guildModel := new(models.Guild)
				coll := db.Connection.Collection("guilds")
				playersCollection := db.Connection.Collection("players")

				result := coll.FindOne(mCtx, bson.M{
					"discord_id": m.GuildID,
					"user_id": m.Author.ID,
				})

				if err := result.Decode(guildModel); err != nil {
					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " Only the owner of the guild can add/remove/update a player!")
					return
				}

				if len(command) < 2 {
					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " Arguments not found, Try this pattern: `-pr @mention_a_user`!")
					return
				}

				if len(m.Mentions) > 1 || m.Mentions[0].Bot {
					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " You sould mention only one member and it can not be a bot!")
					break
				}

				count, err := playersCollection.CountDocuments(mCtx, bson.M{
					"user_discord_id": m.Mentions[0].ID,
					"guild_id": m.GuildID,
				})
				if err != nil {
					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " Something's wrong, Please try again later!")
					break
				}

				if count == 0 {
					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " Player does not exists!")
					break
				}

				_, deleteErr := playersCollection.DeleteOne(mCtx, bson.M{
					"user_discord_id": m.Mentions[0].ID,
					"guild_id": m.GuildID,
				})

				if deleteErr != nil {
					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " Cannot remove the player, Please try again later!")
					break
				}

				_, _ = a.Client.ChannelMessageSend(channel.ID,
					m.Author.Mention() + " Player removed successfully!")
				break

			case "pa":

				guildModel := new(models.Guild)
				coll := db.Connection.Collection("guilds")
				playersCollection := db.Connection.Collection("players")

				result := coll.FindOne(mCtx, bson.M{
					"discord_id": m.GuildID,
					"user_id": m.Author.ID,
				})

				if err := result.Decode(guildModel); err != nil {
					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " Only the owner of the guild can add/remove/update a player!")
					return
				}

				if len(command) < 3 {
					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " Arguments not found, Try this pattern: `-pa @mention_a_user [dota2_friend_id]`!")
					return
				}

				if len(m.Mentions) > 1 || m.Mentions[0].Bot {
					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " You sould mention only one member and it can not be a bot!")
					break
				}

				count, err := playersCollection.CountDocuments(mCtx, bson.M{
					"user_discord_id": m.Mentions[0].ID,
					"guild_id": m.GuildID,
				})
				if err != nil {
					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " Something's wrong, Please try again later!")
					break
				}

				if count > 0 {
					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " Player already added!")
					break
				}

				_, insertErr := playersCollection.InsertOne(mCtx, bson.M{
					"name": m.Mentions[0].Username,
					"account_id": command[2],
					"user_discord_id": m.Mentions[0].ID,
					"guild_id": m.GuildID,
				})

				if insertErr != nil {
					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " Cannot add the player, Please try again later!")
					break
				}

				_, _ = a.Client.ChannelMessageSend(channel.ID,
					m.Author.Mention() + " Player added successfully!")
				break

			case "mh":

				if len(command) == 2 {

					var matchID = command[1]

					if _, err := strconv.ParseInt(matchID,10,64); err != nil {
						_, _ = a.Client.ChannelMessageSend(channel.ID, m.Author.Mention() +
							" The match id should be a number!")
						return
					}

					message, _ := a.Client.ChannelMessageSend(channel.ID, m.Author.Mention() + " Looking for a match ...")

					_ = s.ChannelTyping(channel.ID)

					msg, err := api.GetMatchHistory(matchID, false, false, false, a.Client, guild)
					if err != nil {
						_, _ = a.Client.ChannelMessageEdit(channel.ID, message.ID, m.Author.Mention() +
							" " + err.Error())
						return
					}

					_, _ = a.Client.ChannelMessageEdit(channel.ID, message.ID, m.Author.Mention() + " " + msg)
					return

				} else {

					_, _ = a.Client.ChannelMessageSend(
						channel.ID,
						m.Author.Mention() + " No match_id argument found " + "\n" +
							"***Try with an argument like***  `-mh [match_id]`")
					return
				}

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

	//for _, vs := range g.VoiceStates {
	//
	//	if vs.UserID == m.Author.ID {
	//
	//		var ch *discordgo.Channel
	//
	//		voicech := disc.Channel{ID: vs.ChannelID, Client: a.Client}
	//		_, _, ch, a.VoiceChannel = voicech.Join()
	//
	//		_, _ = a.Client.ChannelMessageSend(channel.ID, m.Author.Mention() + " :white_check_mark: Bot successfully connected to " + ch.Name)
	//		return
	//	}
	//}
	//
	//if a.VoiceChannel == nil {
	//
	//	_, _ = a.Client.ChannelMessageSend(
	//		channel.ID,
	//		m.Author.Mention() + " You must be connected to a voice channel! \n" +
	//			"Connect to a voice channel and try -join!")
	//	return
	//}
}

func (a *Application) ListenAndServeGSIHttpServer()  {

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")

		var (
			response = new(GSIResponse)
			_ = json.NewDecoder(r.Body).Decode(response)

			guild = new(models.Guild)
			collection = db.Connection.Collection("guilds")
			authToken = response.GetAuthToken()
		)

		if authToken == "" {
			_ = json.NewEncoder(w).Encode(map[string] interface{} {
				"code":     0,
				"status":   "failed",
			})
			return
		}

		result := collection.FindOne(context.Background(), bson.M{
			"token": response.GetAuthToken(),
		})

		switch result.Err() {
		case mongo.ErrNoDocuments:
			_ = json.NewEncoder(w).Encode(map[string] interface{} {
				"code":     0,
				"status":   "failed",
			})
			return
		}

		if err := result.Decode(guild); err != nil {
			_ = json.NewEncoder(w).Encode(map[string] interface{} {
				"code":     0,
				"status":   "failed",
			})
			return
		}

		a.GsiChannel <- response
		_ = json.NewEncoder(w).Encode(map[string] interface{} {
			"code":     200,
			"status":   "success",
		})

	})

	if a.GSIHttpPort == "" {
		a.GSIHttpPort = "9001"
	}

	log.Println("Dota 2 GSI Http server running! on :" + a.GSIHttpPort)

	// Listen and serve the gsi application
	log.Println(http.ListenAndServe(":" + a.GSIHttpPort, nil))
}

func (a *Application) CreateEmojiIfNotExists(GuildId string, emojiName, filename string) (emoji *discordgo.Emoji, err error) {

	img, err := imgbase64.FromLocal(filename)
	if err != nil {
		return
	}

	emoji, err = a.Client.GuildEmojiCreate(GuildId, emojiName, img, []string {})
	return
}