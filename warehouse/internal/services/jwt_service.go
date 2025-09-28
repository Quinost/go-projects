package services

import (
	"time"
	"warehouse/internal/cfg"
	"warehouse/internal/constants"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JWTService struct {
	cfg *cfg.AuthCfg
}

func NewJWTService(cfg *cfg.AuthCfg) *JWTService {
	return &JWTService{
		cfg: cfg,
	}
}

func (s *JWTService) GenerateJWT(userId uuid.UUID, userName string) (string, error) {
	claims := jwt.MapClaims{
		constants.ClaimSub:      userId,
		constants.ClaimUserName: userName,
		constants.ClaimExp:      time.Now().Add(time.Duration(s.cfg.Exp) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}

func (s *JWTService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *JWTService) IsPasswordCorrect(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}