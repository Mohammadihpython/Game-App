package authservice

import (
	"GameApp/entity"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type Config struct {
	SignKey               string        `koanf:"sign_key"`
	AccessExpirationTime  time.Duration `koanf:"access_expiration_time"`
	RefreshExpirationTime time.Duration `koanf:"refresh_expiration_time"`

	AccessSubject string `koanf:"access_subject"`

	RefreshSubject string `koanf:"refresh_subject"`
}
type Service struct {
	config Config
}

func New(cfg Config) Service {
	return Service{
		config: cfg}
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
	Role entity.Role `json:"role"`
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, user.Role, s.config.AccessExpirationTime, s.config.AccessSubject)

}
func (s Service) CreateRefreshToken(user entity.User) (string, error) {
	return s.createToken(user.ID, user.Role, s.config.RefreshExpirationTime, s.config.RefreshSubject)

}

func (s Service) createToken(userID uint, role entity.Role, expireDuration time.Duration, subject string) (string, error) {
	// create a signer for rsa 256
	// TODO : replace with rsa 256 RS256
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
		UserID: userID,
		Role:   role,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.SignKey))
	if err != nil {
		return "", err
	}
	return accessTokenString, nil

}

func (s Service) ParsToken(BearerToken string) (*Claims, error) {
	tokenStr := strings.Replace(BearerToken, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SignKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Println(claims.UserID)
		return claims, nil
	} else {
		return nil, err
	}

}
