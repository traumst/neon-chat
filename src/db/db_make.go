package db

import (
	"log"
)

// TODO: update when adding new tables
func (dbConn *DBConn) concatSchema() (schema string, index string) {
	if !dbConn.MigrationsTableExists() {
		log.Println("TRACE concatSchema migration table will be created")
		schema += MigrationSchema + "\n"
		index += MigrationIndex + "\n"
	}

	if !dbConn.UserTableExists() {
		log.Println("TRACE concatSchema user table will be created")
		schema += UserSchema + "\n"
		index += UserIndex + "\n"
	}

	if !dbConn.AuthTableExists() {
		log.Println("TRACE concatSchema auth table will be created")
		schema += AuthSchema + "\n"
		index += AuthIndex + "\n"
	}

	if !dbConn.AvatarTableExists() {
		log.Println("TRACE concatSchema avatar table will be created")
		schema += AvatarSchema + "\n"
		index += AvatarIndex + "\n"
	}

	if !dbConn.ReservationTableExists() {
		log.Println("TRACE concatSchema reservation table will be created")
		schema += ReservationSchema + "\n"
		index += ReservationIndex + "\n"
	}

	if !dbConn.ChatTableExists() {
		log.Println("TRACE concatSchema chat table will be created")
		schema += ChatSchema + "\n"
		index += ChatIndex + "\n"
	}

	if !dbConn.ChatUserTableExists() {
		log.Println("TRACE concatSchema chat_user table will be created")
		schema += ChatUserSchema + "\n"
		index += ChatUserIndex + "\n"
	}

	if !dbConn.MessageTableExists() {
		log.Println("TRACE concatSchema messages table will be created")
		schema += MessageSchema + "\n"
		index += MessageIndex + "\n"
	}

	if !dbConn.QuoteTableExists() {
		log.Println("TRACE concatSchema quotes table will be created")
		schema += QuoteSchema + "\n"
		index += QuoteIndex + "\n"
	}

	if !dbConn.SentimentTableExists() {
		log.Println("TRACE concatSchema sentiment table will be created")
		schema += SentimentSchema + "\n"
		index += SentimentIndex + "\n"
	}

	return schema, index
}
