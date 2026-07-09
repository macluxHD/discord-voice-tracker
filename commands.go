package main

import (
	"discord-voice-tracker/internal/database"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func registerCommands(b *Bot, s *discordgo.Session) {
	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := b.commandHandlers()[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "get-user-count",
			Description: "Get event count of users in the current guild",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "event-types",
					Description: "Comma separated event types",
					Required:    true,
				},
			},
		},
	}
)

func (b *Bot) commandHandlers() map[string]func(*discordgo.Session, *discordgo.InteractionCreate) {
	return map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"get-user-count": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			eventTypes, err := database.ParseEventTypes(options[0].StringValue())
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			counts, err := database.GetUserEventCount(b.DB, eventTypes, i.GuildID)
			if err != nil {
				log.Println("Error getting user event count:", err)
				return
			}

			if len(counts) == 0 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "No events found for the specified event types.",
					},
				})
				return
			}

			var message strings.Builder
			index := 0
			for username, count := range counts {
				index++
				fmt.Fprintf(&message, "%d. %s: %d\n", index, username, count)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: message.String(),
				},
			})
		},
	}
}
