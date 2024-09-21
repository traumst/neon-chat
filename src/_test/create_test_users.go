package test

import (
	"fmt"
	"log"
	"neon-chat/src/app/enum"
	"neon-chat/src/db"
	"neon-chat/src/utils"
	"neon-chat/src/utils/config"
)

// attempts to create test users that do not exist in the database.
// skips errors to attempt to create other users.
// returns a number and an error,
// error can be composite of multiple errors,
// number represents
//   - positive number of created rows
//   - negative number of creation errors
func CreateTestUsers(dbConn *db.DBConn, testUsers config.TestUsers) (int, error) {
	log.Println("Checking test users status...")
	dbUsers, _ := db.SearchUsers(dbConn.Conn, testUsers.GetNames())
	if dbUsers == nil {
		dbUsers = make([]*db.User, 0)
	}
	newUsers := make(config.TestUsers, 0)
	for _, testUser := range testUsers {
		exists := false
		for _, dbUser := range dbUsers {
			if dbUser.Name == testUser.Name {
				exists = true
				break
			}
		}
		if !exists {
			newUsers = append(newUsers, testUser)
		}
	}
	log.Printf("...there are [%d] test users to create\n", len(newUsers))
	insertCounter := 0
	errCounter := 0
	errs := make([]error, 0)
	for _, testUser := range newUsers {
		dbUser := db.User{
			Id:     0,
			Name:   testUser.Name,
			Email:  testUser.Email,
			Type:   string(enum.UserTypeBasic),
			Status: string(enum.UserStatusActive),
			Salt:   utils.GenerateSalt(testUser.Name, string(enum.UserTypeBasic)),
		}
		_, err := db.AddUser(dbConn.Conn, &dbUser)
		if err != nil {
			log.Printf("failed to create test user[%s]: %s", testUser.Name, err)
			errs = append(errs, err)
			errCounter += 1
		} else {
			log.Printf("test user created[%s]", testUser.Name)
			insertCounter += 1
		}
	}
	log.Printf("Out of [%d] test users - created [%d] - failed [%d]\n", len(testUsers), insertCounter, errCounter)
	if errCounter > 0 {
		err := fmt.Errorf("error creating test users")
		for _, e := range errs {
			err = fmt.Errorf("%s, %s", err, e)
		}
		return -errCounter, err
	}
	return insertCounter, nil
}
