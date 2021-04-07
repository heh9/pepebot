package messages

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/mrjoshlab/pepe.bot/db"
	"github.com/mrjoshlab/pepe.bot/models"
)

func AddPlayer(s *discordgo.Session, i *discordgo.InteractionCreate) {

	dbGuild, err := GetDBGuild(i)
	if err != nil {
		return
	}

	if dbGuild.UserID != i.Member.User.ID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Only the owner of the guild can add/remove/update a player!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	user := i.Data.Options[0].UserValue(s)
	steamAccountId := i.Data.Options[1].StringValue()

	if user.Bot {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s You sould mention only one member and it can not be a bot!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	var count int64
	result := db.Connection.Model(&models.Player{}).
		Where("user_discord_id =?", user.ID).
		Where("guild_id =?", i.GuildID).
		Count(&count)

	if result.Error != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Something's wrong, Please try again later!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	if count > 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Player already added!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	insertErr := db.Connection.Create(&models.Player{
		Name:          user.Username,
		AccountID:     steamAccountId,
		UserDiscordID: user.ID,
		GuildID:       i.GuildID,
	}).Error
	if insertErr != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Cannot add the player, Please try again later!",
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
				"%s Player added successfully!",
				i.Member.Mention(),
			),
		},
	})
	return
}

func RemovePlayer(s *discordgo.Session, i *discordgo.InteractionCreate) {

	dbGuild, err := GetDBGuild(i)
	if err != nil {
		return
	}

	if dbGuild.UserID != i.Member.User.ID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Only the owner of the guild can add/remove/update a player!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	user := i.Data.Options[0].UserValue(s)

	if user.Bot {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s You sould mention only one member and it can not be a bot!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	var count int64
	result := db.Connection.Model(&models.Player{}).
		Where("user_discord_id =?", user.ID).
		Where("guild_id =?", i.GuildID).
		Count(&count)

	if result.Error != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Something's wrong, Please try again later!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	if count == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Player does not exists!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	deleteErr := db.Connection.Where("user_discord_id =?", user.ID).
		Where("guild_id =?", i.GuildID).
		Delete(&models.Player{}).Error

	if deleteErr != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Cannot remove the player, Please try again later!",
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
				"%s Player removed successfully!",
				i.Member.Mention(),
			),
		},
	})
	return
}
