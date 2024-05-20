package user

import (
	"context"
	"errors"
	"github.com/isdzulqor/donation-hub/internal/core/model"
	"math"
	"strings"
)

type Storage struct {
	storage DataStorage
}

type Service interface {
	Register(ctx context.Context, input model.UserRegisterInput) (output model.UserRegisterOutput, err error)
	Login(ctx context.Context, input model.UserLoginInput) (output model.UserLoginOutput, err error)
	ListUser(ctx context.Context, input model.ListUserInput) (output *model.ListUserOutput, err error)
}

func NewService(storage DataStorage) Service {
	return &Storage{
		storage: storage,
	}
}

func (s *Storage) Register(ctx context.Context, input model.UserRegisterInput) (output model.UserRegisterOutput, err error) {
	hasEmail, err := s.storage.HasEmail(ctx, input.Email)
	if (err != nil) || (hasEmail) {
		err = errors.New("email already exists")
		return
	}

	hasUsername, err := s.storage.HasUsername(ctx, input.Username)
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
	user, err := s.storage.GetUserByUsername(ctx, input.Username)
	if err != nil || user.Password != input.Password {
		err = errors.New("invalid username or password")
		return
	}

	// assign output
	output.ID = user.ID
	output.Email = user.Email
	output.Username = user.Username

	// todo generate access token
	panic("implement me")
	output.AccessToken = ""
	return
}

func (s *Storage) ListUser(ctx context.Context, input model.ListUserInput) (output *model.ListUserOutput, err error) {
	users, total, err := s.storage.GetUser(ctx, input)
	if err != nil {
		return nil, err
	}

	listUsers := make([]model.ListUser, len(users))
	for i, user := range users {
		roles := strings.Split(user.Roles, ",")
		listUser := model.ListUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Roles:    roles,
		}
		listUsers[i] = listUser
	}

	// pagination
	totalPage := int64(math.Ceil(float64(total / input.Limit)))
	if total%input.Limit != 0 {
		totalPage++
	}

	output = &model.ListUserOutput{
		Users: listUsers,
		Pagination: model.ListUserMeta{
			Page:       input.Page,
			TotalPages: totalPage,
		},
	}

	return
}
