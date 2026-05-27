package user

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FetchAllUsers(ctx context.Context, params *GetUserParams) ([]UserResponse, int64, error)
	RegisterUser(ctx context.Context, req *UserCreateRequest) error
	FindUserByID(ctx context.Context, id *uuid.UUID) (UserResponse, error)
}

type service struct {
	repo          UserRepository
	storageConfig config.Storage
}

func NewUserService(repo UserRepository, storageConfig config.Storage) UserService {
	return &service{
		repo:          repo,
		storageConfig: storageConfig,
	}
}

func (svc *service) FetchAllUsers(ctx context.Context, params *GetUserParams) ([]UserResponse, int64, error) {
	if params == nil {
		params = &GetUserParams{}
	}

	params.Normalize()

	users, total, err := svc.repo.GetAll(ctx, *params)
	if err != nil {
		return []UserResponse{}, 0, err
	}

	return toUserResponses(users, svc.storageConfig), total, nil
}

func (svc *service) RegisterUser(ctx context.Context, req *UserCreateRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var createdBy *uuid.UUID
	if currentUserID, ok := utils.GetUserIDFromContext(ctx); ok {
		createdBy = &currentUserID
	}

	normalizedEmail := strings.TrimSpace(strings.ToLower(req.Email))

	user := &User{
		Name:       req.Name,
		Email:      normalizedEmail,
		Password:   string(hashedPassword),
		Role:       req.Role,
		DivisionID: req.DivisionID,
		CreatedBy:  createdBy,
	}

	if err := svc.repo.Create(ctx, user); err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			return apperror.AlreadyExists("user")
		}
		return err
	}

	return nil
}

func (svc *service) FindUserByID(ctx context.Context, id *uuid.UUID) (UserResponse, error) {
	user, err := svc.repo.GetByID(ctx, *id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return UserResponse{}, apperror.NotFound("user")
		}
		return UserResponse{}, err
	}

	return toUserResponse(*user, svc.storageConfig), nil
}

