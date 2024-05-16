package userstorage

import (
	"fmt"
	"github.com/isdzulqor/donation-hub/internal/core/entity"
	"github.com/jmoiron/sqlx"
	"time"
)

type Storage struct {
	sqlClient *sqlx.DB
}

func New(conn *sqlx.DB) *Storage {
	return &Storage{sqlClient: conn}
}

func (s Storage) CreateUser(user entity.User, role string) (entity.User, error) {
	ts := time.Now().Unix()
	query := `INSERT INTO users (username, email, password,created_at) VALUES (?,?,?,?)`
	resUser, err := s.sqlClient.Exec(query, user.Username, user.Email, user.Password, ts)
	if err != nil {
		return user, err
	}

	userId, _ := resUser.LastInsertId()
	query = `INSERT INTO user_roles (user_id, role) VALUES (?,?)`
	_, err = s.sqlClient.Exec(query, userId, role)

	return entity.User{
		ID:        userId,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: ts,
	}, nil
}

func (s Storage) HasEmail(email string) (has bool, err error) {
	query := "select count(*) from users where email = ?"
	var exists = false
	err = s.sqlClient.Get(&exists, query, email)

	fmt.Println("email: " + email)
	fmt.Println(exists)

	return exists, err
}

func (s Storage) HasUsername(username string) (has bool, err error) {
	query := "select count(*) from users where username = ?"
	var exists = false
	err = s.sqlClient.Get(&exists, query, username)

	fmt.Println("username: " + username)
	fmt.Println(exists)
	return exists, err
}
