package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	collection "github.com/mrjoshlab/go-collection"
	"github.com/mrjoshlab/pepe.bot/api"
	"github.com/mrjoshlab/pepe.bot/api/dota2/responses"
	"github.com/mrjoshlab/pepe.bot/components"
	"github.com/mrjoshlab/pepe.bot/config"
	"github.com/mrjoshlab/pepe.bot/db"
	"github.com/mrjoshlab/pepe.bot/disc"
	"github.com/mrjoshlab/pepe.bot/models"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/polds/imgbase64"
)

type GameEndChannel struct {
	Won     bool
	MatchId string
	GuildId string
}

type Application struct {
	Client           *discordgo.Session
	GsiChannel       chan *components.GSIResponse
	GameEndChannel   chan *GameEndChannel
	GuildLiveMatches cmap.ConcurrentMap
}

func (a *Application) GetGuildMatch(guildID string) (*components.GuildMatch, bool) {
	guildMatch, ok := a.GuildLiveMatches.Get(guildID)
	if ok {
		return guildMatch.(*components.GuildMatch), ok
	}
	return nil, false
}

func (a *Application) ConnectToAuthorVoiceChannel(dg *models.Guild, msg *discordgo.MessageCreate) error {
	if _, ok := a.GetGuildMatch(msg.GuildID); !ok {

		guild, _ := a.Client.Guild(msg.GuildID)
		for _, vs := range guild.VoiceStates {
			if vs.UserID == msg.Author.ID {

				vChannel := &disc.Channel{
					ID:     vs.ChannelID,
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
	if _, ok := a.GetGuildMatch(gameMatch.DiscordGuild.ID); !ok {
		// connect to main voice channel of the guild
		if gameMatch.Guild.MainVoiceChannelID != "" {
			vChannel := &disc.Channel{
				ID:     gameMatch.Guild.MainVoiceChannelID,
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

			case components.StrategyTime:
				break
			case components.HeroSelection:
				break
			case components.WaitForMapToLoad:
				break
			case components.WaitForPlayersToLoad:
				break
			case components.InProgress:

				if gm, ok := a.GetGuildMatch(gameMatch.DiscordGuild.ID); ok {

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

				if gm, ok := a.GetGuildMatch(gameMatch.DiscordGuild.ID); ok {

					if !gm.GameEnded {
						if gameMatch.Map.WinTeam != "none" && gameMatch.Map.WinTeam != "" {

							endStruct := &GameEndChannel{
								MatchId: gameMatch.Map.Matchid,
								Won:     false,
								GuildId: gameMatch.DiscordGuild.ID,
							}

							if gameMatch.Player.TeamName != gameMatch.Map.WinTeam {
								gm.PlaySound(a.getRandomLossSound())
							} else {
								endStruct.Won = true
								gm.PlaySound(a.getRandomWinSound())
							}

							gm.Runes.RuneTimes = []string{}
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

func (a *Application) getRandomLossSound() string {
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)
	sounds := config.Map.Sounds.Loss
	return sounds[random.Intn(len(sounds))]
}

func (a *Application) getRandomWinSound() string {
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)
	sounds := config.Map.Sounds.Win
	return sounds[random.Intn(len(sounds))]
}

func (a *Application) CheckGameEndStatus() {
	for {
		select {
		case gameMatch := <-a.GameEndChannel:

			if gm, ok := a.GetGuildMatch(gameMatch.GuildId); ok {

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
			`    "uri"               "https://pepebot.mrjosh.net"` + "\n" +
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

			"Give us some feedback or write your issues here > https://github.com/mrjoshlab/pepe.bot/issues :heart:"

	_, err = a.Client.ChannelMessageSend(channel.ID, instructions)
	return
}

func (a *Application) RegisterAndServeBot() {

	if err := responses.FeatchHeroes(); err != nil {
		log.Panicf("Could not featch heroes from steam web api :%v", err)
	}

	if err := db.Configure(); err != nil {
		log.Panicf("Could not connect to database :%v", err)
	}

	dbModels := []interface{}{
		&models.Guild{},
		&models.Player{},
	}

	if err := db.Connection.AutoMigrate(dbModels...); err != nil {
		log.Panicf("Could not connect to database :%v", err)
	}

	discord, err := discordgo.New(fmt.Sprintf("Bot %s", config.Map.Discord.Token))
	if err != nil {
		log.Panicf("Could not connect to discord :%v", err)
	}

	a.Client = discord

	a.Client.AddHandler(func(s *discordgo.Session, event *discordgo.Disconnect) {
		log.Println("Disconnected from discord!")
	})

	a.Client.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		log.Println("Logged in as !" + event.User.ID)
	})

	a.Client.AddHandler(func(s *discordgo.Session, event *discordgo.Connect) {
		log.Println("Connected to discord!")
		_ = s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name: "Dota 2",
					Type: discordgo.ActivityTypeGame,
				},
			},
		})
	})

	a.Client.AddHandler(func(s *discordgo.Session, event *discordgo.GuildCreate) {

		var (
			guildModel = new(models.Guild)
			result     = db.Connection.Where("discord_id =?", event.ID).First(&guildModel)
		)

		if err := result.Error; err != nil {
			db.Connection.Create(&models.Guild{
				Name:               event.Name,
				DiscordID:          event.ID,
				UserID:             event.OwnerID,
				MainVoiceChannelID: "",
				MainTextChannelID:  "",
				Token:              components.Random(25),
			})
		}

	})

	a.Client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {

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

			var (
				dbGuild = new(models.Guild)
				result  = db.Connection.Where("discord_id =?", m.GuildID).Where("user_id =?", m.Author.ID).First(&dbGuild)
			)
			if err := result.Error; err != nil {
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
						"  you should use `-dc` command to disconnect it! ** \n\n" +

						"This bot is a runes reminder bot for dota 2 games that works with" +
						" Dota 2 GSI API.\n" +
						"You can get install instruction with command: `-instructions` \n" +
						"Isn't that cool ? "

				_, _ = a.Client.ChannelMessageSend(channel.ID, helpText)
				break

			case "about":
				_, _ = a.Client.ChannelMessageSend(channel.ID,
					m.Author.Mention()+" This bot is a runes reminder bot for dota 2 games that works with"+
						" Dota 2 GSI API. \n"+
						"Isn't that cool ? ")
				break

			case "instructions":
				if m.Author.ID != guild.OwnerID {
					a.ChannelMessageSend(channel.ID,
						m.Author.Mention()+" Only the owner of the server can get instruction!")
					return
				}
				if dbGuild.MainVoiceChannelID == "" {
					a.ChannelMessageSend(channel.ID,
						m.Author.Mention()+" First of all, you need to set your main_voice_channel. \n"+
							"Use `-main_voice_channel {channel_id}` and then ask me for instructions!")
					return
				}
				if err := a.SendInstructions(dbGuild, m.Author.ID); err != nil {
					a.ChannelMessageSend(channel.ID,
						m.Author.Mention()+" Could not send the instructions at the time. Please try again later!")
					return
				}
				a.ChannelMessageSend(channel.ID,
					m.Author.Mention()+" The instructions sent to your private chat successfully!")
				return

			case "disconnect", "dc", "leave":
				gm, ok := a.GetGuildMatch(dbGuild.DiscordID)
				if ok {
					if gm.HasVoiceConnection() {
						_ = gm.VoiceConnection.Disconnect()
						a.GuildLiveMatches.Remove(gm.DiscordGuild.ID)
						_, _ = a.Client.ChannelMessageSend(channel.ID,
							m.Author.Mention()+" Bot disconnected from voice channel!")
					}
				}
				break

			case "main_text_channel":

				if len(command) == 2 {

					if m.Author.ID != guild.OwnerID {
						a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
							"%s Only the owner of the server can change main_text_channel!",
							m.Author.Mention(),
						))
						return
					}

					mainTextChannel, err := a.Client.Channel(command[1])
					if err != nil {
						a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
							"%s Could not find any channels with id: `%s`!",
							m.Author.Mention(),
							command[1],
						))
						return
					}

					updateQuery := db.Connection.Model(&models.Guild{}).Where("discord_id =?", dbGuild.DiscordID).Updates(map[string]interface{}{
						"main_text_channel_id": mainTextChannel.ID,
					})

					if updateQuery.Error != nil {
						a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
							"%s Could not change MainTextChannel, Please try again later!",
							m.Author.Mention(),
						))
						return
					}

					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s MainTextChannel changed to %s successfully!",
						m.Author.Mention(),
						mainTextChannel.Mention(),
					))
					return
				}

				mainTextChannel, err := a.Client.Channel(dbGuild.MainTextChannelID)
				if err != nil {
					return
				}

				a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
					"%s The main text channel for matches is %s\n"+
						"If you want to change it, you can use `-main_text_channel {channel_id}`",
					m.Author.Mention(),
					mainTextChannel.Mention(),
				))

				break

			case "main_voice_channel":

				if len(command) == 2 {

					if m.Author.ID != dbGuild.UserID {
						a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
							"%s Only the owner of the server can change main_voice_channel!",
							m.Author.Mention(),
						))
						return
					}

					mainVoiceChannel, err := a.Client.Channel(command[1])
					if err != nil {
						a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
							"%s Could not find any voice channels with id: `%s`!",
							m.Author.Mention(),
							command[1],
						))
						return
					}

					updateQuery := db.Connection.Model(&models.Guild{}).Where("discord_id =?", dbGuild.DiscordID).Updates(map[string]interface{}{
						"main_voice_channel_id": mainVoiceChannel.ID,
					})

					if updateQuery.Error != nil {
						a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
							"%s Could not change MainVoiceChannel, Please try again later!",
							m.Author.Mention(),
						))
						return
					}

					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s MainVoiceChannel changed to %s successfully!",
						m.Author.Mention(),
						mainVoiceChannel.Name,
					))
					return
				}

				mainVoiceChannel, err := a.Client.Channel(dbGuild.MainVoiceChannelID)
				if err != nil {
					log.Println(err)
					return
				}

				a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
					"%s The main voice channel for matches is %s\n"+
						"If you want to change it, you can use `-main_voice_channel {voice_channel_id}`",
					m.Author.Mention(),
					mainVoiceChannel.Mention(),
				))

				break

			case "join":
				if err := a.ConnectToAuthorVoiceChannel(dbGuild, m); err != nil {
					log.Println(err)
				}
				break

			case "pr":

				if dbGuild.UserID != m.Author.ID {
					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s Only the owner of the guild can add/remove/update a player!",
						m.Author.Mention(),
					))
					return
				}

				if len(command) < 2 {
					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s Arguments not found, Try this pattern: `-pr @mention_a_user`!",
						m.Author.Mention(),
					))
					return
				}

				if len(m.Mentions) > 1 || m.Mentions[0].Bot {
					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s You sould mention only one member and it can not be a bot!",
						m.Author.Mention(),
					))
					break
				}

				var count int64
				result := db.Connection.Model(&models.Player{}).
					Where("user_discord_id =?", m.Mentions[0].ID).
					Where("guild_id =?", m.GuildID).
					Count(&count)
				if result.Error != nil {
					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s Something's wrong, Please try again later!",
						m.Author.Mention(),
					))
					break
				}

				if count == 0 {
					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s Player does not exists!",
						m.Author.Mention(),
					))
					break
				}

				deleteErr := db.Connection.Where("user_discord_id =?", m.Mentions[0].ID).
					Where("guild_id =?", m.GuildID).Delete(&models.Player{}).Error
				if deleteErr != nil {
					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s Cannot remove the player, Please try again later!",
						m.Author.Mention(),
					))
					break
				}

				a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
					"%s Player removed successfully!",
					m.Author.Mention(),
				))
				break

			case "pa":

				if dbGuild.UserID != m.Author.ID {
					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s Only the owner of the guild can add/remove/update a player!",
						m.Author.Mention(),
					))
					return
				}

				if len(command) < 3 {
					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s Arguments not found, Try this pattern: `-pa @mention_a_user [dota2_friend_id]`!",
						m.Author.Mention(),
					))
					return
				}

				if len(m.Mentions) > 1 || m.Mentions[0].Bot {
					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s You sould mention only one member and it can not be a bot!",
						m.Author.Mention(),
					))
					break
				}

				var count int64
				result := db.Connection.Model(&models.Player{}).
					Where("user_discord_id =?", m.Mentions[0].ID).
					Where("guild_id =?", m.GuildID).
					Count(&count)
				if result.Error != nil {
					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s Something's wrong, Please try again later!",
						m.Author.Mention(),
					))
					break
				}
				if count > 0 {
					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s Player already added!",
						m.Author.Mention(),
					))
					break
				}

				insertErr := db.Connection.Create(&models.Player{
					Name:          m.Mentions[0].Username,
					AccountID:     command[2],
					UserDiscordID: m.Mentions[0].ID,
					GuildID:       m.GuildID,
				}).Error
				if insertErr != nil {
					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s Cannot add the player, Please try again later!",
						m.Author.Mention(),
					))
					break
				}

				a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
					"%s Player added successfully!",
					m.Author.Mention(),
				))
				break

			case "mh":

				if len(command) == 2 {

					var matchID = command[1]

					if _, err := strconv.ParseInt(matchID, 10, 64); err != nil {
						a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
							"%s The match id should be a number!",
							m.Author.Mention(),
						))
						return
					}

					message, _ := a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s Looking for a match ...",
						m.Author.Mention(),
					))
					s.ChannelTyping(channel.ID)

					msg, err := api.GetMatchHistory(matchID, false, false, false, a.Client, guild)
					if err != nil {
						a.Client.ChannelMessageEdit(channel.ID, message.ID, fmt.Sprintf(
							"%s %s",
							m.Author.Mention(),
							err.Error(),
						))
						return
					}

					a.Client.ChannelMessageEdit(channel.ID, message.ID, fmt.Sprintf(
						"%s %s",
						m.Author.Mention(),
						msg,
					))
					return

				} else {

					a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
						"%s No match_id argument found\n***Try with an argument like***  `-mh [match_id]`",
						m.Author.Mention(),
					))
					return
				}

			}
		}
	})

	// Open the websocket and begin listening.
	if err = a.Client.Open(); err != nil {
		log.Fatalf("Error opening Discord session: %v ", err)
	}

	log.Println("Pepe.bot is now running!")
}

