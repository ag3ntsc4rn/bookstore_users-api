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
	noRowsInResultSet = "no rows in result set"
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, date_created) VALUES (?,?,?,?);"
	queryGetUser     = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(
			fmt.Sprintf("error occured while preparing statement to database: %v", err.Error()))
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.ID)
	if err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		if strings.Contains(err.Error(), noRowsInResultSet) {
			return errors.NewNotFoundError(fmt.Sprintf("no user found with id %v", user.ID))
		}
		return errors.NewInternalServerError(
			fmt.Sprintf("error occured while retreiving result from database for user id %v:%v", user.ID, err.Error()))
	}
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
