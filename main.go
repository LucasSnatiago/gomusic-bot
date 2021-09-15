package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/LucasSnatiago/gomusic-bot/config"
	"github.com/bwmarrin/discordgo"
)

func main() {
	botConfig := config.ReadConfig()
	if botConfig == nil {
		return
	}

	discord, err := discordgo.New("Bot " + botConfig.Token)
	if err != nil {
		panic(err)
	}

	err = discord.Open()
	if err != nil {
		panic(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("GOmusic is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
	fmt.Println("\nPowering off bot.")
}
