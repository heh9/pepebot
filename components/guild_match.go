package components

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/mrjoshlab/pepe.bot/config"
	"github.com/mrjoshlab/pepe.bot/models"
)

type GuildMatch struct {
	VoiceConnection *discordgo.VoiceConnection
	Guild           *models.Guild
	DiscordGuild    *discordgo.Guild
	GameEnded       bool
	Runes           *Runes
}

func (g *GuildMatch) PlaySound(sound string) error {

	if g.HasVoiceConnection() {

		buffer, err := g.loadSound(sound)
		if err != nil {
			return fmt.Errorf("Error loading sound file : %v", err)
		}

		// Start speaking.
		g.VoiceConnection.Speaking(true)

		// Send the buffer data.
		for _, buff := range buffer {
			g.VoiceConnection.OpusSend <- buff
		}

		g.VoiceConnection.Speaking(false)
	}

	return nil
}

// loadSound attempts to load an encoded sound file from disk.
func (g *GuildMatch) loadSound(sound string) ([][]byte, error) {

	buffer := make([][]byte, 0)

	file, err := os.Open(fmt.Sprintf("%s/%s.dca", config.Map.Sounds.Path, sound))
	if err != nil {
		return nil, fmt.Errorf("Error opening dca file: %v", err)
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if err := file.Close(); err != nil {
				return nil, err
			}
			return buffer, nil
		}

		if err != nil {
			return nil, fmt.Errorf("Error reading from dca file: %v", err)
		}

		// Read encoded pcm from dca file.
		buff := make([]byte, opuslen)
		if err := binary.Read(file, binary.LittleEndian, &buff); err != nil {
			return nil, fmt.Errorf("Error reading from dca file: %v", err)
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, buff)
	}

	return buffer, nil
}

func (g *GuildMatch) HasVoiceConnection() bool {
	return g.VoiceConnection != nil
}
