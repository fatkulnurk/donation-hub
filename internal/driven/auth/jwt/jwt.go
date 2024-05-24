package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/isdzulqor/donation-hub/internal/core/service/auth_token"
	"strconv"
	"time"
)

type service struct {
	secretKey string
	issuer    string
}

type MyCustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func (s service) GenerateToken(i model.AuthPayload) (token string, err error) {
	userIDStr := strconv.FormatInt(i.UserID, 10)
	claims := MyCustomClaims{
		UserID:   userIDStr,
		Username: i.Username,
		Email:    i.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "for user " + i.Username,
			Audience:  jwt.ClaimStrings{},
		},
	}
	fmt.Println(userIDStr)
	fmt.Println(claims)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = jwtToken.SignedString([]byte(s.secretKey))

	return
}

func (s service) ValidateToken(tokenString string) (*model.AuthPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		err = fmt.Errorf("error parsing token: %v", err)
		return nil, err
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if ok && token.Valid {
		fmt.Println("id:"+claims.UserID, "username:"+claims.Username, "email:"+claims.Email)
		id, err := strconv.Atoi(claims.UserID)
		if err != nil {
			return nil, err
		}

		return &model.AuthPayload{
			UserID:   int64(id),
			Username: claims.Username,
			Email:    claims.Email,
		}, nil
	}

	return nil, errors.New("invalid token")
}

func New(cfg model.ConfigMap) auth_token.Service {
	return &service{
		secretKey: cfg.SecretKey,
		issuer:    cfg.Issuer,
	}
}
