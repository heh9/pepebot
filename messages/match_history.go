package messages

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mrjosh/pepebot/api"
)

func ShowMatchHistoryWithMessageCreate(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {

	if len(args) == 2 {

		matchID := strings.TrimSpace(args[1])
		if _, err := strconv.Atoi(matchID); err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				"%s match_id should be a valid number!",
				m.Author.Mention(),
			))
			return
		}

		message, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s Looking for a match ...",
			m.Author.Mention(),
		))
		s.ChannelTyping(m.ChannelID)

		msg, err := api.GetMatchHistory(matchID, false, false, false, s)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				"%s %s",
				m.Author.Mention(),
				err.Error(),
			))
			return
		}

		s.ChannelMessageEdit(m.ChannelID, message.ID, fmt.Sprintf(
			"%s %s",
			m.Author.Mention(),
			msg,
		))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
		"%s Use `-match_history {match_id}`",
		m.Author.Mention(),
	))
	return
}

func ShowMatchHistory(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var user *discordgo.User
	if i.GuildID == "" {
		user = i.User
	} else {
		user = i.Member.User
	}

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
