package category

import (
	"context"
	"log/slog"

	"github.com/riyanamanda/helpdesk-backend/internal/response"
)

type CategoryService interface {
	GetCategories(ctx context.Context, params ListCategoriesParams) (ListCategoriesResult, error)
	Create(ctx context.Context, req *CreateCategoryRequest) (Category, error)
}

type service struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetCategories(ctx context.Context, params ListCategoriesParams) (ListCategoriesResult, error) {
	limit := params.Limit
	offset := params.Offset

	if limit <= 0 || limit > response.MaxLimit {
		limit = response.DefaultLimit
	}
	if offset < 0 {
		offset = 0
	}

	params.Limit = limit
	params.Offset = offset

	categories, total, err := svc.repo.List(ctx, params)
	if err != nil {
		slog.Error("list category failed", "error", err)
		return ListCategoriesResult{}, err
	}

	return ListCategoriesResult{
		Data:  categories,
		Total: total,
	}, nil
}

func (svc *service) Create(ctx context.Context, req *CreateCategoryRequest) (Category, error) {
	category := Category{
		Name: req.Name,
	}

	if err := svc.repo.Create(ctx, &category); err != nil {
		return Category{}, err
	}

	return category, nil
}
