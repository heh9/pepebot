package main

import (
	"os"
	"io"
	"log"
	"sync"
	"strings"
	"strconv"
	"net/http"
	"encoding/json"
	"encoding/binary"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
	"github.com/polds/imgbase64"
)

var (
	dota2DataChannel = make(chan *Response)
	voiceChannel *discordgo.VoiceConnection
	token = os.Getenv("DISCORD_API_TOKEN")
)

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func main()  {

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")

		response := &Response{}
		json.NewDecoder(r.Body).Decode(response)

		if response.CheckAuthToken(os.Getenv("DOTA2_GSI_AUTH_TOKEN")) {
			dota2DataChannel <- response
			json.NewEncoder(w).Encode(map[string] string {
				"status": "Ok",
			})
		} else {
			json.NewEncoder(w).Encode(map[string] string {
				"status": "Failed",
			})
		}
	})

	go func() {

		wg := sync.WaitGroup{}
		wg.Add(1)

		for data := range dota2DataChannel {

			runes := &Runes{ Data: strconv.Itoa(data.Map.ClockTime), Saied: false }

			if runes.AreRunesUp() {
				playSound("runes")
				wg.Done()
			}
		}

		wg.Wait()

	}()

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println(err)
	}

	discord.AddHandler(func (s *discordgo.Session, guild *discordgo.GuildCreate) {

		var exists = false

		for _, guildEmoji := range guild.Emojis {
			if guildEmoji.Name == "peepoblush" {
				exists = true
				break
			}
		}

		if !exists {
			img, err := imgbase64.FromLocal("./emojies/peepoblush.png")
			if err != nil {
				log.Println(err)
				return
			}

			_, err = discord.GuildEmojiCreate(guild.ID, "peepoblush", img, []string {})
			if err != nil {
				log.Println(err)
				return
			}
		}

	})

	discord.AddHandler(func (s *discordgo.Session, event *discordgo.Ready) {
		s.UpdateStatus(0, "Dota 2 [-help]")
		log.Println("Logged in as !" + event.User.ID)
	})

	discord.AddHandler(func (s *discordgo.Session, m *discordgo.MessageCreate) {

		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, "-") {

			channel, err := s.State.Channel(m.ChannelID)
			if err != nil {
				// Could not find channel.
				return
			}

			// Find the guild for that channel.
			g, err := s.State.Guild(channel.GuildID)
			if err != nil {
				// Could not find guild.
				return
			}

			prefixCommand := strings.Split(m.Content, "-")

			command := strings.Split(prefixCommand[1], " ")

			var pepeEmoji *discordgo.Emoji

			for _, guildEmoji := range g.Emojis {
				if guildEmoji.Name == "peepoblush" {
					pepeEmoji = guildEmoji
					break
				}
			}

			switch command[0] {
			case "help":
				discord.ChannelMessageSend(channel.ID, "This bot is a runes reminder bot for dota 2 games that works with" +
					" Dota 2 GSI API. isn't that cool ? :)")
			break
			case "gm":

				var matchID int

				if len(command) == 2 {

					if isNumeric(command[1]) {
						matchID, err = strconv.Atoi(command[1])
					} else {
						discord.ChannelMessageSend(channel.ID, "***Match ID*** should be a number!")
						return
					}

				} else {

					discord.ChannelMessageSend(channel.ID, "No match found <:" + pepeEmoji.Name + ":" + pepeEmoji.ID + ">\n" +
						"***Try with an argument like***  `-gm [match_id]`")
					return
				}

				lookingMsg := "**Looking for a game with match id** `" + strconv.Itoa(matchID) + "`\n"
				message, _ := discord.ChannelMessageSend(channel.ID, lookingMsg + "Please wait ...")

				s.ChannelTyping(channel.ID)

				players, _ := GetPlayerDatasByMatchID(matchID)

				if players["dire"] == nil && players["radiant"] == nil {

					discord.ChannelMessageEdit(channel.ID, message.ID, "No match found <:" + pepeEmoji.Name + ":" + pepeEmoji.ID + ">\n")
					break
				}

				msg := "**Found a game with match id** `" + strconv.Itoa(matchID) + "`\n"

				msg += "\n***Suggest to BAN :octagonal_sign:*** \n"

				for _, direPlayer := range players["dire"] {
					for i, hero := range direPlayer.MostHeroPlayed {
						if i == 3 { break }
						msg += " **`" + hero.Hero.LocalizedName + "`** \n"
					}
				}

				msg += "\n***Dire Players :video_game: *** \n"
				for _, direPlayer := range players["dire"] {
					msg += direPlayer.Name + "`" + direPlayer.RankName + "`\n"
				}

				msg += "\n***Radiant Players :video_game: *** \n"
				for _, radiantPlayer := range players["radiant"] {
					msg += radiantPlayer.Name + "`" + radiantPlayer.RankName + "`\n"
				}

				discord.ChannelMessageEdit(channel.ID, message.ID, msg)
			break
			case "leave":
				if voiceChannel != nil {
					voiceChannel.Disconnect()
					voiceChannel = nil
				}
				break
			case "join":

				if voiceChannel != nil {
					voiceChannel.Disconnect()
					voiceChannel = nil
				}

				for _, vs := range g.VoiceStates {

					if vs.UserID == m.Author.ID {

						voiceChannel, err = s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
						if err != nil {
							log.Println(err)
						}

						ch, err := s.Channel(voiceChannel.ChannelID)
						if err != nil {
							log.Println(err)
						}

						discord.ChannelMessageSend(channel.ID, ":white_check_mark: Bot successfully connected to " + ch.Name)
						return
					}
				}

				if voiceChannel == nil {

					discord.ChannelMessageSend(channel.ID, "You must be connected to a voice channel! \n" +
						"Connect to a voice channel and try -join!")
					return
				}
				break
			case "play runes":
				if !playSound("runes") {

					discord.ChannelMessageSend(channel.ID, "Bot is not connected to a voice channel. \n" +
						"Try -join to join a voice channel youre in.")
				}
				break
			}
		}
	})

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		log.Println("Error opening Discord session: ", err)
	}

	log.Println("Pepe.bot is now running.  Press CTRL-C to exit.")

	// Cleanly close down the Discord session.
	defer discord.Close()

	http.ListenAndServe(":9001", nil)
}

// loadSound attempts to load an encoded sound file from disk.
func loadSound(sound string) ([][]byte, error) {

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

func playSound(sound string) bool {

	if voiceChannel != nil {

		buffer, err := loadSound(sound)

		if err != nil {
			return false
		}

		// Start speaking.
		voiceChannel.Speaking(true)

		// Send the buffer data.
		for _, buff := range buffer {
			voiceChannel.OpusSend <- buff
		}

		return true
	}

	return false
}