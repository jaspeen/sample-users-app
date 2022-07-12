package services

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jaspeen/sample-users-app/config"
	"github.com/jaspeen/sample-users-app/db"
	"github.com/jinzhu/gorm"
)

func GenToken(secret string, user *db.User, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserToClaims(user, jwt.MapClaims{
		"exp": jwt.NewNumericDate(expiration),
	}))
	tokenString, err := token.SignedString([]byte(secret))
	log.Printf("Token gen: %v, err: %v\n", tokenString, err)
	return tokenString, err
}

func ParseAndValidateToken(secret string, rawToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	return token, err
}

func GenerateRefreshAndAccessToken(tx *gorm.DB, user *db.User) (string, string, error) {
	token, err := GenToken(config.C.TokenSecret, user, time.Now().Add(1*time.Hour))
	if err != nil {
		return "", "", err
	}
	refreshExpTime := time.Now().Add(24 * 30 * time.Hour)
	refreshToken, err := GenToken(config.C.RefreshTokenSecret, user, refreshExpTime)
	if err != nil {
		return "", "", err
	}
	err = tx.Save(&db.RefreshToken{
		UserId:     user.ID,
		Token:      refreshToken,
		Expiration: refreshExpTime,
	}).Error

	return token, refreshToken, err
}

func UserToClaims(user *db.User, claims jwt.MapClaims) jwt.MapClaims {
	claims["iss"] = "sample-users-app"
	claims["sub"] = user.Email
	claims["userid"] = user.ID
	claims["adm"] = user.Admin
	return claims
}

func ClaimsToUser(claims jwt.MapClaims, user *db.User) *db.User {
	user.ID = int(claims["userid"].(float64))
	user.Email = claims["sub"].(string)
	user.Admin = claims["adm"].(bool)
	return user
}

func RenewAccessToken(tx *gorm.DB, refreshToken string) (string, error) {
	var token db.RefreshToken
	res := tx.Take(&token, "token = ?", refreshToken)
	if res.Error != nil || res.RowsAffected == 0 {
		return "", Err_NOT_FOUND
	}
	if time.Now().After(token.Expiration) {
		return "", Err_UNAUTHENTICATED
	}

	parsedToken, err := ParseAndValidateToken(config.C.RefreshTokenSecret, refreshToken)
	if err != nil {
		return "", err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {

		fmt.Println(claims["sub"])
		userFromClaims := ClaimsToUser(claims, &db.User{})
		accessToken, err := GenToken(config.C.TokenSecret, userFromClaims, time.Now().Add(2*time.Hour))
		if err != nil {
			return "", err
		}
		return accessToken, nil
	} else {
		return "", Err_UNAUTHENTICATED
	}
}
