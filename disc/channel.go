package disc

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type Channel struct {
	ID               string
	Client           *discordgo.Session
}

func (c *Channel) Join() (bool, error, *discordgo.Channel, *discordgo.VoiceConnection) {

	channel, _ := c.Client.Channel(c.ID)

	// Find the guild for that channel.
	guild, err := c.Client.State.Guild(channel.GuildID)
	if err != nil {
		return false, nil, nil, nil
	}

	voiceChannel, err := c.Client.ChannelVoiceJoin(guild.ID, channel.ID, false, true)
	if err != nil {
		log.Println(err)
	}

	return true, err, channel, voiceChannel
}
