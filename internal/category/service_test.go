package category_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	category "github.com/riyanamanda/helpdesk-backend/internal/category"
	categorymocks "github.com/riyanamanda/helpdesk-backend/internal/category/mocks"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
	testingutil "github.com/riyanamanda/helpdesk-backend/internal/shared/testing"
)

func TestService_RegisterCategory(t *testing.T) {
	testCases := []struct {
		name      string
		req       *category.CreateCategoryRequest
		setupMock func(*categorymocks.CategoryRepository)
		assertFn  func(*testing.T, category.CategoryResponse, error)
	}{
		{
			name: "success",
			req:  &category.CreateCategoryRequest{Name: "Laptop"},
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(data *category.Category) bool {
					return data != nil && data.Name == "Laptop"
				})).Return(nil).Once()
			},
			assertFn: func(t *testing.T, result category.CategoryResponse, err error) {
				require.NoError(t, err)
				assert.Equal(t, "Laptop", result.Name)
			},
		},
		{
			name: "already exists",
			req:  &category.CreateCategoryRequest{Name: "Laptop"},
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("Create", mock.Anything, mock.Anything).Return(category.ErrCategoryAlreadyExists).Once()
			},
			assertFn: func(t *testing.T, result category.CategoryResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, category.CategoryResponse{}, result)
				testingutil.AssertAppError(t, err, apperror.CODE_ALREADY_EXISTS, http.StatusConflict, "category already exists")
			},
		},
		{
			name: "repository error",
			req:  &category.CreateCategoryRequest{Name: "Laptop"},
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result category.CategoryResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, category.CategoryResponse{}, result)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := categorymocks.NewCategoryRepository(t)
			svc := category.NewCategoryService(repo)
			tc.setupMock(repo)

			result, err := svc.RegisterCategory(context.Background(), tc.req)
			tc.assertFn(t, result, err)
		})
	}
}

func TestService_FetchAllCategories(t *testing.T) {
	testCases := []struct {
		name      string
		params    *category.GetCategoryParams
		setupMock func(*categorymocks.CategoryRepository)
		assertFn  func(*testing.T, []category.CategoryResponse, int, error)
	}{
		{
			name:   "success",
			params: &category.GetCategoryParams{},
			setupMock: func(repo *categorymocks.CategoryRepository) {
				items := []category.Category{{ID: 1, Name: "Hardware"}, {ID: 2, Name: "Software"}}
				repo.On("GetAll", mock.Anything, category.GetCategoryParams{Params: pagination.Params{Page: 1, Limit: 10}}).Return(items, 2, nil).Once()
			},
			assertFn: func(t *testing.T, result []category.CategoryResponse, total int, err error) {
				require.NoError(t, err)
				assert.Len(t, result, 2)
				assert.Equal(t, "Hardware", result[0].Name)
				assert.Equal(t, 2, total)
			},
		},
		{
			name:   "repository error",
			params: &category.GetCategoryParams{},
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("GetAll", mock.Anything, mock.Anything).Return(nil, 0, errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result []category.CategoryResponse, total int, err error) {
				require.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, 0, total)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := categorymocks.NewCategoryRepository(t)
			svc := category.NewCategoryService(repo)
			tc.setupMock(repo)

			result, total, err := svc.FetchAllCategories(context.Background(), tc.params)
			tc.assertFn(t, result, total, err)
		})
	}
}

