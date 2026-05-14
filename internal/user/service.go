package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FetchAllUsers(ctx context.Context, params *GetUserParams) ([]UserResponse, int, error)
	RegisterUser(ctx context.Context, req *UserCreateRequest) (UserResponse, error)
	FindUserByID(ctx context.Context, id *uuid.UUID) (UserResponse, error)
	UpdateUserAvatar(ctx context.Context, file multipart.File, header *multipart.FileHeader) error
}

type service struct {
	repo    UserRepository
	storage storage.Storage
}

func NewUserService(repo UserRepository, storage storage.Storage) UserService {
	return &service{
		repo:    repo,
		storage: storage,
	}
}

func (svc *service) FetchAllUsers(ctx context.Context, params *GetUserParams) ([]UserResponse, int, error) {
	if params == nil {
		params = &GetUserParams{}
	}

	page, limit, _ := params.Normalize()
	params.Page = page
	params.Limit = limit

	users, total, err := svc.repo.GetAll(ctx, *params)
	if err != nil {
		slog.Error("List user failed", "error", err)
		return []UserResponse{}, 0, err
	}

	return toUserResponses(users, svc.storage), total, nil
}

func (svc *service) RegisterUser(ctx context.Context, req *UserCreateRequest) (UserResponse, error) {
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
			return UserResponse{}, apperror.AlreadyExists("user")
		}
		return UserResponse{}, err
	}

	return toUserResponse(*user, svc.storage), nil
}

func (svc *service) FindUserByID(ctx context.Context, id *uuid.UUID) (UserResponse, error) {
	user, err := svc.repo.GetByID(ctx, *id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return UserResponse{}, apperror.NotFound("user")
		}
		return UserResponse{}, err
	}

	return toUserResponse(*user, svc.storage), nil
}

func (svc *service) UpdateUserAvatar(ctx context.Context, file multipart.File, header *multipart.FileHeader) error {
	userID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		return apperror.Forbidden("unauthorized")
	}

	contentType := header.Header.Get("Content-Type")
	objectKey := fmt.Sprintf("avatar/%s/%d-%s", userID.String(), time.Now().Unix(), header.Filename)
	if err := svc.storage.Upload(ctx, objectKey, file, header.Size, contentType); err != nil {
		return err
	}

	if err := svc.repo.UpdateAvatar(ctx, userID, objectKey); err != nil {
		return err
	}

	return nil
}
