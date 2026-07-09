package database

import (
	"database/sql"
	"strings"
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
VALUES (?, ?, ?, ?, ?, ?)
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

func GetUserEventCount(db *sql.DB, events []EventType, guildID string) (map[string]int, error) {
	query := `
SELECT COUNT(*) AS "count", username
FROM voice_events
WHERE guild_id = ?
AND event_type IN (?` + strings.Repeat(`, ?`, len(events)-1) + `)
GROUP BY username
ORDER BY count DESC
`

	args := make([]interface{}, 0, len(events)+1)
	args = append(args, guildID)
	for _, event := range events {
		args = append(args, event)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var count int
		var username string
		if err := rows.Scan(&count, &username); err != nil {
			return nil, err
		}
		result[username] = count
	}
	return result, nil
}
