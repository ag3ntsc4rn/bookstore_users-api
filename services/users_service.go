package services

import (
	"github.com/ag3ntsc4rn/bookstore_users-api/domain/users"
	"github.com/ag3ntsc4rn/bookstore_users-api/utils/crypto"
	"github.com/ag3ntsc4rn/bookstore_users-api/utils/date"
	"github.com/ag3ntsc4rn/bookstore_users-api/utils/errors"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersServiceInterface interface {
	CreateUser(users.User) (*users.User, *errors.RestErr)
	GetUser(int64) (*users.User, *errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, *errors.RestErr)
	DeleteUser(int64) *errors.RestErr
	SearchUser(string) (users.Users, *errors.RestErr)
}

type usersService struct {
}

func (s *usersService) CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.DateCreated = date.GetNowDBFormat()
	user.Password = crypto.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) GetUser(userID int64) (*users.User, *errors.RestErr) {
	result := &users.User{
		ID: userID,
	}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {
	current, err := s.GetUser(user.ID)
	if err != nil {
		return nil, err
	}
	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}

	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *usersService) DeleteUser(userID int64) *errors.RestErr {
	current, err := s.GetUser(userID)
	if err != nil {
		return err
	}
	if err = current.Delete(); err != nil {
		return err
	}

	return nil

}

func (s *usersService) SearchUser(status string) (users.Users, *errors.RestErr) {
	user := &users.User{}
	return user.FindByStatus(status)
}
