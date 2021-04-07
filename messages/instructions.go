package messages

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SendInstructions(s *discordgo.Session, i *discordgo.InteractionCreate) {

	guild, err := s.Guild(i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Could not send the instructions at the time. Please try again later!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	if i.Member.User.ID != guild.OwnerID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Only the owner of the server can get instructions!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	dbGuild, err := GetDBGuild(i)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Could not send the instructions at the time. Please try again later!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	if dbGuild.MainVoiceChannelID == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s First of all, you need to set your main_voice_channel. \n"+
						"Use `/main_voice_channel {channel_id}` and then ask me for instructions!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	instructions :=
		"\n" +
			"Dota2 GSI Api installation" +
			"\n\n" +

			"Find your dota2 local directory, Go to `Steam` and right click on `Dota2` then: \n" +
			"> `Properties > Local Files > Browse Local Files` \n\n" +

			"Then go to `game/dota2/cfg` and create a directory called `gamestate_integration`\n" +
			"go to `gamestate_integration` and create a file called `gamestate_integration.cfg`\n" +
			"and paste the below content into it! \n\n```" +

			`"dota2-gsi Configuration"` +
			"\n{\n" +
			`    "uri"               "https://pepebot.mrjosh.net"` + "\n" +
			`    "timeout"           "5.0"` + "\n" +
			`    "buffer"            "0.1"` + "\n" +
			`    "throttle"          "0.1"` + "\n" +
			`    "heartbeat"         "30.0"` + "\n" +
			`    "data"` + "\n" +
			"    {" + "\n" +
			`        "provider"      "1"` + "\n" +
			`        "map"           "1"` + "\n" +
			`        "player"        "1"` + "\n" +
			"    }" + "\n" +
			`    "auth"` + "\n" +
			"    {" + "\n" +
			`         "token"         "` + dbGuild.Token + `"` + "\n" +
			"    }" + "\n" +
			"}```" +

			"\n Restart your game and You're ready to go find some matches! \n" +
			"I will connect to main_voice_channel which is `" + dbGuild.MainVoiceChannelID + "` in your server \n" +
			"and remind you the runes every 5 minutes :sunglasses: ! \n\n" +

			"Give us some feedback or write your issues here > https://github.com/mrjoshlab/pepe.bot/issues :heart:"

	channel, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Could not send the instructions at the time. Please try again later!",
					i.Member.Mention(),
				),
			},
		})
		return
	}

	if _, err = s.ChannelMessageSend(channel.ID, instructions); err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					"%s Could not send the instructions at the time. Please try again later!",
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
				"%s The instructions sent to your private chat successfully!",
				i.Member.Mention(),
			),
		},
	})
}
