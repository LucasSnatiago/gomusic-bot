package music

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ConnectVoiceChannel(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("s.State.Channel:", err)
		return err
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		log.Println("s.State.Guild:", err)
		return err
	}

	voiceChannelID := ""

	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			voiceChannelID = vs.ChannelID
		}
	}

	if voiceChannelID == "" {
		s.ChannelMessageSend(m.ChannelID, "User must specify the Channel ID or join one!")
	}

	// Connect to voice channel.
	// NOTE: Setting mute to false, deaf to true.
	dgv, err := s.ChannelVoiceJoin(m.GuildID, voiceChannelID, false, true)
	if err != nil {
		log.Println(err)
		return err
	}
	// Disconnect after finishing song
	defer dgv.Disconnect()

	dgv.Speaking(true)

	// Play song
	time.Sleep(1 * time.Second)

	dgv.Speaking(false)

	// Wait before leaving
	time.Sleep(250 * time.Millisecond)

	return nil
}
