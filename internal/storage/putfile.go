package storage

import (
	"collegeWaleServer/errz"
	"context"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
)

func PutFile(path string, src multipart.File, contentType string) error {
	client := InitR2Client()
	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("STORAGE_BUCKET")),
		Key:         aws.String(path),
		Body:        src,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Errorf("Failed to upload file: %v", err)
		return errz.NewBadRequest("upload file failed.")
	}
	return nil
}
