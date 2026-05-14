package division_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	division "github.com/riyanamanda/helpdesk-backend/internal/division"
	divisionmocks "github.com/riyanamanda/helpdesk-backend/internal/division/mocks"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
	testingutil "github.com/riyanamanda/helpdesk-backend/internal/shared/testing"
)

func TestService_RegisterDivision(t *testing.T) {
	testCases := []struct {
		name      string
		req       *division.CreateDivisionRequest
		setupMock func(*divisionmocks.DivisionRepository)
		assertFn  func(*testing.T, division.DivisionResponse, error)
	}{
		{
			name: "success",
			req:  &division.CreateDivisionRequest{Name: "IT Support"},
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(data *division.Division) bool {
					return data != nil && data.Name == "IT Support"
				})).Return(nil).Once()
			},
			assertFn: func(t *testing.T, result division.DivisionResponse, err error) {
				require.NoError(t, err)
				assert.Equal(t, "IT Support", result.Name)
			},
		},
		{
			name: "already exists",
			req:  &division.CreateDivisionRequest{Name: "IT Support"},
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				repo.On("Create", mock.Anything, mock.Anything).Return(division.ErrDivisionAlreadyExists).Once()
			},
			assertFn: func(t *testing.T, result division.DivisionResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, division.DivisionResponse{}, result)
				testingutil.AssertAppError(t, err, apperror.CODE_ALREADY_EXISTS, http.StatusConflict, "division already exists")
			},
		},
		{
			name: "repository error",
			req:  &division.CreateDivisionRequest{Name: "IT Support"},
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result division.DivisionResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, division.DivisionResponse{}, result)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := divisionmocks.NewDivisionRepository(t)
			svc := division.NewDivisionService(repo)
			tc.setupMock(repo)

			result, err := svc.RegisterDivision(context.Background(), tc.req)
			tc.assertFn(t, result, err)
		})
	}
}

func TestService_FetchAllDivisions(t *testing.T) {
	testCases := []struct {
		name      string
		params    *division.GetDivisionParams
		setupMock func(*divisionmocks.DivisionRepository)
		assertFn  func(*testing.T, []division.DivisionResponse, int, error)
	}{
		{
			name:   "success",
			params: &division.GetDivisionParams{},
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				items := []division.Division{{ID: 1, Name: "IT"}, {ID: 2, Name: "HR"}}
				repo.On("GetAll", mock.Anything, division.GetDivisionParams{Params: pagination.Params{Page: 1, Limit: 10}}).Return(items, 2, nil).Once()
			},
			assertFn: func(t *testing.T, result []division.DivisionResponse, total int, err error) {
				require.NoError(t, err)
				assert.Len(t, result, 2)
				assert.Equal(t, "IT", result[0].Name)
				assert.Equal(t, 2, total)
			},
		},
		{
			name:   "repository error",
			params: &division.GetDivisionParams{},
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				repo.On("GetAll", mock.Anything, mock.Anything).Return(nil, 0, errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result []division.DivisionResponse, total int, err error) {
				require.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, 0, total)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := divisionmocks.NewDivisionRepository(t)
			svc := division.NewDivisionService(repo)
			tc.setupMock(repo)

			result, total, err := svc.FetchAllDivisions(context.Background(), tc.params)
			tc.assertFn(t, result, total, err)
		})
	}
}

func TestService_FindDivisionByID(t *testing.T) {
	now := time.Now().UTC()

	testCases := []struct {
		name      string
		id        int64
		setupMock func(*divisionmocks.DivisionRepository)
		assertFn  func(*testing.T, division.DivisionResponse, error)
	}{
		{
			name: "success",
			id:   10,
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				item := &division.Division{ID: 10, Name: "IT", IsActive: true, CreatedAt: now, UpdatedAt: now}
				repo.On("GetByID", mock.Anything, int64(10)).Return(item, nil).Once()
			},
			assertFn: func(t *testing.T, result division.DivisionResponse, err error) {
				require.NoError(t, err)
				assert.Equal(t, int64(10), result.ID)
				assert.Equal(t, "IT", result.Name)
			},
		},
		{
			name: "not found",
			id:   11,
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				repo.On("GetByID", mock.Anything, int64(11)).Return(nil, division.ErrDivisionNotFound).Once()
			},
			assertFn: func(t *testing.T, result division.DivisionResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, division.DivisionResponse{}, result)
				testingutil.AssertAppError(t, err, apperror.CODE_NOT_FOUND, http.StatusNotFound, "division not found")
			},
		},
		{
			name: "repository error",
			id:   12,
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				repo.On("GetByID", mock.Anything, int64(12)).Return(nil, errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result division.DivisionResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, division.DivisionResponse{}, result)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := divisionmocks.NewDivisionRepository(t)
			svc := division.NewDivisionService(repo)
			tc.setupMock(repo)

			result, err := svc.FindDivisionByID(context.Background(), tc.id)
			tc.assertFn(t, result, err)
		})
	}
}

