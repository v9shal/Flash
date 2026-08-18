package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type AuthService interface {
	RegisterUser(ctx context.Context, FullName, email, password string, phone *string) (*User, error)
	LoginUser(ctx context.Context, email, password string) (*User, string, string, error)
}
type authService struct {
	repo      UserRepository
	jwtSecret string
}

func NewAuthService(repo UserRepository, jwtSecret string) AuthService {
	return &authService{repo: repo, jwtSecret: jwtSecret}
}

func (s *authService) RegisterUser(
	ctx context.Context,
	fullName string,
	email string,
	password string,
	phone *string,
) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, fmt.Errorf("Error while hashing password %w", err)
	}
	user := &User{
		FullName:     fullName,
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hashedPassword),
	}
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}

	// 4. Return created user
	return user, nil

}
func (s *authService) LoginUser(ctx context.Context, email string,
	password string) (*User, string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", "", ErrInvalidCredentials

	}
	accessToken, err := GenerateAccessToken(user, s.jwtSecret)
	if err != nil {
		return nil, "", "", errors.New("error while generating access token")
	}
	refreshToken, err := GenerateRefreshToken(user, s.jwtSecret)
	if err != nil {
		return nil, "", "", errors.New("error while generating refresh token")
	}
	hashToken := HashToken(refreshToken)
	expiry := time.Now().Add(7 * 24 * time.Hour)
	err = s.repo.StoreRefreshToken(ctx, user.ID, hashToken, expiry)
	if err != nil {
		return nil, "", "", errors.New("Error while storing the refreshtoken in DB")
	}
	return user, accessToken, refreshToken, nil

}
