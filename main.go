package main

import (
	"database/sql"
	"discord-voice-tracker/internal/database"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var trackedGuilds = make(map[string]bool)

type Bot struct {
	DB *sql.DB
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("BOT_TOKEN")
	guilds := os.Getenv("GUILDS")

	for guildID := range strings.SplitSeq(guilds, ",") {
		trackedGuilds[guildID] = true
	}

	db, err := database.Open("./data/voice.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bot := &Bot{
		DB: db,
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting bot...")

	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		panic(err)
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s", r.User.String())
	})

	session.AddHandler(bot.onVoiceStateUpdate)

	err = session.Open()
	if err != nil {
		log.Fatalf("could not open session: %s", err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = session.Close()
	if err != nil {
		log.Printf("could not close session gracefully: %s", err)
	}
}

func (b *Bot) onVoiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if !trackedGuilds[vs.GuildID] {
		return
	}

	var eventType string

	// Voice channel changes
	switch {
	case vs.BeforeUpdate == nil:
		eventType = database.EventTypeJoin
		fmt.Printf("%s joined voice channel %s\n", vs.Member.User.Username, vs.ChannelID)

	case vs.BeforeUpdate.ChannelID != "" && vs.ChannelID == "":
		eventType = database.EventTypeLeave
		fmt.Printf("%s left voice channel %s\n", vs.Member.User.Username, vs.BeforeUpdate.ChannelID)

	case vs.BeforeUpdate.ChannelID != vs.ChannelID:
		eventType = database.EventTypeMove
		fmt.Printf("%s moved from %s to %s\n",
			vs.Member.User.Username,
			vs.BeforeUpdate.ChannelID,
			vs.ChannelID,
		)

	// Speaker mute (self deafen)/undeafen
	case !vs.BeforeUpdate.SelfDeaf && vs.SelfDeaf:
		eventType = database.EventTypeSelfDeaf
		fmt.Printf("%s deafened themselves\n", vs.Member.User.GlobalName)

	case vs.BeforeUpdate.SelfDeaf && !vs.SelfDeaf:
		eventType = database.EventTypeSelfUndeaf
		fmt.Printf("%s undeafened themselves\n", vs.Member.User.GlobalName)

	// Microphone mute/unmute
	case !vs.BeforeUpdate.SelfMute && vs.SelfMute:
		eventType = database.EventTypeSelfMute
		fmt.Printf("%s muted their microphone\n", vs.Member.User.GlobalName)

	case vs.BeforeUpdate.SelfMute && !vs.SelfMute:
		eventType = database.EventTypeSelfUnmute
		fmt.Printf("%s unmuted their microphone\n", vs.Member.User.GlobalName)

	// Moderator mute
	case !vs.BeforeUpdate.Mute && vs.Mute:
		eventType = database.EventTypeServerMute
		fmt.Printf("%s was server muted\n", vs.Member.User.GlobalName)

	case vs.BeforeUpdate.Mute && !vs.Mute:
		eventType = database.EventTypeServerUnmute
		fmt.Printf("%s was server unmuted\n", vs.Member.User.GlobalName)

	// Moderator deafen
	case !vs.BeforeUpdate.Deaf && vs.Deaf:
		eventType = database.EventTypeServerDeaf
		fmt.Printf("%s was server deafened\n", vs.Member.User.GlobalName)

	case vs.BeforeUpdate.Deaf && !vs.Deaf:
		eventType = database.EventTypeServerUndeaf
		fmt.Printf("%s was server undeafened\n", vs.Member.User.GlobalName)
	}

	err := database.InsertVoiceEvent(b.DB, database.VoiceEvent{
		GuildID:   vs.GuildID,
		ChannelID: vs.ChannelID,
		UserID:    vs.UserID,
		Username:  vs.Member.User.GlobalName,
		EventType: eventType,
		Timestamp: time.Now().UTC(),
	})

	if err != nil {
		log.Println(err)
	}
}
