package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// find the previous update to edit
func findMyMessage(s *discordgo.Session, channel, me string) (string, error) {
	before := ""

	for {
		ms, err := s.ChannelMessages(channel, 100, before, "", "")
		if err != nil {
			return "", fmt.Errorf("Error getting messages in channel %s: %w", channel, err)
		}

		if len(ms) == 0 {
			break
		}

		for _, m := range ms {
			if m.Author.ID == me {
				return m.ID, nil
			}
		}

		before = ms[0].ID
	}

	return "", nil
}

func checkAllServers(s *discordgo.Session, me, channel string, servers []string) {
	toUpdate, err := findMyMessage(s, channel, me)
	if err != nil {
		log.Print(err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Type:      "rich",
		Title:     "Currently running servers:",
		Timestamp: time.Now().Format(time.RFC3339),
		Fields:    make([]*discordgo.MessageEmbedField, 0, len(servers)),
		//Footer:    &discordgo.MessageEmbedFooter{Text: "Last updated now"},
	}

	for _, s := range servers {
		server, err := checkServer(s)
		if err != nil {
			continue
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   server.Name,
			Value:  fmt.Sprintf("%2d Players, Map: %s\n", server.Players, server.Map),
			Inline: false,
		})
	}

	if toUpdate == "" {
		_, err = s.ChannelMessageSendEmbed(channel, embed)
		if err != nil {
			log.Printf("Error sending message in channel %s: %s", channel, err)
		}
	} else {
		_, err = s.ChannelMessageEditEmbed(channel, toUpdate, embed)
		if err != nil {
			log.Printf("Error updating message in channel %s: %s", channel, err)
		}
	}
}

func runTracker(s *discordgo.Session, channel string, servers []string, rate time.Duration, closer <-chan bool) {
	t := time.Tick(rate)

	user, err := s.User("@me")
	if err != nil {
		log.Fatalf("Error getting my user: %w", err)
	}

	for {
		checkAllServers(s, user.ID, channel, servers)
		select {
		case <-t:
			continue
		case <-closer:
			return
		}
	}
}
