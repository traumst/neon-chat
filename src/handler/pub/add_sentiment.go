package pub

import (
	"fmt"
	"log"

	"neon-chat/src/app"
	"neon-chat/src/db"
)

func AddSentiment(
	dbConn *db.DBConn,
	reporter *app.User,
	userId uint,
	chatId uint,
	msgId uint,
	sentiment string,
	subText string,
) error {
	log.Printf("TRACE AddSentiment user[%d] flags chat[%d] msg[%d] as [%s] because of [%s]\n",
		reporter.Id, chatId, msgId, sentiment, subText)

	canChat, err := db.UsersCanChat(dbConn.Tx, chatId, reporter.Id)
	if err != nil {
		return fmt.Errorf("failed to check user[%d] can chat[%d]: %s", reporter.Id, chatId, err.Error())
	}
	if !canChat {
		return fmt.Errorf("user is not in chat")
	}
	st := db.SentimentType(sentiment)
	if st == "" {
		return fmt.Errorf("unexpected sentiment type[%s]", sentiment)
	}
	_, err = db.AddSentiment(dbConn.Tx, db.Sentiment{
		ItemType: "message",
		ItemId:   msgId,
		ChatId:   chatId,
		UserId:   userId,
		Type:     db.SentimentType(sentiment),
		Value:    1.000,
	})
	if err != nil {
		return fmt.Errorf("failed to add message to chat[%d]: %s", chatId, err.Error())
	}
	return nil
}
