package main

import (
	"log"
	"errors"
	"strings"
	"strconv"
	"pepe.bot/disc"
	"github.com/polds/imgbase64"
	"github.com/MrJoshLab/arrays"
	"github.com/bwmarrin/discordgo"
)

type Application struct {
	Client               *discordgo.Session
	Guild                *discordgo.Guild
	Emojis               []*discordgo.Emoji
	GsiChannel           chan *GSIResponse
	GameEndChannel       chan *EndStruct
	GameEnded            bool
	VoiceChannel         *discordgo.VoiceConnection
	DiscordAuthToken     string
	MainTextChannelId    string
	MainVoiceChannelId   string
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

					endStruct := &EndStruct{
						MatchId: data.Map.Matchid,
					}

					if data.Player.TeamName != data.Map.WinTeam {
						endStruct.Won = false
						a.PlaySound("loss")
					} else {
						endStruct.Won = true
						a.PlaySound("win")
					}

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

		pepeEmoji := a.GetEmoji("peepoblush")

		var wonText = "lost"
		var StatusText = "Try a bit harder next time " + pepeEmoji.MessageFormat()

		if game.Won {
			wonText = "win"
			StatusText = "Weeeeee Areeeee the championssssssss my friendsss " + pepeEmoji.MessageFormat()
		}

		a.Client.ChannelMessageSend(a.MainTextChannelId,
			"```css\n" +
				"Game ended as " + wonText + " with match id [" + game.MatchId + "]" +
				"```" + StatusText)

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
				a.Client.ChannelMessageSend(channel.ID, "This bot is a runes reminder bot for dota 2 games that works with" +
					" Dota 2 GSI API. \n" +
					"Isn't that cool ? " + a.GetEmoji("peepoblush").MessageFormat())
				break
			case "leave":
				if a.VoiceChannel != nil {
					a.VoiceChannel.Disconnect()
					a.VoiceChannel = nil
				}
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

						a.Client.ChannelMessageSend(channel.ID, ":white_check_mark: Bot successfully connected to " + ch.Name)
						return
					}
				}

				if a.VoiceChannel == nil {

					a.Client.ChannelMessageSend(
						channel.ID,
						"You must be connected to a voice channel! \n" +
							"Connect to a voice channel and try -join!")
					return
				}
				break
			case "play":

				if len(command) == 2 {

					if coll := collection.New(a.Runes.Sounds); !coll.Has(command[1]) {
						a.Client.ChannelMessageSend(
							channel.ID,
							"Could not found sound " + command[1] + " " + a.GetEmoji("peepoblush").MessageFormat())
						return
					}

				} else {

					a.Client.ChannelMessageSend(
						channel.ID,
						"No sound argument found " + a.GetEmoji("peepoblush").MessageFormat() + "\n" +
							"***Try with an argument like***  `-play [sound_name]`")
					return
				}

				if !a.PlaySound(a.Runes.GetRandomVoiceFileName()) {

					a.Client.ChannelMessageSend(
						channel.ID,
						"Bot is not connected to a voice channel. \n" +
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

	defer a.Client.Close()
}

func (a *Application) GetEmoji(name string) *discordgo.Emoji {

	for _, emoji := range a.Emojis {
		if emoji.Name == name {
			return emoji
		}
	}

	return nil
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