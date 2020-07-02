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
	cmap "github.com/orcaman/concurrent-map"
	"github.com/polds/imgbase64"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math/rand"
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

type Application struct {

	// Discord client
	Client               *discordgo.Session

	// gsi chan
	GsiChannel           chan *components.GSIResponse

	// game ended chan
	GameEndChannel       chan *GameEndChannel

	// discord target voice channel
	GuildLiveMatches     cmap.ConcurrentMap

	DiscordAuthToken,
	GSIHttpPort          string

}

func (a *Application) GetGuildMatch(guildID string) (*components.GuildMatch, bool) {
	guildMatch, ok := a.GuildLiveMatches.Get(guildID)
	if ok {
		return guildMatch.(*components.GuildMatch), ok
	}
	return nil, false
}

func (a *Application) ConnectToAuthorVoiceChannel(dg *models.Guild, msg *discordgo.MessageCreate) error {
	_, ok := a.GetGuildMatch(msg.GuildID)
	if !ok {

		guild, _ := a.Client.Guild(msg.GuildID)
		for _, vs := range guild.VoiceStates {
			if vs.UserID == msg.Author.ID {

				vChannel := &disc.Channel{
					ID: vs.ChannelID,
					Client: a.Client,
				}
				voiceConnection, err := vChannel.Join()
				if err != nil {
					return err
				}

				a.GuildLiveMatches.Set(msg.GuildID, &components.GuildMatch{
					VoiceConnection: voiceConnection,
					Guild:           dg,
					DiscordGuild:    guild,
					Runes:           components.NewRunes(),
				})

			}
		}
	}
	return nil
}

func (a *Application) ConnectToVoiceChannelIfNotConnected(gameMatch *components.GSIResponse) error {
	_, ok := a.GetGuildMatch(gameMatch.DiscordGuild.ID)
	if !ok {
		// connect to main voice channel of the guild
		if gameMatch.Guild.MainVoiceChannelID != "" {
			vChannel := &disc.Channel{
				ID: gameMatch.Guild.MainVoiceChannelID,
				Client: a.Client,
			}
			voiceConnection, err := vChannel.Join()
			if err != nil {
				return err
			}

			a.GuildLiveMatches.Set(gameMatch.DiscordGuild.ID, &components.GuildMatch{
				VoiceConnection: voiceConnection,
				Guild:           gameMatch.Guild,
				DiscordGuild:    gameMatch.DiscordGuild,
				Runes:           components.NewRunes(),
			})
		}
	}
	return nil
}

