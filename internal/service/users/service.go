package users

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/coeeter/aniways/internal/mappers"
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooLong  = fmt.Errorf("password too long")
	ErrUserDoesNotExist = fmt.Errorf("user does not exist")
	ErrUsernameTaken    = fmt.Errorf("username taken")
	ErrEmailTaken       = fmt.Errorf("email taken")
	ErrInvalidAuth      = fmt.Errorf("invalid authentication")
)

type UserService struct {
	repo *repository.Queries
	cld  *cloudinary.Cloudinary
}

func NewUserService(repo *repository.Queries, cld *cloudinary.Cloudinary) *UserService {
	return &UserService{
		repo: repo,
		cld:  cld,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (models.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return models.UserResponse{}, err
	}
	return mappers.UserFromRepository(user), nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (models.UserResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return models.UserResponse{}, err
	}
	return mappers.UserFromRepository(user), nil
}

func (s *UserService) CreateUser(ctx context.Context, username, email, password string) (models.UserResponse, error) {
	passwordsBytes := []byte(password)

	if len(passwordsBytes) > 72 {
		return models.UserResponse{}, ErrPasswordTooLong
	}

	hash, err := bcrypt.GenerateFromPassword(passwordsBytes, bcrypt.DefaultCost)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, repository.CreateUserParams{
		Username:       username,
		Email:          email,
		PasswordHash:   string(hash),
		ProfilePicture: pgtype.Text{String: "", Valid: false},
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "users_username_key") {
				return models.UserResponse{}, ErrUsernameTaken
			}
			if strings.Contains(err.Error(), "users_email_key") {
				return models.UserResponse{}, ErrEmailTaken
			}
		}
		return models.UserResponse{}, err
	}
	return mappers.UserFromRepository(user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, id, username, email string) (models.UserResponse, error) {
	user, err := s.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:       id,
		Username: username,
		Email:    email,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "users_username_key") {
				return models.UserResponse{}, ErrUsernameTaken
			}
			if strings.Contains(err.Error(), "users_email_key") {
				return models.UserResponse{}, ErrEmailTaken
			}
		}
		return models.UserResponse{}, err
	}
	return mappers.UserFromRepository(user), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id, password string) error {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return ErrInvalidAuth
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return ErrInvalidAuth
	}

	err = s.repo.DeleteUser(ctx, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (models.UserResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return models.UserResponse{}, ErrInvalidAuth
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.UserResponse{}, ErrInvalidAuth
	}

	return mappers.UserFromRepository(user), nil
}

func (s *UserService) CreateSession(ctx context.Context, userID string) (repository.Session, error) {
	session, err := s.repo.CreateSession(ctx, userID)
	if err != nil {
		return repository.Session{}, err
	}
	return session, nil
}

func (s *UserService) DeleteSession(ctx context.Context, sessionID string) error {
	err := s.repo.DeleteSession(ctx, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserBySessionID(ctx context.Context, sessionID string) (models.UserResponse, error) {
	user, err := s.repo.GetUserBySessionID(ctx, sessionID)
	if err != nil {
		return models.UserResponse{}, err
	}
	return mappers.UserFromRepository(user), nil
}

func (s *UserService) UpdatePassword(ctx context.Context, id, oldPassword, newPassword string) error {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return ErrInvalidAuth
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidAuth
	}

	newPasswordBytes := []byte(newPassword)

	if len(newPasswordBytes) > 72 {
		return ErrPasswordTooLong
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.repo.UpdatePassword(ctx, repository.UpdatePasswordParams{
		ID:           id,
		PasswordHash: string(hash),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, id, newPassword string) error {
	newPasswordBytes := []byte(newPassword)

	if len(newPasswordBytes) > 72 {
		return ErrPasswordTooLong
	}

	hash, err := bcrypt.GenerateFromPassword(newPasswordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.repo.UpdatePassword(ctx, repository.UpdatePasswordParams{
		ID:           id,
		PasswordHash: string(hash),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateProfilePicture(ctx context.Context, id string, image io.Reader) error {
	if _, err := s.repo.GetUserByID(ctx, id); err != nil {
		return ErrUserDoesNotExist
	}

	result, err := s.cld.Upload.Upload(ctx, image, uploader.UploadParams{})
	if err != nil {
		return fmt.Errorf("cloudinary upload: %w", err)
	}

	if err := s.repo.UpdateProfilePicture(ctx, repository.UpdateProfilePictureParams{
		ID:             id,
		ProfilePicture: pgtype.Text{String: result.SecureURL, Valid: len(result.SecureURL) > 0},
	}); err != nil {
		return fmt.Errorf("db update profile picture: %w", err)
	}
	return nil
}

func (s *UserService) RemoveProfilePicture(ctx context.Context, id string) error {
	if _, err := s.repo.GetUserByID(ctx, id); err != nil {
		return ErrUserDoesNotExist
	}

	if err := s.repo.RemoveProfilePicture(ctx, id); err != nil {
		return fmt.Errorf("db remove profile picture: %w", err)
	}
	return nil
}
