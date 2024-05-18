package userstorage

import (
	"context"
	"fmt"
	"github.com/isdzulqor/donation-hub/internal/core/entity"
	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type Storage struct {
	sqlClient *sqlx.DB
}

func New(conn *sqlx.DB) *Storage {
	return &Storage{sqlClient: conn}
}

func (s Storage) CreateUser(ctx context.Context, input model.UserRegisterInput) (entity.User, error) {
	ts := time.Now().Unix()
	query := `INSERT INTO users (username, email, password,created_at) VALUES (?,?,?,?)`
	resUser, err := s.sqlClient.Exec(query, input.Username, input.Email, input.Password, ts)
	if err != nil {
		return entity.User{}, err
	}

	userId, _ := resUser.LastInsertId()
	query = `INSERT INTO user_roles (user_id, role) VALUES (?,?)`
	_, err = s.sqlClient.Exec(query, userId, input.Role)

	return entity.User{
		ID:        userId,
		Username:  input.Username,
		Email:     input.Email,
		Password:  input.Password,
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

func (s Storage) GetUserByUsername(username string) (user entity.User, err error) {
	query := "select * from users where username = ?"
	err = s.sqlClient.Get(&user, query, username)

	return
}

func (s Storage) GetUser(limit int, page int, role string) (users []entity.User, total int64, err error) {
	offset := (page - 1) * limit
	var query, queryCount string

	if role == "" {
		query = `SELECT users.*, GROUP_CONCAT(user_roles.role SEPARATOR ',') AS roles
				FROM users 
				JOIN user_roles ON users.id = user_roles.user_id
				WHERE user_roles.role IN ("donor", "requester")
				GROUP BY users.id LIMIT ? OFFSET ? `

		queryCount = `SELECT count(*)
				FROM users 
				JOIN user_roles ON users.id = user_roles.user_id
				WHERE user_roles.role IN ("donor", "requester")
				GROUP BY users.id LIMIT ? OFFSET ? `

		err = s.sqlClient.Select(&users, query, limit, offset)
		err = s.sqlClient.Get(&total, queryCount)
	} else {
		query = `SELECT users.*, GROUP_CONCAT(user_roles.role SEPARATOR ',') AS roles
				FROM users 
				JOIN user_roles ON users.id = user_roles.user_id
				WHERE user_roles.role = ? GROUP BY users.id LIMIT ? OFFSET ? `
		queryCount = `SELECT count(*)
				FROM users 
				JOIN user_roles ON users.id = user_roles.user_id
				WHERE user_roles.role = ? GROUP BY users.id LIMIT ? OFFSET ? `

		err = s.sqlClient.Select(&users, query, limit, offset)
		err = s.sqlClient.Get(&total, queryCount)
	}

	return
}
