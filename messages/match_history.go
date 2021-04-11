package messages

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/mrjosh/pepebot/api"
)

func ShowMatchHistory(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var user *discordgo.User
	if i.GuildID == "" {
		user = i.User
	} else {
		user = i.Member.User
	}

	//guild, err := s.Guild(i.GuildID)
	//if err != nil {
	//s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	//Type: discordgo.InteractionResponseChannelMessageWithSource,
	//Data: &discordgo.InteractionApplicationCommandResponseData{
	//Content: fmt.Sprintf(
	//"%s %s",
	//user.Mention(),
	//"Could not find guild!!",
	//),
	//},
	//})
	//return
	//}

	matchID := i.Data.Options[0].IntValue()
	matchIDString := strconv.Itoa(int(matchID))

	msg, err := api.GetMatchHistory(matchIDString, false, false, false, s)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s %s",
					user.Mention(),
					err.Error(),
				),
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf(
				"%s %s",
				user.Mention(),
				msg,
			),
		},
	})
	return
}
