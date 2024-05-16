package projectfilestorage

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"log"
	"os"
	"strings"
	"time"
)

type Storage struct {
	s3Client *s3.Client
}

func (s Storage) RequestUploadUrl(mimeType string, fileSize int64) (url string, expiredAt int64, err error) {
	presignClient := s3.NewPresignClient(s.s3Client)
	bucketName := os.Getenv("AWS_BUCKET")
	objectName := fmt.Sprintf("%d_%x.jpg", time.Now().Unix(), makeRandomBytes(8))
	duration := 15 * time.Minute
	expiredAt = time.Now().Add(duration).Unix()
	presignedUrl, err := presignClient.PresignPutObject(context.Background(),
		&s3.PutObjectInput{
			Bucket:        aws.String(bucketName),
			Key:           aws.String(objectName),
			ACL:           types.ObjectCannedACLPublicRead,
			ContentType:   aws.String(mimeType),
			ContentLength: aws.Int64(fileSize),
		},
		s3.WithPresignExpires(duration),
	)

	if err != nil {
		log.Fatalf("failed to create presigned URL, %v", err)
	}

	url = strings.Replace(presignedUrl.URL, "localstack", "localhost", -1)

	return
}

func NewStorage(s3Client *s3.Client) *Storage {
	return &Storage{s3Client: s3Client}
}

func makeRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
