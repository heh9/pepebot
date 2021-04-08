package messages

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SendAboutText(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf(
				"%s This bot is a runes reminder bot for dota 2 games that works with"+
					" Dota 2 GSI API. \n"+
					"Isn't that cool ? ",
				i.Member.Mention(),
			),
		},
	})
}