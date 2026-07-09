package main

import (
	"discord-voice-tracker/internal/database"
	"fmt"
	"log"
	"sort"
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

type UserCount struct {
	Username string
	Count    int
}

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

			var sorted []UserCount
			for username, count := range counts {
				sorted = append(sorted, UserCount{
					Username: username,
					Count:    count,
				})
			}

			// Sort by count descending
			sort.Slice(sorted, func(i, j int) bool {
				if sorted[i].Count == sorted[j].Count {
					return sorted[i].Username < sorted[j].Username
				}
				return sorted[i].Count > sorted[j].Count
			})

			var message strings.Builder
			for i, user := range sorted {
				fmt.Fprintf(&message, "%d. %s: %d\n", i+1, user.Username, user.Count)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: message.String(),
				},
			})
		},
	}
}
