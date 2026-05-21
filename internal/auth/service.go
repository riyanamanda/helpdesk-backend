package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (LoginResponse, error)
	Me(ctx context.Context) (CurrentUserResponse, error)
}

type service struct {
	userRepo      user.UserRepository
	config        config.Auth
	storageConfig config.Storage
}

func NewAuthService(repo user.UserRepository, cfg config.Auth, storageConfig config.Storage) AuthService {
	return &service{
		userRepo:      repo,
		config:        cfg,
		storageConfig: storageConfig,
	}
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (LoginResponse, error) {
	currentUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return LoginResponse{}, apperror.BadRequest("invalid email or password")
		}
		return LoginResponse{}, err
	}

	if !currentUser.IsActive {
		return LoginResponse{}, apperror.Forbidden("user is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(req.Password)); err != nil {
		return LoginResponse{}, apperror.BadRequest("invalid email or password")
	}

	token, err := utils.GenerateToken(currentUser.ID, string(currentUser.Role), s.config.JWTSecret, s.config.JWTExp)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		User:        toCurrentUserResponse(*currentUser, s.storageConfig),
		AccessToken: token,
	}, nil
}

func (s *service) Me(ctx context.Context) (CurrentUserResponse, error) {
	var userID uuid.UUID
	if currentUserID, ok := utils.GetUserIDFromContext(ctx); ok {
		userID = currentUserID
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return CurrentUserResponse{}, err
	}

	return toCurrentUserResponse(*user, s.storageConfig), nil
}
