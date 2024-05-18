package project

import (
	"context"
	"errors"
	"github.com/isdzulqor/donation-hub/internal/core/model"
)

type Storage struct {
	storage     DataStorage
	fileStorage FileStorage
}

type Service interface {
	RequestUploadUrl(ctx context.Context, input model.RequestUploadUrlInput) (output model.RequestUploadUrlOutput, err error)
	SubmitProject(ctx context.Context, input model.SubmitProjectInput) (output model.SubmitProjectOutput, err error)
	ReviewProjectByAdmin(ctx context.Context, input model.ReviewProjectByAdminInput) (ok bool, err error)
	ListProject(ctx context.Context, input model.ListProjectInput) (output model.ListProjectOutput, err error)
	GetProjectById(ctx context.Context, input model.GetProjectByIdInput) (output model.GetProjectByIdOutput, err error)
	DonateToProject(ctx context.Context, input model.DonateToProjectInput) (ok bool, err error)
	ListProjectDonationById(ctx context.Context, input model.ListProjectDonationInput) (output model.ListProjectDonationOutput, err error)
}

func NewService(storage DataStorage, fileStorage FileStorage) Service {
	return &Storage{
		storage:     storage,
		fileStorage: fileStorage,
	}
}

func (s Storage) RequestUploadUrl(ctx context.Context, input model.RequestUploadUrlInput) (output model.RequestUploadUrlOutput, err error) {
	// validate user, make sure role is valid

	// validate size
	if input.FileSize > 1048576 {
		err = errors.New("filesize can't greater than 1MB")
		return
	}

	// validate mimetype
	if input.MimeType != "image/jpeg" && input.MimeType != "image/png" {
		err = errors.New("mimetype must be image/jpg or image/png")
		return
	}

	url, expiredAt, err := s.fileStorage.RequestUploadUrl(input.MimeType, input.FileSize)

	// assign to struct
	output.URL = url
	output.FileSize = input.FileSize
	output.MimeType = input.MimeType
	output.ExpiresAt = expiredAt

	return
}

func (s Storage) SubmitProject(ctx context.Context, input model.SubmitProjectInput) (output model.SubmitProjectOutput, err error) {
	// validate user, make sure role is valid

	// save to database
	_ = s.storage.Submit(ctx, input)
	panic("implement me")
}

func (s Storage) ReviewProjectByAdmin(ctx context.Context, input model.ReviewProjectByAdminInput) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) ListProject(ctx context.Context, input model.ListProjectInput) (output model.ListProjectOutput, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetProjectById(ctx context.Context, input model.GetProjectByIdInput) (output model.GetProjectByIdOutput, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) DonateToProject(ctx context.Context, input model.DonateToProjectInput) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) ListProjectDonationById(ctx context.Context, input model.ListProjectDonationInput) (output model.ListProjectDonationOutput, err error) {
	//TODO implement me
	panic("implement me")
}
