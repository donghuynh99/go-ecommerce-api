package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/database"
	"github.com/donghuynh99/ecommerce_api/models"
	"gorm.io/gorm/clause"
)

func GenerateToken(userId int, role string) (string, time.Time, error) {
	timeExpiration := time.Now().Add(time.Hour * 24 * 7) // 1 week

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userId
	claims["role"] = role
	claims["exp"] = timeExpiration

	signedTolken, err := token.SignedString([]byte(config.GetConfig().GeneralConfig.KeyToken))

	if err != nil {
		return "", time.Time{}, err
	}

	return signedTolken, timeExpiration, nil
}

func CreateSignature(keyString string) string {
	key := []byte(config.GetConfig().GeneralConfig.KeyToken)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(keyString))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return signature
}

func VerifySignature(keyString string, signature string) bool {
	expectedSignature := CreateSignature(keyString)

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func ValidateToken(tokenString string, role string) (models.User, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.GetConfig().GeneralConfig.KeyToken), nil
	})

	if err != nil {
		return models.User{}, err
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)

		if ok {
			userId := claims["user_id"]
			roleClaim := claims["role"].(string)

			var user models.User
			err := database.Database.Preload(clause.Associations).First(&user, userId).Error

			if err != nil {
				return models.User{}, err
			}

			if roleClaim == role && user.ExpiredAt != nil && user.Token != nil && user.ExpiredAt.After(time.Now()) && tokenString == *user.Token {
				return user, nil
			}
		}
	}

	return models.User{}, errors.New(Translation("invalid_token", nil, nil))
}
