package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"JWTService/internal/models"
	"JWTService/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtSecret  = "fdsfsdf"
	accessExp  = 15 * time.Minute
	refreshExp = 24 * time.Hour
)

type AuthService struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.TokenRepository
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, input models.CreateUserInput) (*TokenPair, error) {
	exists, err := s.userRepo.IsEmailExists(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:      input.Username,
		Email:         input.Email,
		Password_hash: string(hashedPassword),
		CreatedAt:     time.Now(),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return s.generateTokenPair(ctx, user.Id)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(password)); err != nil {
		return nil, errors.New("Wrong password")
	}

	return s.generateTokenPair(ctx, user.Id)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	tokenID, ok := claims["jti"].(string)
	if !ok {
		return nil, errors.New("invalid token id")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user id")
	}
	userID := int64(userIDFloat)

	storedToken, err := s.tokenRepo.GetRefreshToken(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	if storedToken == nil || storedToken.Revoked || storedToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired or revoked")
	}

	if err := s.tokenRepo.RevokeRefreshToken(ctx, tokenID); err != nil {
		return nil, err
	}

	return s.generateTokenPair(ctx, userID)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return err
	}

	tokenID, ok := claims["jti"].(string)
	if !ok {
		return errors.New("invalid token id")
	}

	return s.tokenRepo.RevokeRefreshToken(ctx, tokenID)
}

func (s *AuthService) generateTokenPair(ctx context.Context, userID int64) (*TokenPair, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"exp":        time.Now().Add(accessExp).Unix(),
		"token_type": "access",
	})

	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	tokenID := uuid.New().String()
	expiresAt := time.Now().Add(refreshExp)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"jti":        tokenID,
		"exp":        expiresAt.Unix(),
		"token_type": "refresh",
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	if err := s.tokenRepo.SaveRefreshToken(ctx, userID, tokenID, expiresAt); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (s *AuthService) LogoutAll(ctx context.Context, refreshToken string) error {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return errors.New("invalid token")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return errors.New("invalid token payload")
	}
	userID := int64(userIDFloat)

	return s.tokenRepo.DeleteAllByUserID(ctx, userID)
}
