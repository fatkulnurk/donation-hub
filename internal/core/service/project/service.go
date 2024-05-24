package project

import (
	"context"
	"errors"
	"fmt"
	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/isdzulqor/donation-hub/internal/core/service/user"
	_type "github.com/isdzulqor/donation-hub/internal/core/type"
)

type Storage struct {
	storage         DataStorage
	fileStorage     FileStorage
	userDataStorage user.DataStorage
}

type Service interface {
	RequestUploadUrl(ctx context.Context, input model.RequestUploadUrlInput) (output model.RequestUploadUrlOutput, err error)
	SubmitProject(ctx context.Context, input model.SubmitProjectInput) (output model.SubmitProjectOutput, err error)
	ReviewProjectByAdmin(ctx context.Context, input model.ReviewProjectByAdminInput) (ok bool, err error)
	ListProject(ctx context.Context, input model.ListProjectInput) (output model.ListProjectOutput, err error)
	GetProjectById(ctx context.Context, input model.GetProjectByIdInput) (output model.GetProjectByIdOutput, err error)
	DonateToProject(ctx context.Context, input model.DonateToProjectInput) (ok bool, err error)
	ListDonationByProjectId(ctx context.Context, input model.ListProjectDonationInput) (output model.ListProjectDonationOutput, err error)
}

func NewService(storage DataStorage, fileStorage FileStorage, userDataStorage user.DataStorage) Service {
	return &Storage{
		storage:         storage,
		fileStorage:     fileStorage,
		userDataStorage: userDataStorage,
	}
}

func (s *Storage) RequestUploadUrl(ctx context.Context, input model.RequestUploadUrlInput) (output model.RequestUploadUrlOutput, err error) {
	// validate user, make sure role is valid
	ok, err := s.userDataStorage.UserHasRole(ctx, input.UserID, _type.ROLE_REQUESTER)
	if !ok {
		err = errors.New("ERR_FORBIDDEN_ACCESS")
		return
	}

	// validate size
	if input.FileSize > 1048576 {
		err = errors.New("filesize can't greater than 1MB")
		return
	}

	if input.FileSize <= 0 {
		err = errors.New("filesize must greater than 0Kb")
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

func (s *Storage) SubmitProject(ctx context.Context, input model.SubmitProjectInput) (output model.SubmitProjectOutput, err error) {
	// validate user, make sure role is valid
	ok, err := s.userDataStorage.UserHasRole(ctx, input.UserID, _type.ROLE_REQUESTER)
	if !ok {
		err = errors.New("ERR_FORBIDDEN_ACCESS")
		return
	}

	// save to database
	projectId, err := s.storage.Submit(ctx, input)

	if err != nil {
		return
	}

	// mapping to output
	output.ID = projectId
	output.Title = input.Title
	output.Description = input.Description
	output.ImageURLs = input.ImageURLs
	output.Currency = input.Currency
	output.DueAt = input.DueAt
	output.TargetAmount = input.TargetAmount

	return
}

func (s *Storage) ReviewProjectByAdmin(ctx context.Context, input model.ReviewProjectByAdminInput) (ok bool, err error) {
	// validate user, make sure role is valid
	ok, err = s.userDataStorage.UserHasRole(ctx, input.UserID, _type.ROLE_REQUESTER)
	if !ok {
		err = errors.New("ERR_FORBIDDEN_ACCESS")
		return
	}

	if input.Status != _type.PROJECT_APPROVED && input.Status != _type.PROJECT_REJECTED {
		ok = false
		err = errors.New("status must be approved or rejected")
		return
	}

	err = s.storage.ReviewByAdmin(ctx, input)

	return
}

func (s *Storage) ListProject(ctx context.Context, input model.ListProjectInput) (output model.ListProjectOutput, err error) {
	// make sure user has role admin if status need_review
	if input.Status == _type.PROJECT_NEED_REVIEW {
		ok, _ := s.userDataStorage.UserHasRole(ctx, input.UserID, _type.ROLE_ADMIN)
		if !ok {
			err = errors.New("ERR_FORBIDDEN_ACCESS")
			return
		}
	}

	output, err = s.storage.ListProject(ctx, input)

	return
}

func (s *Storage) GetProjectById(ctx context.Context, input model.GetProjectByIdInput) (output model.GetProjectByIdOutput, err error) {
	output, err = s.storage.GetProjectById(ctx, input)

	return
}

func (s *Storage) DonateToProject(ctx context.Context, input model.DonateToProjectInput) (ok bool, err error) {
	// make sure user has role donor
	ok, _ = s.userDataStorage.UserHasRole(ctx, input.UserID, _type.ROLE_DONOR)
	if !ok {
		fmt.Println("Masuk sini 3")
		err = errors.New("ERR_FORBIDDEN_ACCESS")
		return
	}

	ok = false
	p, err := s.storage.GetProjectById(ctx, model.GetProjectByIdInput{ProjectId: input.ProjectId})

	fmt.Println("Masuk sini 2")
	if err != nil {
		return
	}

	fmt.Println(p)
	fmt.Println(input)
	fmt.Println("Masuk sini 1")
	if float64(input.Amount) > p.TargetAmount || float64(input.Amount) > p.CollectionAmount {
		err = errors.New("ERR_TOO_MUCH_DONATION")
		return
	}

	fmt.Println("Masuk sini")
	err = s.storage.DonateToProject(ctx, input)
	if err == nil {
		ok = true
	}

	return
}

func (s *Storage) ListDonationByProjectId(ctx context.Context, input model.ListProjectDonationInput) (output model.ListProjectDonationOutput, err error) {
	output, err = s.storage.ListDonationByProjectId(ctx, input)

	return
}
