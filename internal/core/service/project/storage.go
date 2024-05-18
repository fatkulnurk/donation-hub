package project

import (
	"context"
	"github.com/isdzulqor/donation-hub/internal/core/model"
)

type FileStorage interface {
	RequestUploadUrl(mimeType string, fileSize int64) (url string, expiredAt int64, err error)
}

type DataStorage interface {
	Submit(ctx context.Context, input model.SubmitProjectInput) (err error)
	ReviewByAdmin(ctx context.Context, input model.ReviewProjectByAdminInput) (err error)
	ListProject(ctx context.Context, input model.ListProjectInput) (err error)
	GetProjectById(ctx context.Context, input model.GetProjectByIdInput) (err error)
	DonateToProject(ctx context.Context, input model.DonateToProjectInput) (err error)
	GetDonationById(ctx context.Context, input model.GetProjectByIdInput) (err error)
}
