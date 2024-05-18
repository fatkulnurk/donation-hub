package user

import (
	"context"
	"github.com/isdzulqor/donation-hub/internal/core/entity"
	"github.com/isdzulqor/donation-hub/internal/core/model"
)

// TokenStorage for driver jwt
type TokenStorage interface {
}

// DataStorage for driven
type DataStorage interface {
	CreateUser(ctx context.Context, user model.UserRegisterInput) (entity.User, error)
	HasEmail(ctx context.Context, email string) (has bool, err error)
	HasUsername(ctx context.Context, username string) (has bool, err error)
	GetUserByUsername(ctx context.Context, username string) (user entity.User, err error)
	GetUser(ctx context.Context, input model.ListUserInput) (users []entity.User, total int64, err error)
	UserHasRole(ctx context.Context, userId int64, role string) (ok bool, err error)
}
