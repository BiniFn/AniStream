package mappers

import (
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
)

func TokenFromRepository(token repository.GetResetPasswordTokenRow) models.ChangePasswordTokenResponse {
	return models.ChangePasswordTokenResponse{
		User:      UserFromRepository(token.User),
		ExpiresAt: token.ResetPasswordToken.ExpiresAt.Time.UnixMilli(),
	}
}
