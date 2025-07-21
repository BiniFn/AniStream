package auth

import "github.com/coeeter/aniways/internal/repository"

type AuthService struct {
	repo *repository.Queries
}

func NewAuthService(repo *repository.Queries) *AuthService {
	return &AuthService{
		repo: repo,
	}
}
