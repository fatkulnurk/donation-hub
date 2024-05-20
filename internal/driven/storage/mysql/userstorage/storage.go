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

func (s Storage) HasEmail(ctx context.Context, email string) (has bool, err error) {
	query := "select count(*) from users where email = ?"
	var exists = false
	err = s.sqlClient.Get(&exists, query, email)

	fmt.Println("email: " + email)
	fmt.Println(exists)

	return exists, err
}

func (s Storage) HasUsername(ctx context.Context, username string) (has bool, err error) {
	query := "select count(*) from users where username = ?"
	var exists = false
	err = s.sqlClient.Get(&exists, query, username)

	fmt.Println("username: " + username)
	fmt.Println(exists)
	return exists, err
}

func (s Storage) GetUserByUsername(ctx context.Context, username string) (user entity.User, err error) {
	query := "select * from users where username = ?"
	err = s.sqlClient.Get(&user, query, username)

	return
}

type DatabaseUser struct {
	ID        int64  `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"password"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
	Roles     string `db:"roles" json:"roles"`
}

type UsersCount struct {
	Total int64 `json:"total"`
}

// GetUser total is a total data, not pagination
func (s Storage) GetUser(ctx context.Context, input model.ListUserInput) (users []entity.User, total int64, err error) {
	offset := (input.Page - 1) * input.Limit
	var query string
	var count UsersCount

	if input.Role == "" {
		query = `SELECT users.*, GROUP_CONCAT(user_roles.role) AS roles
				FROM users 
				JOIN user_roles ON users.id = user_roles.user_id
				WHERE user_roles.role IN ("donor", "requester")
				GROUP BY users.id LIMIT ? OFFSET ? `

		err = s.sqlClient.Select(&users, query, input.Limit, offset)
		err = s.sqlClient.Get(&count, "SELECT COUNT(*) as total FROM users u JOIN user_roles ur ON u.id = ur.user_id WHERE ur.role IN ('donor', 'requester')")
	} else {
		query = `SELECT users.*, GROUP_CONCAT(user_roles.role) AS roles
				FROM users 
				JOIN user_roles ON users.id = user_roles.user_id
				WHERE user_roles.role = ? GROUP BY users.id LIMIT ? OFFSET ? `

		err = s.sqlClient.Select(&users, query, input.Role, input.Limit, offset)
		err = s.sqlClient.Get(&count, "SELECT COUNT(*) as total FROM users u JOIN user_roles ur ON u.id = ur.user_id WHERE ur.role = ? GROUP BY u.id")
	}

	total = count.Total
	return
}

func (s Storage) UserHasRole(ctx context.Context, userId int64, role string) (ok bool, err error) {
	query := "select count(*) from user_roles where user_id = ? and role = ?"
	var exists = false
	err = s.sqlClient.Get(&exists, query, userId, role)

	return exists, err
}