func TestService_FindCategoryByID(t *testing.T) {
	now := time.Now().UTC()

	testCases := []struct {
		name      string
		id        int64
		setupMock func(*categorymocks.CategoryRepository)
		assertFn  func(*testing.T, category.CategoryResponse, error)
	}{
		{
			name: "success",
			id:   10,
			setupMock: func(repo *categorymocks.CategoryRepository) {
				item := &category.Category{ID: 10, Name: "Hardware", IsActive: true, CreatedAt: now, UpdatedAt: now}
				repo.On("GetByID", mock.Anything, int64(10)).Return(item, nil).Once()
			},
			assertFn: func(t *testing.T, result category.CategoryResponse, err error) {
				require.NoError(t, err)
				assert.Equal(t, int64(10), result.ID)
				assert.Equal(t, "Hardware", result.Name)
			},
		},
		{
			name: "not found",
			id:   11,
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("GetByID", mock.Anything, int64(11)).Return(nil, category.ErrCategoryNotFound).Once()
			},
			assertFn: func(t *testing.T, result category.CategoryResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, category.CategoryResponse{}, result)
				testingutil.AssertAppError(t, err, apperror.CODE_NOT_FOUND, http.StatusNotFound, "category not found")
			},
		},
		{
			name: "repository error",
			id:   12,
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("GetByID", mock.Anything, int64(12)).Return(nil, errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result category.CategoryResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, category.CategoryResponse{}, result)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := categorymocks.NewCategoryRepository(t)
			svc := category.NewCategoryService(repo)
			tc.setupMock(repo)

			result, err := svc.FindCategoryByID(context.Background(), tc.id)
			tc.assertFn(t, result, err)
		})
	}
}

func TestService_EditCategory(t *testing.T) {
	testCases := []struct {
		name      string
		id        int64
		req       *category.UpdateCategoryRequest
		setupMock func(*categorymocks.CategoryRepository)
		assertFn  func(*testing.T, category.CategoryResponse, error)
	}{
		{
			name: "success",
			id:   20,
			req:  &category.UpdateCategoryRequest{Name: "Peripheral"},
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("Update", mock.Anything, int64(20), mock.MatchedBy(func(data *category.Category) bool {
					return data != nil && data.Name == "Peripheral"
				})).Return(nil).Once()
			},
			assertFn: func(t *testing.T, result category.CategoryResponse, err error) {
				require.NoError(t, err)
				assert.Equal(t, "Peripheral", result.Name)
			},
		},
		{
			name: "not found",
			id:   21,
			req:  &category.UpdateCategoryRequest{Name: "Peripheral"},
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("Update", mock.Anything, int64(21), mock.Anything).Return(category.ErrCategoryNotFound).Once()
			},
			assertFn: func(t *testing.T, result category.CategoryResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, category.CategoryResponse{}, result)
				testingutil.AssertAppError(t, err, apperror.CODE_NOT_FOUND, http.StatusNotFound, "category not found")
			},
		},
		{
			name: "already exists",
			id:   22,
			req:  &category.UpdateCategoryRequest{Name: "Peripheral"},
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("Update", mock.Anything, int64(22), mock.Anything).Return(category.ErrCategoryAlreadyExists).Once()
			},
			assertFn: func(t *testing.T, result category.CategoryResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, category.CategoryResponse{}, result)
				testingutil.AssertAppError(t, err, apperror.CODE_ALREADY_EXISTS, http.StatusConflict, "category already exists")
			},
		},
		{
			name: "repository error",
			id:   23,
			req:  &category.UpdateCategoryRequest{Name: "Peripheral"},
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("Update", mock.Anything, int64(23), mock.Anything).Return(errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result category.CategoryResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, category.CategoryResponse{}, result)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := categorymocks.NewCategoryRepository(t)
			svc := category.NewCategoryService(repo)
			tc.setupMock(repo)

			result, err := svc.EditCategory(context.Background(), tc.id, tc.req)
			tc.assertFn(t, result, err)
		})
	}
}

func TestService_DeleteCategory(t *testing.T) {
	testCases := []struct {
		name      string
		id        int64
		setupMock func(*categorymocks.CategoryRepository)
		assertFn  func(*testing.T, error)
	}{
		{
			name: "success",
			id:   30,
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("Delete", mock.Anything, int64(30)).Return(nil).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "not found",
			id:   31,
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("Delete", mock.Anything, int64(31)).Return(category.ErrCategoryNotFound).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				testingutil.AssertAppError(t, err, apperror.CODE_NOT_FOUND, http.StatusNotFound, "category not found")
			},
		},
		{
			name: "repository error",
			id:   32,
			setupMock: func(repo *categorymocks.CategoryRepository) {
				repo.On("Delete", mock.Anything, int64(32)).Return(errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := categorymocks.NewCategoryRepository(t)
			svc := category.NewCategoryService(repo)
			tc.setupMock(repo)

			err := svc.DeleteCategory(context.Background(), tc.id)
			tc.assertFn(t, err)
		})
	}
}
