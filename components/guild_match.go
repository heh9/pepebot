package components

import (
	"encoding/binary"
	"github.com/MrJoshLab/pepe.bot/models"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"os"
)

type GuildMatch struct {
	VoiceConnection   *discordgo.VoiceConnection
	Guild             *models.Guild
	DiscordGuild      *discordgo.Guild
	GameEnded         bool
	Runes             *Runes
}

func (g *GuildMatch) PlaySound(sound string) bool {

	if g.HasVoiceConnection() {

		buffer, err := g.loadSound(sound)

		if err != nil {
			return false
		}

		// Start speaking.
		_ = g.VoiceConnection.Speaking(true)

		// Send the buffer data.
		for _, buff := range buffer {
			g.VoiceConnection.OpusSend <- buff
		}

		_ = g.VoiceConnection.Speaking(false)

		return true
	}

	return false
}

// loadSound attempts to load an encoded sound file from disk.
func (g *GuildMatch) loadSound(sound string) ([][]byte, error) {

	buffer := make([][]byte, 0)

	file, err := os.Open("./sounds/" + sound + ".dca")
	if err != nil {
		log.Println("Error opening dca file :", err)
		return nil, err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return nil, err
			}
			return buffer, err
		}

		if err != nil {
			log.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			log.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}

	return buffer, nil
}

func (g *GuildMatch) HasVoiceConnection() bool {
	return g.VoiceConnection != nil
}
