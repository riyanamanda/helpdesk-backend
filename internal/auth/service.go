package auth

import (
	"context"
	"errors"
	"time"

	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (LoginResponse, error)
}

type service struct {
	userRepo     user.UserRepository
	jwtSecret    string
	jwtExpiresIn time.Duration
}

func NewAuthService(repo user.UserRepository, jwtSecret string, jwtExpiresIn time.Duration) AuthService {
	return &service{
		userRepo:     repo,
		jwtSecret:    jwtSecret,
		jwtExpiresIn: jwtExpiresIn,
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

	token, err := utils.GenerateToken(currentUser.ID, currentUser.Role, s.jwtSecret, s.jwtExpiresIn)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken: token,
	}, nil
}
