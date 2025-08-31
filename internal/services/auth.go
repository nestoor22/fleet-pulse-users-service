package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fleet-pulse-users-service/internal/config"
	"fleet-pulse-users-service/internal/errors"
	"fleet-pulse-users-service/internal/models"
	"fleet-pulse-users-service/internal/repositories"
	"fleet-pulse-users-service/internal/schemas"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(config.Get().Auth.JwtSecret)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type AuthService struct {
	refreshTokenRepository *repositories.RefreshTokenRepository
	userRepository         *repositories.UserRepository
}

func NewAuthService(refreshTokenRepository *repositories.RefreshTokenRepository, userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{refreshTokenRepository: refreshTokenRepository, userRepository: userRepository}
}

func (s AuthService) GenerateJWT(userID string, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s AuthService) ParseJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.ErrInvalidToken
	}
	return claims, nil
}

func (s AuthService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 256)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s AuthService) LoginUser(loginPayload schemas.LoginUserRequest) (string, string, error) {
	userObj, err := s.userRepository.GetUserByEmail(loginPayload.Email)
	settings := config.Get()

	if err != nil || userObj == nil {
		return "", "", errors.ErrUserNotFound
	}
	if !CheckPassword(userObj.Password, loginPayload.Password) {
		return "", "", errors.ErrInvalidCredentials
	}

	accessToken, err := s.GenerateJWT(userObj.ID.String(), time.Duration(settings.Auth.JwtAccessTokenExpireInMinutes)*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	s.refreshTokenRepository.DeletePreviousTokens(userObj.ID)
	_, err = s.refreshTokenRepository.Create(&models.RefreshToken{
		UserID:    userObj.ID,
		Token:     HashRefreshToken(refreshToken),
		ExpiresAt: time.Now().Add(time.Duration(settings.Auth.JwtRefreshTokenExpireInHours) * time.Hour),
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s AuthService) RefreshAccessToken(rawRefreshToken string) (newAccessToken string, newRefreshToken string, err error) {
	// Hash the incoming refresh token
	hashedToken := HashRefreshToken(rawRefreshToken)
	settings := config.Get()

	tokenObj, err := s.refreshTokenRepository.GetByToken(hashedToken)
	if err != nil || tokenObj == nil {
		return "", "", errors.ErrInvalidToken
	}

	if time.Now().After(tokenObj.ExpiresAt) {
		s.refreshTokenRepository.Delete(tokenObj)
		return "", "", errors.ErrExpiredToken
	}

	newAccessToken, err = s.GenerateJWT(
		tokenObj.UserID.String(),
		time.Duration(settings.Auth.JwtAccessTokenExpireInMinutes)*time.Minute,
	)
	if err != nil {
		return "", "", err
	}

	newRefreshTokenRaw, err := s.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	tokenObj.Token = HashRefreshToken(newRefreshTokenRaw)
	tokenObj.ExpiresAt = time.Now().Add(time.Duration(settings.Auth.JwtRefreshTokenExpireInHours) * time.Hour)
	if _, err := s.refreshTokenRepository.Update(tokenObj); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshTokenRaw, nil
}

func HashRefreshToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}
