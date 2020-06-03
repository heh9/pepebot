package main

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

// loadSound attempts to load an encoded sound file from disk.
func (a *Application) LoadSound(sound string) (buffer [][]byte, err error) {

	buffer = make([][]byte, 0)

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

func (a *Application) PlaySound(sound string) bool {

	if a.VoiceChannel != nil {

		buffer, err := a.LoadSound(sound)

		if err != nil {
			return false
		}

		// Start speaking.
		_ = a.VoiceChannel.Speaking(true)

		// Send the buffer data.
		for _, buff := range buffer {
			a.VoiceChannel.OpusSend <- buff
		}

		_ = a.VoiceChannel.Speaking(false)

		return true
	}

	return false
}
