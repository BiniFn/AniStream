package auth

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/coeeter/aniways/internal/email"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/template"
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
