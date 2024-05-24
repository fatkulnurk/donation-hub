package project

import (
	"context"
	"github.com/isdzulqor/donation-hub/internal/core/model"
)

type FileStorage interface {
	RequestUploadUrl(mimeType string, fileSize int64) (url string, expiredAt int64, err error)
}

type DataStorage interface {
	Submit(ctx context.Context, input model.SubmitProjectInput) (projectId int64, err error)
	ReviewByAdmin(ctx context.Context, input model.ReviewProjectByAdminInput) (err error)
	ListProject(ctx context.Context, input model.ListProjectInput) (o model.ListProjectOutput, err error)
	GetProjectById(ctx context.Context, input model.GetProjectByIdInput) (o model.GetProjectByIdOutput, err error)
	DonateToProject(ctx context.Context, input model.DonateToProjectInput) (err error)
	ListDonationByProjectId(ctx context.Context, input model.ListProjectDonationInput) (o model.ListProjectDonationOutput, err error)
}
