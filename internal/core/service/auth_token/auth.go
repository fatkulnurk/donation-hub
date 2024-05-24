package auth_token

import "github.com/isdzulqor/donation-hub/internal/core/model"

type Service interface {
	GenerateToken(i model.AuthPayload) (token string, err error)
	ValidateToken(token string) (*model.AuthPayload, error)
}
