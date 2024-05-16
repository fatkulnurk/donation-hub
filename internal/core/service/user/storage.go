package user

import (
	"github.com/isdzulqor/donation-hub/internal/core/entity"
)

// TokenStorage for driver jwt
type TokenStorage interface {
}

// DataStorage for driven
type DataStorage interface {
	CreateUser(user entity.User, role string) (entity.User, error)
	HasEmail(email string) (has bool, err error)
	HasUsername(username string) (has bool, err error)
	GetUserByUsername(username string) (user entity.User, err error)
	GetUser(limit int, page int, role string) (users []entity.User, total int64, err error)
}
