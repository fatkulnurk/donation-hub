package userstorage

import "github.com/jmoiron/sqlx"

type User struct {
	ID        uint64 `db:"id"`
	Username  string `db:"username"`
	Email     string `db:"email"`
	Password  string `db:"-"`
	CreatedAt uint32 `db:"created_at"`
}

type Storage struct {
	connection *sqlx.DB
}

func New(conn *sqlx.DB) *Storage {
	return &Storage{connection: conn}
}

func (s *Storage) Register() {

}

func (s *Storage) FindByEmail() {

}

func (s *Storage) FindByUsername() {

}

func (s *Storage) Lists() {

}
