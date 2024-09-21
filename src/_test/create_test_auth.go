package test

import (
	"fmt"
	"log"
	"neon-chat/src/app/enum"
	"neon-chat/src/db"
	"neon-chat/src/utils"
	"neon-chat/src/utils/config"
)

// attempts to create auth for test users.
// skips errors to attempt to create other auth.
// returns a number and an error,
// error can be composite of multiple errors,
// number represents either:
//   - zero means - no changes in the db
//   - positive number of created rows
//   - negative number of creation errors
func CreateTestAuth(dbConn *db.DBConn, TestUsers config.TestUsers) (int, error) {
	log.Println("Checking test users status...")
	dbUsers, err := db.SearchUsers(dbConn.Conn, TestUsers.GetNames()...)
	if err != nil {
		return 0, fmt.Errorf("failed to search for test users: %s", err)
	}
	insertCounter := 0
	errCounter := 0
	errs := make([]error, 0)
	for _, dbUser := range dbUsers {
		for _, testUser := range TestUsers {
			if testUser.Name == dbUser.Name {
				hash, err := utils.HashPassword(testUser.Pass, testUser.Salt)
				if err != nil {
					log.Printf("ERROR failed hashing pass, %s", err)
				}
				_, err = db.AddAuth(dbConn.Conn, db.Auth{
					Id:     0,
					UserId: dbUser.Id,
					Type:   string(enum.AuthTypeEmail),
					Hash:   hash,
				})
				if err != nil {
					log.Printf("ERROR failed to create test auth for [%s]: %s", testUser.Name, err)
					errs = append(errs, err)
					errCounter += 1
				} else {
					log.Printf("test auth created for user[%s]", testUser.Name)
					insertCounter += 1
				}
			}
		}
	}
	log.Printf("Out of [%d] test auth - created [%d] - failed [%d]\n", len(TestUsers), insertCounter, errCounter)
	if errCounter > 0 {
		err = fmt.Errorf("error creating auth for test users")
		for _, e := range errs {
			err = fmt.Errorf("%s, %s", err, e)
		}
		return -errCounter, err
	}
	return insertCounter, nil
}
