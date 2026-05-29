package user

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperror"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
)

type UserService interface {
	ListUsers(ctx context.Context, params *GetUserParams) ([]UserResponse, int64, error)

	CreateUser(ctx context.Context, req *UserCreateRequest) error

	GetUser(ctx context.Context, id *uuid.UUID) (UserResponse, error)
}

type service struct {
	repo UserRepository

	storageConfig config.Storage
}

func NewUserService(repo UserRepository, storageConfig config.Storage) UserService {

	return &service{

		repo: repo,

		storageConfig: storageConfig,
	}

}

func (s *service) ListUsers(ctx context.Context, params *GetUserParams) ([]UserResponse, int64, error) {

	if params == nil {

		params = &GetUserParams{}

	}

	params.Normalize()

	users, total, err := s.repo.GetAll(ctx, *params)

	if err != nil {

		return []UserResponse{}, 0, err

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

		Name: req.Name,

		Email: normalizedEmail,

		Password: string(hashedPassword),

		Role: req.Role,

		DivisionID: req.DivisionID,

		CreatedBy: createdBy,
	}

	if err := s.repo.Create(ctx, user); err != nil {

		if errors.Is(err, ErrUserAlreadyExists) {

			return apperror.AlreadyExists("user")

		}

		return err

	}

	return nil

}

func (s *service) GetUser(ctx context.Context, id *uuid.UUID) (UserResponse, error) {

	user, err := s.repo.GetByID(ctx, *id)

	if err != nil {

		if errors.Is(err, ErrUserNotFound) {

			return UserResponse{}, apperror.NotFound("user")

		}

		return UserResponse{}, err

	}

	return toUserResponse(*user, s.storageConfig), nil

}