func TestService_EditDivision(t *testing.T) {
	testCases := []struct {
		name      string
		id        int64
		req       *division.UpdateDivisionRequest
		setupMock func(*divisionmocks.DivisionRepository)
		assertFn  func(*testing.T, division.DivisionResponse, error)
	}{
		{
			name: "success",
			id:   20,
			req:  &division.UpdateDivisionRequest{Name: "Finance"},
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				repo.On("Update", mock.Anything, int64(20), mock.MatchedBy(func(data *division.Division) bool {
					return data != nil && data.Name == "Finance"
				})).Return(nil).Once()
			},
			assertFn: func(t *testing.T, result division.DivisionResponse, err error) {
				require.NoError(t, err)
				assert.Equal(t, "Finance", result.Name)
			},
		},
		{
			name: "not found",
			id:   21,
			req:  &division.UpdateDivisionRequest{Name: "Finance"},
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				repo.On("Update", mock.Anything, int64(21), mock.Anything).Return(division.ErrDivisionNotFound).Once()
			},
			assertFn: func(t *testing.T, result division.DivisionResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, division.DivisionResponse{}, result)
				testingutil.AssertAppError(t, err, apperror.CODE_NOT_FOUND, http.StatusNotFound, "division not found")
			},
		},
		{
			name: "already exists",
			id:   22,
			req:  &division.UpdateDivisionRequest{Name: "Finance"},
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				repo.On("Update", mock.Anything, int64(22), mock.Anything).Return(division.ErrDivisionAlreadyExists).Once()
			},
			assertFn: func(t *testing.T, result division.DivisionResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, division.DivisionResponse{}, result)
				testingutil.AssertAppError(t, err, apperror.CODE_ALREADY_EXISTS, http.StatusConflict, "division already exists")
			},
		},
		{
			name: "repository error",
			id:   23,
			req:  &division.UpdateDivisionRequest{Name: "Finance"},
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				repo.On("Update", mock.Anything, int64(23), mock.Anything).Return(errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result division.DivisionResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, division.DivisionResponse{}, result)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := divisionmocks.NewDivisionRepository(t)
			svc := division.NewDivisionService(repo)
			tc.setupMock(repo)

			result, err := svc.EditDivision(context.Background(), tc.id, tc.req)
			tc.assertFn(t, result, err)
		})
	}
}

func TestService_DeleteDivision(t *testing.T) {
	testCases := []struct {
		name      string
		id        int64
		setupMock func(*divisionmocks.DivisionRepository)
		assertFn  func(*testing.T, error)
	}{
		{
			name: "success",
			id:   30,
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				repo.On("Delete", mock.Anything, int64(30)).Return(nil).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "not found",
			id:   31,
			setupMock: func(repo *divisionmocks.DivisionRepository) {
				repo.On("Delete", mock.Anything, int64(31)).Return(division.ErrDivisionNotFound).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				testingutil.AssertAppError(t, err, apperror.CODE_NOT_FOUND, http.StatusNotFound, "division not found")
			},
		},
		{
			name: "repository error",
			id:   32,
			setupMock: func(repo *divisionmocks.DivisionRepository) {
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
			repo := divisionmocks.NewDivisionRepository(t)
			svc := division.NewDivisionService(repo)
			tc.setupMock(repo)

			err := svc.DeleteDivision(context.Background(), tc.id)
			tc.assertFn(t, err)
		})
	}
}
