package user

import (
	"context"
	"errors"
	"github.com/isdzulqor/donation-hub/internal/core/entity"
	"github.com/isdzulqor/donation-hub/internal/driver/rest/req"
)

type Storage struct {
	storage DataStorage
}

type Service interface {
	Register(ctx context.Context, rb req.RegisterReqBody) (user entity.User, err error)
	Login(ctx context.Context, rb req.LoginReqBody) (user entity.User, err error)
	Get(ctx context.Context) (users []entity.User, err error)
}

func NewService(storage DataStorage) Service {
	return &Storage{
		storage: storage,
	}
}

func (s *Storage) Register(ctx context.Context, rb req.RegisterReqBody) (user entity.User, err error) {
	userEntity := entity.User{
		Username: rb.Username,
		Email:    rb.Email,
		Password: rb.Password,
	}

	hasEmail, err := s.storage.HasEmail(ctx, userEntity.Email)
	if (err != nil) || (hasEmail) {
		err = errors.New("email already exists")
		return
	}

	hasUsername, err := s.storage.HasUsername(ctx, userEntity.Username)
	if (err != nil) || (hasUsername) {
		err = errors.New("username already exists")
		return
	}

	_ = s.storage.Store(ctx, &userEntity)

	return userEntity, err
}

func (s *Storage) Login(ctx context.Context, rb req.LoginReqBody) (user entity.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Get(ctx context.Context) (users []entity.User, err error) {
	//TODO implement me
	panic("implement me")
}
