package user

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	apperrors "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUser(ctx context.Context, params *GetUserParams) ([]UserResponse, int, error)
	Create(ctx context.Context, req *UserCreateRequest) (UserResponse, error)
	GetById(ctx context.Context, id *uuid.UUID) (UserResponse, error)
}

type service struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetUser(ctx context.Context, params *GetUserParams) ([]UserResponse, int, error) {
	if params == nil {
		params = &GetUserParams{}
	}

	page, limit, _ := params.Normalize()
	params.Page = page
	params.Limit = limit

	users, total, err := svc.repo.List(ctx, *params)
	if err != nil {
		slog.Error("List user failed", "error", err)
		return []UserResponse{}, 0, err
	}

	return toUserResponses(users), total, nil
}

func (svc *service) Create(ctx context.Context, req *UserCreateRequest) (UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, err
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
			return UserResponse{}, apperrors.AlreadyExists("user")
		}
		return UserResponse{}, err
	}

	return toUserResponse(*user), nil
}

func (svc *service) GetById(ctx context.Context, id *uuid.UUID) (UserResponse, error) {
	user, err := svc.repo.GetByID(ctx, *id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return UserResponse{}, apperror.NotFound("user")
		}
		return UserResponse{}, err
	}

	return toUserResponse(*user), nil
}
