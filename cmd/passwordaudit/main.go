package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/riyanamanda/helpdesk-backend/internal/event"
	"github.com/riyanamanda/helpdesk-backend/internal/outbox"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/redis"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/jwtutil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/passwordtoken"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

const (
	knownDefaultPassword = "password"
	resetLinkTTL         = 72 * time.Hour
	resetLinkExpiryText  = "72 jam"
)

type account struct {
	ID          uuid.UUID  `db:"id"`
	Name        string     `db:"name"`
	Email       string     `db:"email"`
	Password    string     `db:"password"`
	IsActive    bool       `db:"is_active"`
	LastLoginAt *time.Time `db:"last_login_at"`
}

func main() {
	if err := run(); err != nil {
		slog.Error("password audit stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	apply := flag.Bool("apply", false, "replace the passwords of affected accounts, sign them out and email a set-password link")
	flag.Parse()

	ctx := context.Background()
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.Database.ConnString())
	if err != nil {
		return err
	}
	defer db.Close()

	var accounts []account
	if err := db.SelectContext(ctx, &accounts, `SELECT id, name, email, password, is_active, last_login_at FROM users ORDER BY name`); err != nil {
		return fmt.Errorf("list users: %w", err)
	}

	var affected []account
	for _, a := range accounts {
		if bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(knownDefaultPassword)) == nil {
			affected = append(affected, a)
		}
	}

	fmt.Printf("Checked %d accounts. %d still use the default password.\n\n", len(accounts), len(affected))
	for _, a := range affected {
		lastLogin := "never"
		if a.LastLoginAt != nil {
			lastLogin = a.LastLoginAt.Format("2006-01-02 15:04")
		}
		fmt.Printf("  %s  active=%-5t  last_login=%-16s  %s\n", a.ID, a.IsActive, lastLogin, a.Name)
	}

	if len(affected) == 0 {
		return nil
	}

	if !*apply {
		fmt.Println("\nDry run: nothing was changed. Run again with --apply to fix these accounts.")
		return nil
	}

	redisClient, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		return err
	}
	defer redisClient.Close()

	store := cache.NewRedisCache(redisClient)
	userRepo := user.NewRepository(db)
	outboxRepo := outbox.NewRepository(db)
	txManager := database.NewTxManager(db)

	fmt.Println()
	failed := 0
	for _, a := range affected {
		if err := fix(ctx, a, cfg.App.URL, store, userRepo, outboxRepo, txManager); err != nil {
			failed++
			fmt.Printf("  FAILED  %s  %s: %v\n", a.ID, a.Name, err)
			continue
		}

		if a.IsActive {
			fmt.Printf("  fixed   %s  %s (signed out, email queued)\n", a.ID, a.Name)
		} else {
			fmt.Printf("  fixed   %s  %s (inactive, no email)\n", a.ID, a.Name)
		}
	}

	fmt.Printf("\nDone: %d fixed, %d failed.\n", len(affected)-failed, failed)
	if failed > 0 {
		return fmt.Errorf("%d accounts could not be fixed", failed)
	}

	return nil
}

func fix(
	ctx context.Context,
	a account,
	appURL string,
	store cache.Cache,
	userRepo *user.Repository,
	outboxRepo outbox.Repository,
	txManager *database.TxManager,
) error {
	random, err := passwordtoken.New()
	if err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(random), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash random password: %w", err)
	}

	if err := userRepo.UpdatePassword(ctx, a.ID, string(hashed)); err != nil {
		return err
	}

	if err := jwtutil.RevokeUserSessions(ctx, store, a.ID); err != nil {
		return fmt.Errorf("revoke sessions: %w", err)
	}

	if !a.IsActive {
		return nil
	}

	resetURL, err := passwordtoken.Issue(ctx, store, appURL, a.ID, resetLinkTTL)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(event.PasswordResetRequestedEvent{
		Name:      a.Name,
		Email:     a.Email,
		ResetURL:  resetURL,
		ExpiresIn: resetLinkExpiryText,
		Forced:    true,
	})
	if err != nil {
		return fmt.Errorf("encode reset event: %w", err)
	}

	tx, err := txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := outboxRepo.Create(ctx, tx, outbox.OutboxEvent{
		EventType:   event.PasswordResetRequested,
		AggregateID: a.ID.String(),
		Payload:     payload,
	}); err != nil {
		return err
	}

	return tx.Commit()
}
