package mappers

import (
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
)

func UserFromRepository(user repository.User) models.UserResponse {
	return models.UserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture.String,
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
	}
}

