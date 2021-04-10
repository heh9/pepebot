package messages

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/mrjosh/pepebot/components"
	"github.com/mrjosh/pepebot/disc"
	cmap "github.com/orcaman/concurrent-map"
)

func GetGuildMatch(a cmap.ConcurrentMap, guildID string) (*components.GuildMatch, bool) {
	guildMatch, ok := a.Get(guildID)
	if ok {
		return guildMatch.(*components.GuildMatch), ok
	}
	return nil, false
}

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

func ConnectVoiceChannel(a cmap.ConcurrentMap) CommandHandler {

	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		guild, err := s.Guild(i.GuildID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: fmt.Sprintf(
						"%s Could not get guild from database! Contact administrator for this",
						i.Member.Mention(),
					),
				},
			})
			return
		}

		var (
			voiceState         *discordgo.VoiceState
			userInVoiceChannel bool
		)

		log.Println(guild.VoiceStates)

		for _, vs := range guild.VoiceStates {

			log.Println(vs.UserID, i.Member.User.ID)

			if vs.UserID == i.Member.User.ID {
				voiceState = vs
				userInVoiceChannel = true
				log.Println("User found :", vs)
				break
			}
		}

		if userInVoiceChannel {

			vChannel := &disc.Channel{
				ID:     voiceState.ChannelID,
				Client: s,
			}

			voiceConnection, err := vChannel.Join()
			if err != nil {

				voiceChannel, err := s.Channel(vChannel.ID)
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionApplicationCommandResponseData{
							Content: fmt.Sprintf(
								"%s Cannot join channel %s",
								i.Member.Mention(),
								voiceChannel.Mention(),
							),
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Content: fmt.Sprintf(
							"%s Cannot join channel %s",
							i.Member.Mention(),
							voiceChannel.Mention(),
						),
					},
				})
				return
			}

			dbGuild, err := GetDBGuild(i)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Content: fmt.Sprintf(
							"%s Could not get guild from database! Contact administrator for this",
							i.Member.Mention(),
						),
					},
				})
				return
			}

			a.Set(i.GuildID, &components.GuildMatch{
				VoiceConnection: voiceConnection,
				Guild:           dbGuild,
				DiscordGuild:    guild,
				Runes:           components.NewRunes(),
			})

		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: fmt.Sprintf(
						"%s You should connect to a voice channel first then use /join",
						i.Member.Mention(),
					),
				},
			})
			return
		}

	}
}

func DisconnectVoiceChannel(a cmap.ConcurrentMap) CommandHandler {

	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		dbGuild, err := GetDBGuild(i)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: fmt.Sprintf(
						"%s Could not get guild from database! Contact administrator for this",
						i.Member.Mention(),
					),
				},
			})
			return
		}

		gm, ok := GetGuildMatch(a, dbGuild.DiscordID)
		if ok {
			if gm.HasVoiceConnection() {

				if err := gm.VoiceConnection.Disconnect(); err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionApplicationCommandResponseData{
							Content: fmt.Sprintf(
								"%s Cannot leave channel!",
								i.Member.Mention(),
							),
						},
					})
					return
				}

				a.Remove(gm.DiscordGuild.ID)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Content: fmt.Sprintf(
							"%s Bot disconnected from voice channel!",
							i.Member.Mention(),
						),
					},
				})
				return
			}
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: fmt.Sprintf(
						"%s No channel to leave",
						i.Member.Mention(),
					),
				},
			})
			return
		}
	}
}
