package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/LucasSnatiago/gomusic-bot/external"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "pong!")
}

func PlaySong(song []byte, v *discordgo.VoiceConnection) error {
	// Cria um canal para enviar dados PCM
	send := make(chan []int16, 2)
	defer close(send)

	// Inicia o envio de áudio em uma goroutine
	go dgvoice.SendPCM(v, send)

	// Converte os bytes do áudio para PCM 16-bit estéreo
	reader := bytes.NewReader(song)
	buffer := make([]int16, 1920) // Tamanho do buffer (ajuste conforme necessário)

	for {
		// Lê os dados PCM do áudio
		err := binary.Read(reader, binary.LittleEndian, &buffer)
		if err != nil {
			if err == io.EOF {
				break // Fim do arquivo
			}
			return fmt.Errorf("erro ao ler dados PCM: %w", err)
		}

		// Envia os dados PCM para o canal
		send <- buffer

		// Atraso entre frames (20ms é o padrão para áudio em tempo real)
		time.Sleep(20 * time.Millisecond)
	}

	return nil
}

// Takes inbound audio and sends it right back out.
func Echo(v *discordgo.VoiceConnection) {
	recv := make(chan *discordgo.Packet, 2)
	go dgvoice.ReceivePCM(v, recv)

	send := make(chan []int16, 2)
	go dgvoice.SendPCM(v, send)

	v.Speaking(true)
	defer v.Speaking(false)

	for {
		p, ok := <-recv
		if !ok {
			return
		}

		send <- p.PCM
	}
}

func PlayTestSong(v *discordgo.VoiceConnection) error {
	// Cria um reader a partir dos dados embutidos
	reader := bytes.NewReader(external.Airhorn_dca)
	var opuslen int16

	for {
		// Lê o tamanho do frame Opus (2 bytes little-endian)
		err := binary.Read(reader, binary.LittleEndian, &opuslen)
		if err != nil {
			if err == io.EOF {
				break // Fim do arquivo
			}
			fmt.Println("Erro ao ler tamanho do frame:", err)
			return err
		}

		// Lê os dados do frame Opus
		frame := make([]byte, opuslen)
		err = binary.Read(reader, binary.LittleEndian, &frame)
		if err != nil {
			fmt.Println("Erro ao ler dados do frame:", err)
			return err
		}

		// Envia o frame para o canal OpusSend
		v.OpusSend <- frame

		// Atraso entre frames (exemplo: 20ms para áudio em tempo real)
		time.Sleep(20 * time.Millisecond)
	}

	return nil
}
