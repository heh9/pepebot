package messages

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mrjosh/pepebot/db"
	"github.com/mrjosh/pepebot/models"
)

func GetDBGuild(i *discordgo.InteractionCreate) (*models.Guild, error) {
	var (
		dbGuild = new(models.Guild)
		result  = db.Connection.Where("discord_id =?", i.GuildID).Where("user_id =?", i.Member.User.ID).First(&dbGuild)
	)
	if err := result.Error; err != nil {
		return nil, err
	}
	return dbGuild, nil
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

	dbGuild, err := GetDBGuild(i)
	if err != nil {
		return
	}

	helpText := "\n"

	if dbGuild.MainTextChannelID == "" {
		helpText += "Set MainTextChannel with `/main_text_channel channel_id: {channel_id}`"
	} else {
		mainTextChannel, err := s.Channel(dbGuild.MainTextChannelID)
		if err == nil {
			helpText +=
				"MainTextChannel is " + mainTextChannel.Mention() + " \n" +
					"If you want to change this you can use: `/main_text_channel channel_id: {channel_id}`\n\n"
		} else {
			helpText += "Set MainTextChannel with `/main_text_channel channel_id: {channel_id}`"
		}
	}

	if dbGuild.MainVoiceChannelID == "" {
		helpText += "Set MainVoiceChannel with `/main_voice_channel channel_id: {voice_channel_id}`\n"
	} else {
		mainVoiceChannel, err := s.Channel(dbGuild.MainVoiceChannelID)
		if err == nil {
			helpText +=
				"MainVoiceChannel is " + mainVoiceChannel.Mention() + " \n" +
					"If you want to change this you can use: `/main_voice_channel channel_id: {channel_id}`\n\n"
		} else {
			helpText += "Set MainVoiceChannel with `/main_voice_channel channel_id: {channel_id}`\n"
		}
	}

	helpText +=

		"Connect a dota2 player `/player_add user: [user] steam_account_id: [account_id]`\n" +
			"Disconnect a dota2 player  `/player_remove user: [user]`\n" +
			"Get a summary of a dota2 match `/match_history match_id: [match_id]`\n" +
			"Join voice channel that youre in `/join`\n" +
			"Disconnect from voice channel that bot connected to `/leave`\n\n " +

			"** If you having some issues with disconnecting the bot manually\n" +
			"  you should use `/leave` command to disconnect it! ** \n\n" +

			"This bot is a runes reminder bot for dota 2 games that works with" +
			" Dota 2 GSI API.\n" +
			"You can get install instruction with command: `/instructions` \n" +
			"Isn't that cool ? "

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: helpText,
		},
	})
}
