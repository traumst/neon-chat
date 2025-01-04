package pub

import (
	"fmt"
	"log"

	"neon-chat/src/app"
	"neon-chat/src/db"
)

// This function
//   - checks if both reporter and reported user are members of the specified chat,
//   - parses the sentiment type, and adds the sentiment to the database.
//   - if any step fails, it returns an appropriate error.
//
// Parameters:
//   - dbConn: Database connection object.
//   - reporter: The user reporting the sentiment.
//   - userId: The ID of the user to whom the sentiment is directed.
//   - chatId: The ID of the chat containing the message.
//   - msgId: The ID of the message to which the sentiment is added.
//   - class: The type of sentiment being added.
//   - subText: Additional text explaining the sentiment.
//
// Returns:
//   - error: An error object if an error occurs, otherwise nil.
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
	sentimentType, err := db.ParseSentimentType(sentiment)
	if err != nil || sentimentType == "" {
		return fmt.Errorf("unexpected sentiment type[%s], %s", sentiment, err)
	}
	_, err = db.AddSentiment(dbConn.Tx, db.Sentiment{
		ItemType: "message",
		ItemId:   msgId,
		ChatId:   chatId,
		UserId:   userId,
		Type:     sentimentType,
		Value:    1.000, // user input, source of truth
	})
	if err != nil {
		return fmt.Errorf("failed to add message to chat[%d]: %s", chatId, err.Error())
	}
	return nil
}
