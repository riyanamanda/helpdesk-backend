package category

import (
	"context"
	"errors"
	"log/slog"

	apperrors "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
)

type CategoryService interface {
	GetCategories(ctx context.Context, params *GetCategoryParams) ([]CategoryResponse, int, error)
	Create(ctx context.Context, req *CreateCategoryRequest) (CategoryResponse, error)
}

type service struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetCategories(ctx context.Context, params *GetCategoryParams) ([]CategoryResponse, int, error) {
	categories, total, err := svc.repo.List(ctx, *params)
	if err != nil {
		slog.Error("list category failed", "error", err)
		return []CategoryResponse{}, 0, err
	}

	return toCategoryResponses(categories), total, nil
}

func (svc *service) Create(ctx context.Context, req *CreateCategoryRequest) (CategoryResponse, error) {
	category := Category{
		Name: req.Name,
	}

	if err := svc.repo.Create(ctx, &category); err != nil {
		if errors.Is(err, ErrCategoryAlreadyExists) {
			return CategoryResponse{}, apperrors.AlreadyExists("category")
		}
		return CategoryResponse{}, err
	}

	return toCategoryResponse(category), nil
}
