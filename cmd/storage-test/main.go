package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/riyanamanda/helpdesk-backend/internal/storage"
)

func main() {
	s, err := storage.NewMinioStorage()
	if err != nil {
		panic(err)
	}

	file, err := os.Open("avatar.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	err = s.Upload(context.Background(), "avatar/test.jpg", file, stat.Size(), "image/jpeg")
	if err != nil {
		panic(err)
	}
	slog.Info("upload success")
}