func (a *Application) ListenAndServeGSIHttpServer(host string, port int) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")

		var (
			response  = new(components.GSIResponse)
			guild     = new(models.Guild)
			authToken = response.GetAuthToken()
		)

		if err := json.NewDecoder(r.Body).Decode(response); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":   0,
				"status": "failed",
			})
			return
		}

		if response.Provider.Appid != 570 {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":   0,
				"status": "failed",
			})
			return
		}

		if response.Map.Name != "start" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":   0,
				"status": "failed",
			})
			return
		}

		if authToken == "" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":   0,
				"status": "failed",
			})
			return
		}

		result := db.Connection.Where("token =?", response.GetAuthToken()).First(&guild)
		if err := result.Error; err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":   0,
				"status": "failed",
			})
			return
		}

		discordGuild, err := a.Client.Guild(guild.DiscordID)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":   0,
				"status": "failed",
			})
			return
		}

		response.DiscordGuild = discordGuild
		response.Guild = guild

		a.GsiChannel <- response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":   200,
			"status": "success",
		})

	})

	log.Println(fmt.Sprintf("Dota 2 GSI Http server running! on %s:%d", host, port))

	// Listen and serve the gsi application
	log.Println(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
}

func (a *Application) CreateEmojiIfNotExists(GuildId string, emojiName, filename string) (emoji *discordgo.Emoji, err error) {
	img, err := imgbase64.FromLocal(filename)
	if err != nil {
		return
	}
	emoji, err = a.Client.GuildEmojiCreate(GuildId, emojiName, img, []string{})
	return
}
