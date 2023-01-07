package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type gpressJWTClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func newJWTToken(userId string, info map[string]interface{}) (string, error) {
	if userId == "" {
		return "", errors.New("userId不能为空")
	}
	mapClaims := jwt.MapClaims{}
	if info != nil {
		for k, v := range info {
			mapClaims[k] = v
		}
	}
	mapClaims["exp"] = jwt.NewNumericDate(time.Now().Add(time.Duration(timeout) * time.Second))
	mapClaims["iat"] = jwt.NewNumericDate(time.Now())
	mapClaims["nbf"] = jwt.NewNumericDate(time.Now())
	mapClaims["iss"] = defaultName
	mapClaims["jti"] = userId
	mapClaims["userId"] = userId

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, mapClaims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	return tokenString, err

}

func userIdByToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if !token.Valid {
		return "", errors.New("token is not Valid")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return "", fmt.Errorf("that's not even a token:%w", err)
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		// Token is either expired or not active yet
		return "", fmt.Errorf("timing is everything:%w", err)
	} else if err != nil {
		return "", fmt.Errorf("couldn't handle this token:%w", err)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId := mapClaims["userId"].(string)
		return userId, nil
	}
	return "", errors.New("token错误或过期")
}
