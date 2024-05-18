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
	HasEmail(email string) (has bool, err error)
	HasUsername(username string) (has bool, err error)
	GetUserByUsername(username string) (user entity.User, err error)
	GetUser(limit int, page int, role string) (users []entity.User, total int64, err error)
}
