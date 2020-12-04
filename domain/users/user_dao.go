package users

import (
	"fmt"
	"strings"

	"github.com/ag3ntsc4rn/bookstore_users-api/datasources/mysql/usersdb"
	"github.com/ag3ntsc4rn/bookstore_users-api/utils/date"
	"github.com/ag3ntsc4rn/bookstore_users-api/utils/errors"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, date_created) VALUES (?,?,?,?);"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	if err := usersdb.Client.Ping(); err != nil {
		panic(err)
	}
	result := usersDB[user.ID]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user id %v not found", user.ID))
	}
	user.ID = result.ID
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(
			fmt.Sprintf("error occured while preparing statement to database: %v", err.Error()))
	}
	defer stmt.Close()
	user.DateCreated = date.GetNowString()
	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewBadRequestError(
				fmt.Sprintf("email %v already exists", user.Email))
		}
		return errors.NewInternalServerError(
			fmt.Sprintf("error occured while inserting to database:%v", err.Error()))
	}
	userID, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError("error when trying to get last inserted id")
	}
	user.ID = userID
	return nil
}
