package category

import (
	"context"
	"errors"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperror"
)

type CategoryService interface {
	ListCategories(ctx context.Context, params *GetCategoryParams) ([]CategoryResponse, int64, error)
	CreateCategory(ctx context.Context, req *CreateCategoryRequest) (CategoryResponse, error)
	GetCategory(ctx context.Context, id int64) (CategoryResponse, error)
	UpdateCategory(ctx context.Context, id int64, req *UpdateCategoryRequest) error
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

func (s *service) ListCategories(ctx context.Context, params *GetCategoryParams) ([]CategoryResponse, int64, error) {
	if params == nil {
		params = &GetCategoryParams{}
	}

	params.Normalize()

	categories, total, err := s.repo.GetAll(ctx, *params)
	if err != nil {
		return []CategoryResponse{}, 0, err
	}

	return toCategoryResponses(categories), total, nil
}

func (s *service) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (CategoryResponse, error) {
	category := Category{
		Name: req.Name,
	}

	if err := s.repo.Create(ctx, &category); err != nil {
		if errors.Is(err, ErrCategoryAlreadyExists) {
			return CategoryResponse{}, apperror.AlreadyExists("category")
		}
		return CategoryResponse{}, err
	}

	return toCategoryResponse(category), nil
}

func (s *service) GetCategory(ctx context.Context, id int64) (CategoryResponse, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return CategoryResponse{}, apperror.NotFound("category")
		}
		return CategoryResponse{}, err
	}

	return toCategoryResponse(*category), nil
}

func (s *service) UpdateCategory(ctx context.Context, id int64, req *UpdateCategoryRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return apperror.NotFound("category")
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

	if err := s.repo.Update(ctx, id, &category); err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return apperror.NotFound("category")
		}
		if errors.Is(err, ErrCategoryAlreadyExists) {
			return apperror.AlreadyExists("category")
		}
		return err
	}

	return nil
}
func (s *service) DeleteCategory(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return apperror.NotFound("category")
		}
		return err
	}

	return nil
}
