package authservice

import (
	"GameApp/entity"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type Service struct {
	SignKey               string
	accessExpirationTime  time.Duration
	refreshExpirationTime time.Duration
	accessSubject         string
	refreshSubject        string
}

func New(signKey string, accessExpirationTime time.Duration, refreshExpirationTime time.Duration, accessSubject string, refreshSubject string) Service {
	return Service{
		SignKey:               signKey,
		accessExpirationTime:  accessExpirationTime,
		refreshExpirationTime: refreshExpirationTime,
		accessSubject:         accessSubject,
		refreshSubject:        refreshSubject,
	}
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.accessExpirationTime, s.accessSubject)

}
func (s Service) CreateRefreshToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.refreshExpirationTime, s.refreshSubject)

}

func (s Service) createToken(userID uint, expireDuration time.Duration, subject string) (string, error) {
	// create a signer for rsa 256
	// TODO : replace with rsa 256 RS256
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
		UserID: userID,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(s.SignKey))
	if err != nil {
		return "", err
	}
	return accessTokenString, nil

}

func (s Service) ParsToken(BearerToken string) (*Claims, error) {
	tokenStr := strings.Replace(BearerToken, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.SignKey), nil
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
