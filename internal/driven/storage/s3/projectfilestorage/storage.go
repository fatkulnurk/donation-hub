package projectfilestorage

import "github.com/aws/aws-sdk-go-v2/service/s3"

type Storage struct {
	s3Client *s3.Client
}

func (s Storage) RequestUploadUrl(mimeType string, fileSize int) (url string, err error) {
	//TODO implement me
	panic("implement me")
}

func NewStorage(s3Client *s3.Client) *Storage {
	return &Storage{s3Client: s3Client}
}
