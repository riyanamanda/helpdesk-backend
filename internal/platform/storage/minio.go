package storage

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type minioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(client *minio.Client, bucketName string) Storage {
	return &minioStorage{
		client:     client,
		bucketName: bucketName,
	}
}

func (m *minioStorage) Upload(ctx context.Context, key string, file *File) error {
	_, err := m.client.PutObject(ctx, m.bucketName, key, file.Content, file.Size, minio.PutObjectOptions{
		ContentType: file.ContentType,
	})

	return err
}

func (m *minioStorage) Delete(ctx context.Context, key string) error {
	return m.client.RemoveObject(ctx, m.bucketName, key, minio.RemoveObjectOptions{})
}
