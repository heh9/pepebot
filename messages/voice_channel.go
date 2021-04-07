package messages

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/mrjoshlab/pepe.bot/db"
	"github.com/mrjoshlab/pepe.bot/models"
)

func SetMainVoiceChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {

	dbGuild, err := GetDBGuild(i)
	if err != nil {
		return
	}

	if i.Member.User.ID != dbGuild.UserID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Only the owner of the server can change main_voice_channel!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	channelId := i.Data.Options[0].StringValue()

	mainVoiceChannel, err := s.Channel(channelId)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Could not find any voice channels with id: `%s`!",
					i.Member.Mention(),
					channelId,
				),
			},
		})
		return
	}

	updateQuery := db.Connection.Model(&models.Guild{}).Where("discord_id =?", dbGuild.DiscordID).Updates(map[string]interface{}{
		"main_voice_channel_id": mainVoiceChannel.ID,
	})

	if updateQuery.Error != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Could not change MainVoiceChannel, Please try again later!",
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
				"%s MainVoiceChannel changed to %s successfully!",
				i.Member.Mention(),
				mainVoiceChannel.Name,
			),
		},
	})
	return
}
