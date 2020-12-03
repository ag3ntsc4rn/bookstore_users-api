package services

import (
	"github.com/ag3ntsc4rn/bookstore_users-api/domain/users"
	"github.com/ag3ntsc4rn/bookstore_users-api/utils/errors"
)

func CreateUser(u users.User) (*users.User, *errors.RestErr) {
	return &u, nil
}
