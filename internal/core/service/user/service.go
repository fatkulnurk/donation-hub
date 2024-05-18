package user

import (
	"context"
	"errors"
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

	user, err := s.storage.CreateUser(ctx, input)
	output.ID = user.ID
	output.Username = user.Username
	output.Email = user.Email

	return
}

func (s *Storage) Login(ctx context.Context, input model.UserLoginInput) (output model.UserLoginOutput, err error) {
	user, err := s.storage.GetUserByUsername(input.Username)
	if err != nil || user.Password != input.Password {
		err = errors.New("invalid username or password")
		return
	}

	// assign output
	output.ID = user.ID
	output.Email = user.Email
	output.Username = user.Username
	output.AccessToken = ""

	// todo generate access token
	panic("implement me")
	return
}

func (s *Storage) ListUser(ctx context.Context, input model.ListUserInput) (output model.ListUserOutput, err error) {
	_, _, err = s.storage.GetUser(int(input.Limit), int(input.Page), input.Role)
	if err != nil {
		return model.ListUserOutput{}, err
	}

	// todo parse output to listUser

	return
}
