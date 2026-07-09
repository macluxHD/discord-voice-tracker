package database

import "database/sql"

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS voice_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    guild_id TEXT NOT NULL,
    channel_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    username TEXT NOT NULL,

    event_type TEXT NOT NULL,

    timestamp TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_user
ON voice_events(user_id);

CREATE INDEX IF NOT EXISTS idx_events_time
ON voice_events(timestamp);

CREATE INDEX IF NOT EXISTS idx_events_guild
ON voice_events(guild_id);
`)
	return err
}
