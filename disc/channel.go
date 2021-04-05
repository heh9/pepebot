package disc

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Channel struct {
	ID     string
	Client *discordgo.Session
}

func (c *Channel) Join() (*discordgo.VoiceConnection, error) {

	channel, err := c.Client.Channel(c.ID)
	if err != nil {
		return nil, err
	}

	// Find the guild for that channel.
	guild, err := c.Client.State.Guild(channel.GuildID)
	if err != nil {
		return nil, err
	}

	voiceConnection, err := c.Client.ChannelVoiceJoin(guild.ID, channel.ID, false, true)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return voiceConnection, nil
}
