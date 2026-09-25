package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type rustfsStorage struct {
	client     *s3.Client
	bucketName string
}

func NewRustFSStorage(client *s3.Client, bucketName string) Storage {
	return &rustfsStorage{
		client:     client,
		bucketName: bucketName,
	}
}

func (r *rustfsStorage) Upload(ctx context.Context, key string, file *File) error {
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(key),
		Body:        file.Content,
		ContentType: aws.String(file.ContentType),
	})

	return err
}

func (r *rustfsStorage) Delete(ctx context.Context, key string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})

	return err
}
