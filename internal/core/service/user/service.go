package user

import (
	"errors"
	"github.com/isdzulqor/donation-hub/internal/core/entity"
	"github.com/isdzulqor/donation-hub/internal/driver/rest/req"
)

type Storage struct {
	storage DataStorage
}

type Service interface {
	Register(rb req.RegisterReqBody) (user entity.User, err error)
	Login(username string, password string) (user entity.User, accessToken string, err error)
	ListUser(limit int, page int, role string) (users []entity.User, totalPage int64, err error)
}

func NewService(storage DataStorage) Service {
	return &Storage{
		storage: storage,
	}
}

func (s *Storage) Register(rb req.RegisterReqBody) (user entity.User, err error) {
	hasEmail, err := s.storage.HasEmail(rb.Email)
	if (err != nil) || (hasEmail) {
		err = errors.New("email already exists")
		return
	}

	hasUsername, err := s.storage.HasUsername(rb.Username)
	if (err != nil) || (hasUsername) {
		err = errors.New("username already exists")
		return
	}
	user = entity.User{
		Username: rb.Username,
		Password: rb.Password,
		Email:    rb.Email,
	}
	user, err = s.storage.CreateUser(user, rb.Role)
	return user, err
}

func (s *Storage) Login(username string, password string) (user entity.User, accessToken string, err error) {
	user, err = s.storage.GetUserByUsername(username)
	if err != nil || user.Password != password {
		err = errors.New("invalid username or password")
		return
	}

	// todo generate access token

	return
}

func (s *Storage) ListUser(limit int, page int, role string) (users []entity.User, totalPage int64, err error) {
	users, totalPage, err = s.storage.GetUser(limit, page, role)
	if err != nil {
		return nil, 0, err
	}

	return
}
