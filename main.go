package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token   string
	Channel string
	Servers []string
	Rate    time.Duration
}

var defaultConfig = Config{
	Rate: 5 * time.Minute,
}

func main() {
	var c Config = defaultConfig
	envconfig.Process("arkbot", &c)

	if c.Token == "" {
		log.Fatal("No token provided.")
	}

	if c.Channel == "" {
		log.Fatal("No channel provided.")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + c.Token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %s", err)
		return
	}

	log.Printf("Starting tracker in channel %s", c.Channel)
	for _, s := range c.Servers {
		log.Printf("Monitoring server %s", s)
	}

	closer := make(chan bool)
	go runTracker(dg, c.Channel, c.Servers, closer)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %s", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Print("arkbot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	closer <- true
	dg.Close()
}
