package music

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ConnectVoiceChannel(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.VoiceConnection, error) {
	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("s.State.Channel:", err)
		return nil, err
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		log.Println("s.State.Guild:", err)
		return nil, err
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
	dgv, err := s.ChannelVoiceJoin(m.GuildID, voiceChannelID, false, false)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dgv.Speaking(true)

	return dgv, nil
}

func DisconnectVoiceChat(dgv *discordgo.VoiceConnection) {
	dgv.Speaking(false)
	time.Sleep(250 * time.Millisecond)
	dgv.Disconnect()
}
