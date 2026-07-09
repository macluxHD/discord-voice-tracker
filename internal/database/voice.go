package database

import (
	"database/sql"
)

func InsertVoiceEvent(db *sql.DB, e VoiceEvent) error {
	_, err := db.Exec(`
INSERT INTO voice_events (
    guild_id,
    channel_id,
    user_id,
	username,
    event_type,
    timestamp
)
VALUES (?, ?, ?, ?, ?)
`,
		e.GuildID,
		e.ChannelID,
		e.UserID,
		e.Username,
		e.EventType,
		e.Timestamp,
	)

	return err
}
