package projectstorage

import "github.com/jmoiron/sqlx"

type Storage struct {
	sqlClient *sqlx.DB
}

func New(conn *sqlx.DB) *Storage {
	return &Storage{sqlClient: conn}
}

func (s Storage) Submit() (err error) {
	panic("implement me")
}

func (s Storage) ReviewByAdmin() (err error) {
	panic("implement me")
}

func (s Storage) Get() (err error) {
	panic("implement me")
}

func (s Storage) GetById() (err error) {
	panic("implement me")
}

func (s Storage) DonateToProject() (err error) {
	panic("implement me")
}

func (s Storage) GetDonationById() (err error) {
	panic("implement me")
}
