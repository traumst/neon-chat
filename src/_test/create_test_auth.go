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
// returns a number and an error.
//
// number represents either:
//   - positive number of created rows
//   - negative number of creation errors
//   - zero means - no changes in the db
//
// error can be composite of multiple errors.
func CreateTestAuth(dbConn *db.DBConn, testUsers config.TestUsers) (int, error) {
	log.Println("Checking test users status...")
	userNames := testUsers.GetNames()
	log.Println("Searching for test users", userNames)
	dbUsers, err := db.SearchUsers(dbConn.Conn, userNames)
	if err != nil {
		return 0, fmt.Errorf("failed to search for test users: %s", err)
	}
	var usersWithoutAuth []*db.User
	for i, dbUser := range dbUsers {
		dbUserAuth, _ := db.GetUserAuth(dbConn.Conn, dbUser.Id)
		if dbUserAuth == nil {
			usersWithoutAuth = append(usersWithoutAuth, dbUsers[i])
		}
	}
	insertCounter := 0
	errCounter := 0
	errs := make([]error, 0)
	for _, dbUser := range usersWithoutAuth {
		for _, testUser := range testUsers {
			if testUser.Name != dbUser.Name {
				continue
			}
			hash, err := utils.HashPassword(testUser.Pass, dbUser.Salt)
			if err != nil {
				log.Printf("ERROR failed hashing pass, %s", err)
				errs = append(errs, fmt.Errorf("failed hashing user[%s] pass[isEmpty:%t] with salt[isEmpty:%t], %s", testUser.Name, testUser.Pass == "", dbUser.Salt == "", err))
				break
			}
			_, err = db.AddAuth(dbConn.Conn, db.Auth{
				Id:     0,
				UserId: dbUser.Id,
				Type:   string(enum.AuthTypeEmail),
				Hash:   hash,
			})
			if err != nil {
				log.Printf("ERROR failed to create test auth for [%s]: %s", testUser.Name, err)
				errs = append(errs, fmt.Errorf("failed to create test auth for [%s]: %s", testUser.Name, err))
				errCounter += 1
			} else {
				log.Printf("test auth created for user[%s]", testUser.Name)
				insertCounter += 1
			}
			break
		}
	}
	log.Printf("Out of [%d] test auth - created [%d] - failed [%d]\n", len(testUsers), insertCounter, errCounter)
	if errCounter > 0 {
		err = fmt.Errorf("error creating auth for test users")
		for _, e := range errs {
			err = fmt.Errorf("%s, %s", err, e)
		}
		return -errCounter, err
	}
	return insertCounter, nil
}
