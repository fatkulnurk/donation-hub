package projectstorage

import (
	"context"
	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	sqlClient *sqlx.DB
}

func New(conn *sqlx.DB) *Storage {
	return &Storage{sqlClient: conn}
}

func (s Storage) Submit(ctx context.Context, input model.SubmitProjectInput) (err error) {
	panic("implement me")
}

func (s Storage) ReviewByAdmin(ctx context.Context, input model.ReviewProjectByAdminInput) (err error) {
	panic("implement me")
}

func (s Storage) ListProject(ctx context.Context, input model.ListProjectInput) (err error) {
	panic("implement me")
}

func (s Storage) GetProjectById(ctx context.Context, input model.GetProjectByIdInput) (err error) {
	panic("implement me")
}

func (s Storage) DonateToProject(ctx context.Context, input model.DonateToProjectInput) (err error) {
	panic("implement me")
}

func (s Storage) GetDonationById(ctx context.Context, input model.GetProjectByIdInput) (err error) {
	panic("implement me")
}
