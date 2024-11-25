package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/MarcusXavierr/api-banco-cooperativa/internal/db"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	DB *db.DBPool
}

func (service AuthService) generateTokenJWT(user db.User) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_HMAC")))
}
