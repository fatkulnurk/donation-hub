package user

import (
	"context"
	"errors"
	"github.com/isdzulqor/donation-hub/internal/core/entity"
	"github.com/isdzulqor/donation-hub/internal/core/model"
)

type Storage struct {
	storage DataStorage
}

type Service interface {
	Register(ctx context.Context, input model.UserRegisterInput) (output model.UserRegisterOutput, err error)
	Login(ctx context.Context, input model.UserLoginInput) (output model.UserLoginOutput, err error)
	ListUser(ctx context.Context, input model.ListUserInput) (output model.ListUserOutput, err error)
}

func NewService(storage DataStorage) Service {
	return &Storage{
		storage: storage,
	}
}

func (s *Storage) Register(ctx context.Context, input model.UserRegisterInput) (output model.UserRegisterOutput, err error) {
	hasEmail, err := s.storage.HasEmail(input.Email)
	if (err != nil) || (hasEmail) {
		err = errors.New("email already exists")
		return
	}

	hasUsername, err := s.storage.HasUsername(input.Username)
	if (err != nil) || (hasUsername) {
		err = errors.New("username already exists")
		return
	}
	user := entity.User{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	}
	user, err = s.storage.CreateUser(user, input.Role)
	return model.UserRegisterOutput{}, err
}

func (s *Storage) Login(ctx context.Context, input model.UserLoginInput) (output model.UserLoginOutput, err error) {
	user, err = s.storage.GetUserByUsername(username)
	if err != nil || user.Password != password {
		err = errors.New("invalid username or password")
		return
	}

	// todo generate access token
	panic("implement me")
	return
}

func (s *Storage) ListUser(ctx context.Context, input model.ListUserInput) (output model.ListUserOutput, err error) {
	users, totalPage, err = s.storage.GetUser(limit, page, role)
	if err != nil {
		return nil, 0, err
	}

	return
}
