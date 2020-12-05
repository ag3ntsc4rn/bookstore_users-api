package users

import (
	"fmt"

	"github.com/ag3ntsc4rn/bookstore_users-api/datasources/mysql/usersdb"
	"github.com/ag3ntsc4rn/bookstore_users-api/utils/errors"
	mysqlutils "github.com/ag3ntsc4rn/bookstore_users-api/utils/mysql"
)

const (
	queryInsertUser        = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES (?,?,?,?,?,?);"
	queryGetUser           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser        = "UPDATE users SET first_name=?,last_name=?,email=? WHERE id=?;"
	queryDeleteUser        = "DELETE FROM users where id=?;"
	queryFindUsersByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users where status=?;"
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
	if err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
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

	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
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
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		return mysqlutils.ParseError(err)
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(
			fmt.Sprintf("error occured while preparing statement to database: %v", err.Error()))
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.ID)
	if err != nil {
		return mysqlutils.ParseError(err)
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := usersdb.Client.Prepare(queryFindUsersByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(
			fmt.Sprintf("error occured while preparing statement to database: %v", err.Error()))
	}
	rows, err := stmt.Query(status)
	if err != nil {
		return nil, mysqlutils.ParseError(err)
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, mysqlutils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("No users found for status %v", status))
	}
	return results, nil
}
