package users

import (
	"fmt"

	"github.com/ag3ntsc4rn/bookstore_users-api/datasources/mysql/usersdb"
	"github.com/ag3ntsc4rn/bookstore_users-api/utils/date"
	"github.com/ag3ntsc4rn/bookstore_users-api/utils/errors"
	mysqlutils "github.com/ag3ntsc4rn/bookstore_users-api/utils/mysql"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES (?,?,?,?);"
	queryGetUser    = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
	queryUpdateUser = "UPDATE users SET first_name=?,last_name=?,email=? WHERE id=?;"
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
		return mysqlutils.ParseError(err)
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
		return mysqlutils.ParseError(err)
	}
	userID, err := insertResult.LastInsertId()
	if err != nil {
		return mysqlutils.ParseError(err)
	}
	user.ID = userID
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(
			fmt.Sprintf("error occured while preparing statement to database: %v", err.Error()))
	}
	defer stmt.Close()
	user.DateCreated = date.GetNowString()
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		return mysqlutils.ParseError(err)
	}
	return nil
}
