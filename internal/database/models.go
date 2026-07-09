package database

import (
	"fmt"
	"strings"
	"time"
)

type EventType string

const (
	EventTypeJoin         EventType = "join"
	EventTypeLeave        EventType = "leave"
	EventTypeMove         EventType = "move"
	EventTypeSelfDeaf     EventType = "self_deaf"
	EventTypeSelfUndeaf   EventType = "self_undeaf"
	EventTypeSelfMute     EventType = "self_mute"
	EventTypeSelfUnmute   EventType = "self_unmute"
	EventTypeServerMute   EventType = "server_mute"
	EventTypeServerUnmute EventType = "server_unmute"
	EventTypeServerDeaf   EventType = "server_deaf"
	EventTypeServerUndeaf EventType = "server_undeaf"
)

func (e EventType) Valid() bool {
	switch e {
	case EventTypeJoin,
		EventTypeLeave,
		EventTypeMove,
		EventTypeSelfDeaf,
		EventTypeSelfUndeaf,
		EventTypeSelfMute,
		EventTypeSelfUnmute,
		EventTypeServerMute,
		EventTypeServerUnmute,
		EventTypeServerDeaf,
		EventTypeServerUndeaf:
		return true
	default:
		return false
	}
}

func ParseEventTypes(s string) ([]EventType, error) {
	var out []EventType

	for part := range strings.SplitSeq(s, ",") {
		e := EventType(strings.TrimSpace(part))
		if !e.Valid() {
			return nil, fmt.Errorf("invalid event type %q", part)
		}
		out = append(out, e)
	}

	return out, nil
}

type VoiceEvent struct {
	GuildID   string
	ChannelID string
	UserID    string
	Username  string

	EventType EventType

	Timestamp time.Time
}
