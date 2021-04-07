package messages

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/mrjoshlab/pepe.bot/api"
)

func ShowMatchHistory(s *discordgo.Session, i *discordgo.InteractionCreate) {

	guild, err := s.Guild(i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s %s",
					i.Member.Mention(),
					"Could not find guild!!",
				),
			},
		})
		return
	}

	matchID := i.Data.Options[0].IntValue()
	matchIDString := strconv.Itoa(int(matchID))

	msg, err := api.GetMatchHistory(matchIDString, false, false, false, s, guild)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s %s",
					i.Member.Mention(),
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
				i.Member.Mention(),
				msg,
			),
		},
	})
	return
}
