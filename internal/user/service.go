package user

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/event"
	"github.com/riyanamanda/helpdesk-backend/internal/outbox"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/jwtutil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/passwordtoken"
	"golang.org/x/crypto/bcrypt"
)

type repository interface {
	GetAll(ctx context.Context, params GetUserParams) ([]UserProjection, int64, error)
	Create(ctx context.Context, tx database.Tx, user *User) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*UserProjection, error)
	UpdateByID(ctx context.Context, id uuid.UUID, user User) error
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	AssignableUser(ctx context.Context) ([]AssignableUserProjection, error)
}

type presenceReader interface {
	Online(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]bool, error)
}

type Service struct {
	repo          repository
	outboxRepo    outbox.Repository
	txManager     *database.TxManager
	storageConfig config.Storage
	appConfig     config.App
	cache         cache.Cache
	presence      presenceReader
	publisher     event.Publisher
}

func NewService(repo repository, outboxRepo outbox.Repository, txManager *database.TxManager, storageConfig config.Storage, appConfig config.App, cache cache.Cache, presence presenceReader, publisher event.Publisher) *Service {
	return &Service{
		repo:          repo,
		outboxRepo:    outboxRepo,
		txManager:     txManager,
		storageConfig: storageConfig,
		appConfig:     appConfig,
		cache:         cache,
		presence:      presence,
		publisher:     publisher,
	}
}

func (s *Service) markOnline(ctx context.Context, users []UserResponse) {
	ids := make([]uuid.UUID, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}

	online, err := s.presence.Online(ctx, ids)
	if err != nil {
		slog.WarnContext(ctx, "read user presence", "error", err)
		return
	}

	for i := range users {
		users[i].IsOnline = online[users[i].ID]
	}
}

func (s *Service) ListUsers(ctx context.Context, params *GetUserParams) ([]UserResponse, int64, error) {
	if params == nil {
		params = &GetUserParams{}
	}
	params.Normalize()

	users, total, err := s.repo.GetAll(ctx, *params)
	if err != nil {
		return nil, 0, err
	}

	result := toUserResponses(users, s.storageConfig)
	s.markOnline(ctx, result)

	return result, total, nil
}

func (s *Service) CreateUser(ctx context.Context, req *UserCreateRequest) error {
	password := req.Password
	if password == "" {
		random, err := passwordtoken.New()
		if err != nil {
			return err
		}
		password = random
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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

	// start db transaction from platform/database/transaction.go
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	// defer rollback db transaction
	defer tx.Rollback()

	userID, err := s.repo.Create(ctx, tx, user)
	if err != nil {

		return err
	}

	setPasswordURL, err := passwordtoken.Issue(ctx, s.cache, s.appConfig.URL, userID, SetPasswordTokenTTL)
	if err != nil {
		return err
	}

	event := event.UserCreatedEvent{
		Name:           user.Name,
		Email:          user.Email,
		SetPasswordURL: setPasswordURL,
		ExpiresIn:      SetPasswordExpiryText,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	outboxEvent := outbox.OutboxEvent{
		EventType:   "user.created",
		AggregateID: userID.String(),
		Payload:     payload,
	}

	if err := s.outboxRepo.Create(ctx, tx, outboxEvent); err != nil {
		return err
	}

	// commit db transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	InvalidateCache(ctx, s.cache)

	return nil
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {

		return nil, err
	}

	result := []UserResponse{toUserResponse(*user, s.storageConfig)}
	s.markOnline(ctx, result)

	return &result[0], nil
}

func (s *Service) UpdateUser(ctx context.Context, userID uuid.UUID, req *UserUpdateRequest) error {
	existing, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user := User{
		Name:       req.Name,
		Email:      req.Email,
		RoleID:     req.Role,
		DivisionID: req.Division,
		Gender:     req.Gender,
		IsActive:   *req.IsActive,
	}

	if err := s.repo.UpdateByID(ctx, userID, user); err != nil {
		return err
	}

	InvalidateCache(ctx, s.cache)
	rbac.InvalidateUserPermissionsCache(ctx, s.cache, userID)
	rbac.InvalidateUserRoleCache(ctx, s.cache, userID)

	if existing.IsActive && !user.IsActive {
		return jwtutil.RevokeUserSessions(ctx, s.cache, userID)
	}

	return nil
}

func (s *Service) UpdatePassword(ctx context.Context, userID uuid.UUID, req *UserUpdatePasswordRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return err
	}

	if err := jwtutil.RevokeUserSessions(ctx, s.cache, userID); err != nil {
		return err
	}

	event.NotifyPasswordChanged(ctx, s.publisher, userID, true)

	return nil
}

func (s *Service) ListAssignableUser(ctx context.Context) ([]UserBrief, error) {
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
