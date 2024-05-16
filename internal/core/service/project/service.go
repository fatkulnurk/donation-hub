package project

import (
	"context"
	"github.com/isdzulqor/donation-hub/internal/driver/rest/req"
)

type Storage struct {
	storage     DataStorage
	fileStorage FileStorage
}

type Service interface {
	RequestUploadUrl(mimeType string, fileSize int) (url string, expiredAt int, err error)
	Submit(rb req.SubmitProjectReqBody) (err error)
	ReviewByAdmin(ctx context.Context, rb req.ReviewProjectReqBody) (err error)
	Get(ctx context.Context) (err error)
	GetById(ctx context.Context, id uint32) (err error)
	DonateToProject(ctx context.Context, id uint32, rb req.DonateToProjectReqBody) (err error)
	GetDonationById(ctx context.Context, id uint32) (err error)
}

func NewService(storage DataStorage, fileStorage FileStorage) Service {
	return &Storage{
		storage:     storage,
		fileStorage: fileStorage,
	}
}

func (s Storage) RequestUploadUrl(mimeType string, fileSize int) (url string, expiredAt int, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) Submit(rb req.SubmitProjectReqBody) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) ReviewByAdmin(ctx context.Context, rb req.ReviewProjectReqBody) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) Get(ctx context.Context) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetById(ctx context.Context, id uint32) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) DonateToProject(ctx context.Context, id uint32, rb req.DonateToProjectReqBody) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetDonationById(ctx context.Context, id uint32) (err error) {
	//TODO implement me
	panic("implement me")
}
