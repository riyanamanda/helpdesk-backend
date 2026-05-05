package category

import (
	"context"
	"log/slog"
)

type CategoryService interface {
	GetCategories(ctx context.Context, params ListCategoriesParams) ([]Category, error)
}

type service struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetCategories(ctx context.Context, params ListCategoriesParams) ([]Category, error) {
	const (
		defaultLimit = 10
		maxLimit     = 100
	)

	limit := params.Limit
	offset := params.Offset

	if limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}
	if offset < 0 {
		offset = 0
	}

	params.Limit = limit
	params.Offset = offset

	categories, err := svc.repo.List(ctx, params)
	if err != nil {
		slog.Error("list category failed", "error", err)
		return nil, err
	}

	return categories, nil
}
