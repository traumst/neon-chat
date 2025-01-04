package db

import (
	"fmt"

	"neon-chat/src/consts"

	"github.com/jmoiron/sqlx"
)

type Sentiment struct {
	ItemType string               `db:"item_type"`
	ItemId   uint                 `db:"item_id"`
	ChatId   uint                 `db:"chat_id"`
	UserId   uint                 `db:"user_id"`
	Type     consts.SentimentType `db:"score_type"`
	Value    float64              `db:"score_value"`
}

const SentimentSchema = `
	CREATE TABLE IF NOT EXISTS sentiment (
		item_type TEXT,
		item_id INTEGER, 
		chat_id INTEGER FOR, 
		user_id INTEGER, 
		score_type TEXT, 
		score_value REAL,
		FOREIGN KEY(chat_id) REFERENCES chats(id) ON DELETE CASCADE,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

const SentimentIndex = `
	CREATE INDEX IF NOT EXISTS idx_sentiment_fragment ON sentiment(user_id, chat_id);
	CREATE INDEX IF NOT EXISTS idx_sentiment_fragment ON sentiment(item_type, item_id, score_type);
	CREATE INDEX IF NOT EXISTS idx_sentiment_value ON sentiment(score_value);`

func (dbConn *DBConn) SentimentTableExists() bool {
	return dbConn.TableExists("sentiment")
}

func AddSentiment(dbConn sqlx.Ext, fragment Sentiment) error {
	if fragment.ItemType == "" {
		return fmt.Errorf("sentiment is missing itemType")
	} else if fragment.ItemId == 0 {
		return fmt.Errorf("sentiment is missing itemId")
	} else if fragment.Type == "" {
		return fmt.Errorf("sentiment is missing scoreType")
	}

	result, err := dbConn.Exec(`INSERT INTO sentiment (item_type, item_id, chat_id, user_id, score_type, score_value) VALUES (?, ?, ?, ?, ?, ?)`,
		fragment.ItemType, fragment.ItemId, fragment.ChatId, fragment.UserId, fragment.Type, fragment.Value)
	if err != nil {
		return fmt.Errorf("error adding sentiment: %s", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %s", err)
	} else if rows == 0 {
		return fmt.Errorf("0 rows affected")
	}

	return nil
}
