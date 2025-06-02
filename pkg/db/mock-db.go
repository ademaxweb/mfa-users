package db

import (
	"github.com/ademaxweb/mfa-go-core/pkg/data"
	"reflect"
	"slices"
	"time"
)

type MockDB struct {
	users []data.User
}

func NewMockDB() *MockDB {
	return &MockDB{}
}

func (m *MockDB) Open() error {
	m.users = []data.User{
		{
			ID:         1,
			Name:       "Alex",
			Email:      "Alex@example.com",
			Password:   "123",
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		},
		{
			ID:         2,
			Name:       "Sam",
			Email:      "Sam@example.com",
			Password:   "123",
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		},
		{
			ID:         3,
			Name:       "Sarah",
			Email:      "Sarah@example.com",
			Password:   "123",
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		},
	}
	return nil
}

func (m *MockDB) Close() error {
	m.users = []data.User{}
	return nil
}

func (m *MockDB) CreateUser(d data.User) (int, error) {
	if len(m.users) < 1 {
		d.ID = 1
		m.users = []data.User{d}
		return d.ID, nil
	}

	lastUser := m.users[len(m.users)-1]
	d.ID = lastUser.ID + 1
	m.users = append(m.users, d)
	return d.ID, nil
}

func (m *MockDB) DeleteUser(id int) error {
	i := m.getUserIndex(id)
	if i == -1 {
		return NotFound
	}
	m.users = slices.Delete(m.users, i, i+1)
	return nil
}

func (m *MockDB) UpdateUser(id int, d data.User) error {
	i := m.getUserIndex(id)
	if i == -1 {
		return NotFound
	}
	target := &m.users[i]
	sourceValue := reflect.ValueOf(d)
	targetValue := reflect.ValueOf(target).Elem()

	for i := 0; i < sourceValue.NumField(); i++ {
		field := sourceValue.Field(i)
		if !field.IsZero() {
			targetValue.Field(i).Set(field)
		}
	}

	return nil
}

func (m *MockDB) GetUser(id int) (*data.User, error) {
	i := m.getUserIndex(id)
	if i == -1 {
		return nil, NotFound
	}
	return &m.users[i], nil
}

func (m *MockDB) GetAllUsers() ([]data.User, error) {
	return m.users, nil
}

func (m *MockDB) getUserIndex(id int) int {
	return slices.IndexFunc(m.users, func(u data.User) bool {
		return u.ID == id
	})
}

func (m *MockDB) getUserIndexByEmail(email string) int {
	return slices.IndexFunc(m.users, func(u data.User) bool {
		return u.Email == email
	})
}

func (m *MockDB) GetUserByEmail(email string) (*data.User, error) {
	i := m.getUserIndexByEmail(email)
	if i == -1 {
		return nil, NotFound
	}
	return &m.users[i], nil
}
