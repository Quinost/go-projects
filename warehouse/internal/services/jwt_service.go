package services

import (
	"time"
	"warehouse/internal/cfg"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	cfg *cfg.AuthCfg
}

func NewJWTService(cfg *cfg.AuthCfg) *JWTService {
	return &JWTService{
		cfg: cfg,
	}
}

func (s *JWTService) GenerateJWT(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}