func (a *Application) CheckRunes() {

	for {
		select {
		case gameMatch := <-a.GsiChannel:

			switch gameMatch.Map.GameState {

			case components.PreGame:
				if err := a.ConnectToVoiceChannelIfNotConnected(gameMatch); err != nil {
					log.Println(err)
					continue
				}
				break

			case components.StrategyTime: break
			case components.HeroSelection: break
			case components.WaitForMapToLoad: break
			case components.WaitForPlayersToLoad: break
			case components.InProgress:

				gm, ok := a.GetGuildMatch(gameMatch.DiscordGuild.ID)
				if ok {

					if gm.GameEnded {
						gm.GameEnded = false
					}

					gm.Runes.ClockTime = strconv.Itoa(gameMatch.Map.ClockTime)
					if ok, clock := gm.Runes.Up(); ok {
						if coll := collection.New(gm.Runes.RuneTimes); !coll.Has(clock) {
							gm.Runes.RuneTimes = append(gm.Runes.RuneTimes, clock)
							gm.PlaySound(gm.Runes.GetRandomVoiceFileName())
						}
					}

				}

				break
			case components.PostGame:

				gm, ok := a.GetGuildMatch(gameMatch.DiscordGuild.ID)
				if ok {

					if !gm.GameEnded {
						if gameMatch.Map.WinTeam != "none" && gameMatch.Map.WinTeam != "" {

							endStruct := &GameEndChannel{
								MatchId: gameMatch.Map.Matchid,
								Won:     false,
								GuildId: gameMatch.DiscordGuild.ID,
							}

							if gameMatch.Player.TeamName != gameMatch.Map.WinTeam {
								gm.PlaySound("loss")
							} else {
								endStruct.Won = true
								gm.PlaySound(a.getRandomWinSound())
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

}

func (a *Application) getRandomWinSound() string {
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)
	return []string{"gta", "win"}[random.Intn(2)]
}

func (a *Application) CheckGameEndStatus() {
	for {
		select {
		case gameMatch := <- a.GameEndChannel:

			gm, ok := a.GetGuildMatch(gameMatch.GuildId)
			if ok {

				if gm.HasVoiceConnection() {
					_ = gm.VoiceConnection.Disconnect()
				}

				msg, _ := api.GetMatchHistory(gameMatch.MatchId, true, gameMatch.Won, true, a.Client, gm.DiscordGuild)
				_, _ = a.Client.ChannelMessageSend(gm.Guild.MainTextChannelID, msg)

				a.GuildLiveMatches.Remove(gm.DiscordGuild.ID)

			}

		}
	}
}

func (a *Application) ChannelMessageSend(channelId, content string) {
	_, _ = a.Client.ChannelMessageSend(channelId, content)
}

func (a *Application) SendInstructions(guild *models.Guild, userId string) (err error) {

	channel, err := a.Client.UserChannelCreate(userId)
	if err != nil {
		return err
	}

	instructions :=
		"\n" +
		"Dota2 GSI Api installation" +
		"\n\n" +

		"Find your dota2 local directory, Go to `Steam` and right click on `Dota2` then: \n" +
		"> `Properties > Local Files > Browse Local Files` \n\n" +

		"Then go to `game/dota2/cfg` and create a directory called `gamestate_integration`\n" +
		"go to `gamestate_integration` and create a file called `gamestate_integration.cfg`\n" +
		"and paste the below content into it! \n\n```" +

		`"dota2-gsi Configuration"` +
		"\n{\n" +
		`    "uri"               "https://pepebot.irgeek.ir/"` + "\n" +
		`    "timeout"           "5.0"` + "\n" +
		`    "buffer"            "0.1"` + "\n" +
		`    "throttle"          "0.1"` + "\n" +
		`    "heartbeat"         "30.0"` + "\n" +
		`    "data"` + "\n" +
		"    {" + "\n" +
		`        "provider"      "1"` + "\n" +
		`        "map"           "1"` + "\n" +
		`        "player"        "1"` + "\n" +
		"    }" + "\n" +
		`    "auth"` + "\n" +
		"    {" + "\n" +
		`         "token"         "` + guild.Token + `"` + "\n" +
		"    }" + "\n" +
		"}```" +

		"\n Restart your game and You're ready to go find some matches! \n" +
		"I will connect to main_voice_channel which is `" + guild.MainVoiceChannelID + "` in your server \n" +
		"and remind you the runes every 5 minutes :sunglasses: ! \n\n" +

		"Give us some feedback or write your issues here > https://github.com/MrJoshLab/pepe.bot/issues :heart:"

	_, err = a.Client.ChannelMessageSend(channel.ID, instructions)
	return
}

func (a *Application) RegisterAndServeBot() {

	discord, err := discordgo.New("Bot " + a.DiscordAuthToken)
	if err != nil {
		log.Println(err)
	}

	a.Client = discord

	a.Client.AddHandler(func(s *discordgo.Session, event *discordgo.Disconnect) {
		log.Println("Disconnected from discord!")
	})

	a.Client.AddHandler(func (s *discordgo.Session, event *discordgo.Ready) {
		log.Println("Logged in as !" + event.User.ID)
	})

	a.Client.AddHandler(func (s *discordgo.Session, event *discordgo.Connect) {
		log.Println("Connected to discord!")
		_ = s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Game: &discordgo.Game{
				Name: "Dota 2",
				Type: discordgo.GameTypeGame,
			},
		})
	})

	a.Client.AddHandler(func (s *discordgo.Session, event *discordgo.GuildCreate) {

		var (
			guildModel = new(models.Guild)
			coll       = db.Connection.Collection("guilds")
			mCtx, _    = context.WithTimeout(context.Background(), 10 * time.Second)
			result     = coll.FindOne(mCtx, bson.M{ "discord_id": event.ID })
		)

		if err := result.Decode(guildModel); err != nil {
			guild := bson.M{
				"name": event.Name,
				"discord_id": event.ID,
				"user_id": event.OwnerID,
				"deleted": false,
				"main_voice_channel_id": nil,
				"main_text_channel_id": event.SystemChannelID,
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

			dbGuild := new(models.Guild)
			result  := db.Connection.Collection("guilds").FindOne(mCtx, bson.M{
				"discord_id": m.GuildID,
				"user_id": m.Author.ID,
			})

			if err := result.Decode(dbGuild); err != nil {
				// Could not find guild from database.
				return
			}

			prefixCommand := strings.Split(strings.TrimSpace(m.Content), "-")

			command := strings.Split(prefixCommand[1], " ")

			switch command[0] {
			case "help":

				helpText := "\n"

				if dbGuild.MainTextChannelID == "" {
					helpText += "Set MainTextChannel with `-main_text_channel {channel_id}`"
				} else {
					mainTextChannel, err := a.Client.Channel(dbGuild.MainTextChannelID)
					if err == nil {
						helpText +=
							"MainTextChannel is " + mainTextChannel.Mention() + " \n" +
								"If you want to change this you can use: `-main_text_channel {channel_id}`\n\n"
					} else {
						helpText += "Set MainTextChannel with `-main_text_channel {channel_id}`"
					}
				}

				if dbGuild.MainVoiceChannelID == "" {
					helpText += "Set MainVoiceChannel with `-main_voice_channel {voice_channel_id}`\n"
				} else {
					mainVoiceChannel, err := a.Client.Channel(dbGuild.MainVoiceChannelID)
					if err == nil {
						helpText +=
							"MainVoiceChannel is " + mainVoiceChannel.Mention() + " \n" +
								"If you want to change this you can use: `-main_voice_channel {channel_id}`\n\n"
					} else {
						helpText += "Set MainVoiceChannel with `-main_voice_channel {channel_id}`\n"
					}
				}

				helpText +=

					"Connect a dota2 player `-pa @mention_user [dota2_friend_id]`\n" +
					"Disconnect a dota2 player  `-pr @mention_user`\n" +
					"Get a summary of a dota2 match `-mh [match_id]`\n" +
					"Join voice channel that youre in `-join`\n" +
					"Disconnect from voice channel that bot connected to `-dc`\n\n " +

					"** If you having some issues with disconnecting the bot manually\n" +
					"  you should use `-dc` command to disconnect it! ** \n\n"+

					"This bot is a runes reminder bot for dota 2 games that works with" +
						" Dota 2 GSI API.\n" +
						"You can get install instruction with command: `-instruction` \n" +
						"Isn't that cool ? "

				_, _ = a.Client.ChannelMessageSend(channel.ID, helpText)
				break

			case "about":
				_, _ = a.Client.ChannelMessageSend(channel.ID,
					m.Author.Mention() + " This bot is a runes reminder bot for dota 2 games that works with" +
						" Dota 2 GSI API. \n" +
						"Isn't that cool ? ")
				break

			case "instructions":
				if m.Author.ID != guild.OwnerID {
					a.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " Only the owner of the server can get instruction!")
					return
				}
				if dbGuild.MainVoiceChannelID == "" {
					a.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " First of all, you need to set your main_voice_channel. \n" +
						"Use `-main_voice_channel {channel_id}` and then ask me for instructions!")
					return
				}
				if err := a.SendInstructions(dbGuild, m.Author.ID); err != nil {
					a.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " Could not send the instructions at the time. Please try again later!")
					return
				}
				a.ChannelMessageSend(channel.ID,
					m.Author.Mention() + " The instructions sent to your private chat successfully!")
				return

			case "disconnect", "dc", "leave":
				gm, ok := a.GetGuildMatch(dbGuild.DiscordID)
				if ok {
					if gm.HasVoiceConnection() {
						_ = gm.VoiceConnection.Disconnect()
						a.GuildLiveMatches.Remove(gm.DiscordGuild.ID)
						_, _ = a.Client.ChannelMessageSend(channel.ID,
							m.Author.Mention() + " Bot disconnected from voice channel!")
					}
				}
				break

			case "main_text_channel":

				if len(command) == 2 {

					if m.Author.ID != guild.OwnerID {
						_, _ = a.Client.ChannelMessageSend(channel.ID,
							m.Author.Mention() + " Only the owner of the server can change main_text_channel!")
						return
					}

					mainTextChannel, err := a.Client.Channel(command[1])
					if err != nil {
						_, _ = a.Client.ChannelMessageSend(channel.ID,
							m.Author.Mention() + " Could not find any channels with id: `" + command[1] + "`!")
						return
					}

					guildsCollection := db.Connection.Collection("guilds")

					_, err = guildsCollection.UpdateOne(mCtx, bson.M{"discord_id": dbGuild.DiscordID}, bson.M{
						"$set": bson.M{
							"main_text_channel_id": mainTextChannel.ID,
						},
					})

					if err != nil {
						log.Println(err)
						_, _ = a.Client.ChannelMessageSend(channel.ID,
							m.Author.Mention() + " Could not change MainTextChannel, Please try again later!")
						return
					}

					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " MainTextChannel changed to " + mainTextChannel.Mention() + " successfully!")
					return
				}

				mainTextChannel, err := a.Client.Channel(dbGuild.MainTextChannelID)
				if err != nil {
					log.Println(err)
					return
				}

				_, _ = a.Client.ChannelMessageSend(channel.ID,
					m.Author.Mention() + " The main text channel for matches is " + mainTextChannel.Mention() + " \n" +
					"If you want to change it, you can use `-main_text_channel {channel_id}`")

				break

			case "main_voice_channel":

				if len(command) == 2 {

					if m.Author.ID != dbGuild.UserID {
						_, _ = a.Client.ChannelMessageSend(channel.ID,
							m.Author.Mention() + " Only the owner of the server can change main_voice_channel!")
						return
					}

					mainTextChannel, err := a.Client.Channel(command[1])
					if err != nil {
						_, _ = a.Client.ChannelMessageSend(channel.ID,
							m.Author.Mention() + " Could not find any voice channels with id: `" + command[1] + "`!")
						return
					}

					guildsCollection := db.Connection.Collection("guilds")

					_, err = guildsCollection.UpdateOne(mCtx, bson.M{"discord_id": dbGuild.DiscordID}, bson.M{
						"$set": bson.M{
							"main_voice_channel_id": mainTextChannel.ID,
						},
					})

					if err != nil {
						log.Println(err)
						_, _ = a.Client.ChannelMessageSend(channel.ID,
							m.Author.Mention() + " Could not change MainVoiceChannel, Please try again later!")
						return
					}

					_, _ = a.Client.ChannelMessageSend(channel.ID,
						m.Author.Mention() + " MainVoiceChannel changed to " + mainTextChannel.Name + " successfully!")
					return
				}

				mainVoiceChannel, err := a.Client.Channel(dbGuild.MainVoiceChannelID)
				if err != nil {
					log.Println(err)
					return
				}

				_, _ = a.Client.ChannelMessageSend(channel.ID,
					m.Author.Mention() + " The main voice channel for matches is " + mainVoiceChannel.Mention() + " \n" +
						"If you want to change it, you can use `-main_voice_channel {voice_channel_id}`")

				break

			case "join":
				if err := a.ConnectToAuthorVoiceChannel(dbGuild, m); err != nil {
					log.Println(err)
				}
				break

			case "pr":

				playersCollection := db.Connection.Collection("players")

				if dbGuild.UserID != m.Author.ID {
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

				playersCollection := db.Connection.Collection("players")

				if dbGuild.UserID != m.Author.ID {
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

func (a *Application) ListenAndServeGSIHttpServer()  {

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")

		var (
			response = new(components.GSIResponse)
			_ = json.NewDecoder(r.Body).Decode(response)
			guild = new(models.Guild)
			coll = db.Connection.Collection("guilds")
			authToken = response.GetAuthToken()
		)

		if response.Provider.Appid != 570 {
			_ = json.NewEncoder(w).Encode(map[string] interface{} {
				"code":     0,
				"status":   "failed",
			})
			return
		}

		if authToken == "" {
			_ = json.NewEncoder(w).Encode(map[string] interface{} {
				"code":     0,
				"status":   "failed",
			})
			return
		}

		mCtx, _ := context.WithTimeout(r.Context(), 10 * time.Second)

		result := coll.FindOne(mCtx, bson.M{ "token": response.GetAuthToken() })

		if err := result.Err(); err != nil {
			_ = json.NewEncoder(w).Encode(map[string] interface{} {
				"code":     0,
				"status":   "failed",
			})
			return
		}

		if err := result.Decode(guild); err != nil {
			log.Println(err)
			_ = json.NewEncoder(w).Encode(map[string] interface{} {
				"code":     0,
				"status":   "failed",
			})
			return
		}

		discordGuild, err := a.Client.Guild(guild.DiscordID)
		if err != nil {
			log.Println(err)
			_ = json.NewEncoder(w).Encode(map[string] interface{} {
				"code":     0,
				"status":   "failed",
			})
			return
		}

		response.DiscordGuild = discordGuild
		response.Guild = guild

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