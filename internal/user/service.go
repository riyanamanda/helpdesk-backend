package user

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	ListUsers(ctx context.Context, params *GetUserParams) ([]UserResponse, int64, error)
	CreateUser(ctx context.Context, req *UserCreateRequest) error
	GetUser(ctx context.Context, id *uuid.UUID) (*UserResponse, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req *UserUpdateRequest) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, req *UserUpdatePasswordRequest) error
	ListAssignableUser(ctx context.Context) ([]UserBrief, error)
}

type service struct {
	repo          UserRepository
	storageConfig config.Storage
	cache         cache.Cache
}

func NewUserService(repo UserRepository, storageConfig config.Storage, cache cache.Cache) UserService {
	return &service{
		repo:          repo,
		storageConfig: storageConfig,
		cache:         cache,
	}
}

func (s *service) ListUsers(ctx context.Context, params *GetUserParams) ([]UserResponse, int64, error) {
	if params == nil {
		params = &GetUserParams{}
	}
	params.Normalize()

	users, total, err := s.repo.GetAll(ctx, *params)
	if err != nil {
		return nil, 0, err
	}

	return toUserResponses(users, s.storageConfig), total, nil
}

func (s *service) CreateUser(ctx context.Context, req *UserCreateRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var createdBy *uuid.UUID
	if currentUserID, ok := ctxkey.GetUserIDFromContext(ctx); ok {
		createdBy = &currentUserID
	}

	normalizedEmail := strings.TrimSpace(strings.ToLower(req.Email))
	user := &User{
		Name:       req.Name,
		Email:      normalizedEmail,
		Password:   string(hashedPassword),
		RoleID:     req.Role,
		DivisionID: req.Division,
		Gender:     strings.ToUpper(req.Gender),
		CreatedBy:  createdBy,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			return apperr.AlreadyExists("user")
		}

		return err
	}

	InvalidateCache(ctx, s.cache)

	return nil
}

func (s *service) GetUser(ctx context.Context, id *uuid.UUID) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, *id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, apperr.NotFound("user")
		}

		return nil, err
	}

	result := toUserResponse(*user, s.storageConfig)

	return &result, nil
}

func (s *service) UpdateUser(ctx context.Context, userID uuid.UUID, req *UserUpdateRequest) error {
	user := User{
		Name:       req.Name,
		Email:      req.Email,
		RoleID:     req.Role,
		DivisionID: req.Division,
		Gender:     req.Gender,
		IsActive:   *req.IsActive,
	}

	if err := s.repo.UpdateByID(ctx, userID, user); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return apperr.NotFound("user")
		}
		if errors.Is(err, ErrUserAlreadyExists) {
			return apperr.AlreadyExists("user")
		}
		return err
	}

	InvalidateCache(ctx, s.cache)

	return nil
}

func (s *service) UpdatePassword(ctx context.Context, userID uuid.UUID, req *UserUpdatePasswordRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return apperr.NotFound("user")
		}
		return err
	}

	return nil
}

func (s *service) ListAssignableUser(ctx context.Context) ([]UserBrief, error) {
	cached, err := s.cache.Get(ctx, AssignableCacheKey)
	if err == nil {
		var users []UserBrief

		if err := json.Unmarshal([]byte(cached), &users); err == nil {
			return users, nil
		}
	}

	projection, err := s.repo.AssignableUser(ctx)
	if err != nil {
		return nil, err
	}

	users := toUserBriefs(projection)

	if len(users) > 0 {
		data, err := json.Marshal(users)
		if err == nil {
			_ = s.cache.Set(ctx, AssignableCacheKey, string(data), 24*time.Hour)
		}
	}

	return users, nil
}
