package oauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/big"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type Provider interface {
	Name() string
	AuthURL(ctx context.Context, state string) (string, error)
	ExchangeToken(ctx context.Context, state, code string) (TokenResponse, error)
}

var (
	ErrInvalidCodeVerifierLength = errors.New("invalid code verifier length")
)

type generateCodeVerifierParams struct {
	length int
}

func generateCodeVerifier(params generateCodeVerifierParams) (string, error) {
	length := params.length
	if length == 0 {
		length = 43
	}

	if length < 43 || length > 128 {
		return "", ErrInvalidCodeVerifierLength
	}

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
	res := make([]byte, length)
	for i := range res {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		res[i] = charset[n.Int64()]
	}
	return string(res), nil
}

func createCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
