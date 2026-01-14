package service

import (
	"effective-project/internal/repository"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTConfig struct {
	SecretKey      string
	AccessTokenTTL time.Duration
	logger         *slog.Logger
}

var ErrInvalidCredentials = errors.New("неправильный email или пароль")

type UserClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Login(email, password string) (string, error)

	GenerateToken(userID uuid.UUID, role string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtCfg   JWTConfig
	logger   *slog.Logger
}

func NewAuthService(
	userRepo repository.UserRepository,
	jwtCfg JWTConfig,
	logger *slog.Logger,
) AuthService {
	return &authService{userRepo: userRepo, jwtCfg: jwtCfg, logger: logger}
}

func (s *authService) Login(email, password string) (string, error) {
	s.logger.Debug("Попытка входа", "email", email)

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		s.logger.Warn("неверные учетные данные — пользователь не найден", "email", email)
		return "", ErrInvalidCredentials
	}

	if err := checkPassword(user.Password, password); err != nil {
		s.logger.Warn("неверные учетные данные — неверный пароль", "email", email)
		return "", ErrInvalidCredentials
	}

	now := time.Now()
	claims := UserClaims{
		UserID: user.ID,
		Role:   string(user.Roles),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtCfg.AccessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(s.jwtCfg.SecretKey))
	if err != nil {
		s.logger.Error("ошибка при подписании JWT", "error", err, "user_id", user.ID)
		return "", err
	}

	s.logger.Info("пользователь вошёл", "user_id", user.ID, "email", email)
	return signed, nil
}

func (s *authService) GenerateToken(userID uuid.UUID, role string) (string, error) {
	s.logger.Debug("GenerateToken вызван", "user_id", userID)
	now := time.Now()

	claims := UserClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtCfg.AccessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(s.jwtCfg.SecretKey))
	if err != nil {
		s.logger.Error("ошибка при подписании токена", "error", err, "user_id", userID)
		return "", err
	}
	s.logger.Info("токен сгенерирован", "user_id", userID)
	return signed, nil
}
