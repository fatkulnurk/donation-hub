package user

import (
	"context"
	"github.com/isdzulqor/donation-hub/internal/core/entity"
)

type DataStorage interface {
	Store(ctx context.Context, user *entity.User) (err error)
	Get()
	GetById()
	HasEmail(ctx context.Context, email string) (has bool, err error)
	HasUsername(ctx context.Context, username string) (has bool, err error)
}
