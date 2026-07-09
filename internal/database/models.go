package database

import "time"

var (
	EventTypeJoin         = "join"
	EventTypeLeave        = "leave"
	EventTypeMove         = "move"
	EventTypeSelfDeaf     = "self_deaf"
	EventTypeSelfUndeaf   = "self_undeaf"
	EventTypeSelfMute     = "self_mute"
	EventTypeSelfUnmute   = "self_unmute"
	EventTypeServerMute   = "server_mute"
	EventTypeServerUnmute = "server_unmute"
	EventTypeServerDeaf   = "server_deaf"
	EventTypeServerUndeaf = "server_undeaf"
)

type VoiceEvent struct {
	GuildID   string
	ChannelID string
	UserID    string
	Username  string

	EventType string

	Timestamp time.Time
}
