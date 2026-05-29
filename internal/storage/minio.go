package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type minioStorage struct {
	client *minio.Client

	bucketName string

	publicURL string
}

func NewMinioStorage(client *minio.Client, bucketName, publicURL string) Storage {

	return &minioStorage{

		client: client,

		bucketName: bucketName,

		publicURL: publicURL,
	}

}

func (m *minioStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {

	_, err := m.client.PutObject(ctx, m.bucketName, key, reader, size, minio.PutObjectOptions{

		ContentType: contentType,
	})

	return err

}

func (m *minioStorage) Delete(ctx context.Context, key string) error {

	return m.client.RemoveObject(ctx, m.bucketName, key, minio.RemoveObjectOptions{})

}
