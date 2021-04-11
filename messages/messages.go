package messages

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/mrjosh/pepebot/db"
	"github.com/mrjosh/pepebot/models"
)

func GetDBGuild(guildID string) (*models.Guild, error) {
	var (
		dbGuild = new(models.Guild)
		result  = db.Connection.Where("discord_id =?", guildID).First(&dbGuild)
	)
	if err := result.Error; err != nil {
		return nil, err
	}
	return dbGuild, nil
}

func SendHelpTextWithMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s This command is only allowed in a discord server!",
			m.Author.Mention(),
		))
		return
	}

	helpText, err := getHelpText(s, guild.ID)
	if err != nil {
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
		"%s %s",
		m.Author.Mention(),
		helpText,
	))
}

func SendHelpText(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.GuildID == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "Use this command in a discord server!",
			},
		})
		return
	}

	helpText, err := getHelpText(s, i.GuildID)
	if err != nil {
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: helpText,
		},
	})
}

func getHelpText(s *discordgo.Session, guildID string) (string, error) {

	dbGuild, err := GetDBGuild(guildID)
	if err != nil {
		return "", err
	}

	var (
		helpText = "\n"
		mtc      = dbGuild.MainTextChannelID
		mvc      = dbGuild.MainVoiceChannelID
	)

	if mtc == "" {
		helpText += "Set MainTextChannel with `-main_text_channel {channel_id}`\n"
	} else {
		mainTextChannel, err := s.Channel(mtc)
		if err == nil {
			helpText +=
				"MainTextChannel is " + mainTextChannel.Mention() + " \n" +
					"If you want to change this you can use: `-main_text_channel {channel_id}`\n\n"
		} else {
			helpText += "Set MainTextChannel with `-main_text_channel {channel_id}`\n"
		}
	}

	if mvc == "" {
		helpText += "Set MainVoiceChannel with `-main_voice_channel {voice_channel_id}`\n"
	} else {
		mainVoiceChannel, err := s.Channel(mvc)
		if err == nil {
			helpText +=
				"MainVoiceChannel is " + mainVoiceChannel.Mention() + " \n" +
					"If you want to change this you can use: `-main_voice_channel {channel_id}`\n\n"
		} else {
			helpText += "Set MainVoiceChannel with `-main_voice_channel {channel_id}`\n"
		}
	}

	helpText +=

		"Get a summary of a dota2 match `-match_history {match_id}`\n\n" +

			"This bot is a runes reminder bot for dota 2 games that works with" +
			" Dota 2 GSI API.\n" +
			"You can get install instruction with command: `-instructions` \n" +
			"Isn't that cool ? "

	return helpText, nil
}
