package messages

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/mrjosh/pepebot/db"
	"github.com/mrjosh/pepebot/models"
)

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

	dbGuild, err := GetDBGuild(i)
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

//mainTextChannel, err := a.Client.Channel(dbGuild.MainTextChannelID)
//if err != nil {
//return
//}

//a.Client.ChannelMessageSend(channel.ID, fmt.Sprintf(
//"%s The main text channel for matches is %s\n"+
//"If you want to change it, you can use `-main_text_channel {channel_id}`",
//m.Author.Mention(),
//mainTextChannel.Mention(),
//))

//}
