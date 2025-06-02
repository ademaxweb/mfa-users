package db

import (
	"errors"
	"github.com/ademaxweb/mfa-go-core/pkg/data"
)

type Interface interface {
	CreateUser(d data.User) (int, error)
	DeleteUser(id int) error
	UpdateUser(id int, d data.User) error
	GetUser(id int) (*data.User, error)
	GetAllUsers() ([]data.User, error)
	GetUserByEmail(email string) (*data.User, error)
}

var NotFound = errors.New("not found")
var NoFieldsToUpdate = errors.New("no fields to update")
