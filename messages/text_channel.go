package messages

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mrjosh/pepebot/db"
	"github.com/mrjosh/pepebot/models"
)

func SetMainTextChannelWithMessageCreate(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s This command is only allowed in a discord server!",
			m.Author.Mention(),
		))
		return
	}

	dbGuild, err := GetDBGuild(guild.ID)
	if err != nil {
		return
	}

	if len(args) == 2 {

		channelId := strings.TrimSpace(args[1])
		if _, err := strconv.Atoi(channelId); err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				"%s channel_id should be a valid text channel id number!",
				m.Author.Mention(),
			))
			return
		}

		if m.Author.ID != dbGuild.UserID {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				"%s Only the owner of the server can change main_voice_channel!",
				m.Author.Mention(),
			))
			return
		}

		mainTextChannel, err := s.Channel(channelId)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				"%s Could not find any text channel with id: `%s`!",
				m.Author.Mention(),
				channelId,
			))
			return
		}

		if mainTextChannel.Type != discordgo.ChannelTypeGuildText {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				"%s Channel type should be a text channel!",
				m.Author.Mention(),
			))
			return
		}

		updateQuery := db.Connection.Model(&models.Guild{}).Where("discord_id =?", dbGuild.DiscordID).Updates(map[string]interface{}{
			"main_text_channel_id": mainTextChannel.ID,
		})

		if updateQuery.Error != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				"%s Could not change MainTextChannel, Please try again later!",
				m.Author.Mention(),
			))
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s MainTextChannel changed to %s successfully!",
			m.Author.Mention(),
			mainTextChannel.Mention(),
		))
		return
	}

	if dbGuild.MainTextChannelID == "" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s Use `-main_text_channel {text_channel_id}` to set!",
			m.Author.Mention(),
		))
		return
	}

	mainTextChannel, err := s.Channel(dbGuild.MainTextChannelID)
	if err != nil {
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
		"%s The main text channel for matches is %s\n"+
			"If you want to change it, you can use `-main_text_channel {channel_id}`",
		m.Author.Mention(),
		mainTextChannel.Mention(),
	))
	return
}

func SetMainTextChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.GuildID == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "Use this command in a discord server!",
			},
		})
		return
	}

	guild, err := s.Guild(i.GuildID)
	if err != nil {
		return
	}

	if i.Member.User.ID != guild.OwnerID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Only the owner of the server can change main_text_channel!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	channelId := i.Data.Options[0].StringValue()
	mainTextChannel, err := s.Channel(channelId)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Could not find any channels with id: `%s`!",
					i.Member.Mention(),
					channelId,
				),
			},
		})
		return
	}

	dbGuild, err := GetDBGuild(i.GuildID)
	if err != nil {
		return
	}

	updateQuery := db.Connection.Model(&models.Guild{}).Where("discord_id =?", dbGuild.DiscordID).Updates(map[string]interface{}{
		"main_text_channel_id": mainTextChannel.ID,
	})

	if updateQuery.Error != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Could not change MainTextChannel, Please try again later!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf(
				"%s MainTextChannel changed to %s successfully!",
				i.Member.Mention(),
				mainTextChannel.Mention(),
			),
		},
	})
	return
}
