package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/LucasSnatiago/gomusic-bot/commands"
	"github.com/LucasSnatiago/gomusic-bot/config"
	"github.com/LucasSnatiago/gomusic-bot/music"
	"github.com/bwmarrin/discordgo"
)

var BotConfig *config.Config

func init() {
	BotConfig = config.ReadConfig()
	if BotConfig == nil {
		log.Fatalf("Could not read config file")
	}

	music.SetupYtdlp()
}

func main() {
	dg, err := discordgo.New("Bot " + BotConfig.Token)
	if err != nil {
		panic(err)
	}

	dg.AddHandler(HandleCommands)
	dg.AddHandler(ready)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	err = dg.Open()
	if err != nil {
		panic(err)
	}
	// Cleanly close down the Discord session.
	defer dg.Close()
	defer fmt.Println("\nPowering off bot.")

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("GOmusic is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	s.UpdateGameStatus(0, fmt.Sprintf("%shelp", BotConfig.BotPrefix))
}

func HandleCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages sent by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore all messages without the prefix
	if !strings.HasPrefix(m.Content, BotConfig.BotPrefix) {
		return
	}

	// Split message in the command and its arguments
	message := strings.Replace(m.Content, BotConfig.BotPrefix, "", 1)
	cmd_and_args := strings.Split(message, " ")

	// All available commands
	switch strings.ToLower(cmd_and_args[0]) {
	case "ping":
		commands.Ping(s, m)
	case "play":
		url := strings.Join(cmd_and_args[1:], " ")
		buf := music.YtGetBytes(url)
		dvc, _ := music.ConnectVoiceChannel(s, m)
		commands.PlaySong(buf, dvc)
		music.DisconnectVoiceChat(dvc)
	case "echo":
		dvc, _ := music.ConnectVoiceChannel(s, m)
		commands.Echo(dvc)
		music.DisconnectVoiceChat(dvc)
	case "playtest":
		dvc, _ := music.ConnectVoiceChannel(s, m)
		commands.PlayTestSong(dvc)
		music.DisconnectVoiceChat(dvc)
	}
}
