package category

import (
	"context"
	"errors"

	apperrors "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
)

type CategoryService interface {
	FetchAllCategories(ctx context.Context, params *GetCategoryParams) ([]CategoryResponse, int64, error)
	RegisterCategory(ctx context.Context, req *CreateCategoryRequest) (CategoryResponse, error)
	FindCategoryByID(ctx context.Context, id int64) (CategoryResponse, error)
	EditCategory(ctx context.Context, id int64, req *UpdateCategoryRequest) error
	DeleteCategory(ctx context.Context, id int64) error
}

type service struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &service{
		repo: repo,
	}
}

func (svc *service) FetchAllCategories(ctx context.Context, params *GetCategoryParams) ([]CategoryResponse, int64, error) {
	if params == nil {
		params = &GetCategoryParams{}
	}

	params.Normalize()

	categories, total, err := svc.repo.GetAll(ctx, *params)
	if err != nil {
		return []CategoryResponse{}, 0, err
	}

	return toCategoryResponses(categories), total, nil
}

func (svc *service) RegisterCategory(ctx context.Context, req *CreateCategoryRequest) (CategoryResponse, error) {
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

func (svc *service) FindCategoryByID(ctx context.Context, id int64) (CategoryResponse, error) {
	category, err := svc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return CategoryResponse{}, apperrors.NotFound("category")
		}
		return CategoryResponse{}, err
	}

	return toCategoryResponse(*category), nil
}

func (svc *service) EditCategory(ctx context.Context, id int64, req *UpdateCategoryRequest) error {
	existing, err := svc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return apperrors.NotFound("category")
		}
		return err
	}

	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	category := Category{
		Name:     req.Name,
		IsActive: isActive,
	}
	if err := svc.repo.Update(ctx, id, &category); err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return apperrors.NotFound("category")
		}
		if errors.Is(err, ErrCategoryAlreadyExists) {
			return apperrors.AlreadyExists("category")
		}
		return err
	}

	return nil
}

func (svc *service) DeleteCategory(ctx context.Context, id int64) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return apperrors.NotFound("category")
		}
		return err
	}

	return nil
}
