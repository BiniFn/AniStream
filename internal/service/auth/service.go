package auth

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/coeeter/aniways/internal/infra/email"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/service/users"
	"github.com/coeeter/aniways/internal/template"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	repo        *repository.Queries
	emailClient email.EmailClient
	frontendURL string
}

func NewAuthService(repo *repository.Queries, emailClient email.EmailClient, frontendURL string) *AuthService {
	return &AuthService{
		repo:        repo,
		emailClient: emailClient,
		frontendURL: frontendURL,
	}
}

func (s *AuthService) buildSendEmailParams(to, token string) email.SendSimpleEmailParams {
	tokenURL := s.frontendURL + "/reset-password/" + token

	html := fmt.Sprintf(template.ForgetPasswordEmailTemplate, tokenURL)

	return email.SendSimpleEmailParams{
		To:      []string{to},
		Subject: "Reset your password",
		Html:    html,
	}
}

func (s *AuthService) SendForgetPasswordEmail(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}

	token, err := s.repo.CreateResetPasswordToken(ctx, user.ID)
	if err != nil {
		return err
	}

	err = s.emailClient.SendSimpleEmail(ctx, s.buildSendEmailParams(email, token.Token))
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetUserByForgetPasswordToken(ctx context.Context, token string) (users.User, error) {
	user, err := s.repo.GetUserByResetPasswordToken(ctx, token)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return users.User{}, users.ErrUserDoesNotExist
	}
	if err != nil {
		return users.User{}, err
	}

	return users.User{}.FromRepository(user), nil
}

var ErrInvalidToken = errors.New("invalid token")

func (s *AuthService) ResetPassword(ctx context.Context, userService *users.UserService, token, password string) error {
	user, err := s.GetUserByForgetPasswordToken(ctx, token)
	if err != nil {
		return ErrInvalidToken
	}

	if err := userService.ResetPassword(ctx, user.ID, password); err != nil {
		return err
	}

	err = s.repo.DeleteResetPasswordToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetConnectedProviders(ctx context.Context, userID string) ([]string, error) {
	providers, err := s.repo.GetAllOauthTokensOfUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	out := make([]string, len(providers))
	for i, provider := range providers {
		out[i] = string(provider.Provider)
	}
	return out, nil
}

func (s *AuthService) DisconnectProvider(ctx context.Context, userID string, provider string) error {
	return s.repo.DeleteOauthToken(ctx, repository.DeleteOauthTokenParams{
		UserID:   userID,
		Provider: repository.Provider(provider),
	})
}
