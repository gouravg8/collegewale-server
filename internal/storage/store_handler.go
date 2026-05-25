package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetUploadURL(client *s3.Client, bucket, filename string) (string, error) {
	presigner := s3.NewPresignClient(client)
	req, err := presigner.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get sign upload request: %w", err)
	}
	return req.URL, err
}

func GetViewURL(client *s3.Client, bucket, filename string) (string, error) {
	presigner := s3.NewPresignClient(client)
	req, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", fmt.Errorf("failed to sign request: %w", err)
	}
	return req.URL, err
}
