package route

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

// 请求响应函数
func funcIndex(ctx context.Context, c *app.RequestContext) {
	c.HTML(http.StatusOK, ThemePath+"index.html", map[string]string{"name": "test"})
}

// newJWTToken 创建一个jwtToken
func newJWTToken(userId string, info map[string]interface{}) (string, error) {
	if userId == "" {
		return "", errors.New("userId不能为空")
	}
	mapClaims := jwt.MapClaims{}
	for k, v := range info {
		mapClaims[k] = v
	}
	mapClaims["exp"] = jwt.NewNumericDate(time.Now().Add(time.Duration(service.Config.Timeout) * time.Second))
	mapClaims["iat"] = jwt.NewNumericDate(time.Now())
	mapClaims["nbf"] = jwt.NewNumericDate(time.Now())
	mapClaims["iss"] = constant.DEFAULT_NAME
	mapClaims["jti"] = userId

	mapClaims[constant.TOKEN_USER_ID] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, mapClaims)
	tokenString, err := token.SignedString([]byte(service.Config.JwtSecret))
	return tokenString, err
}

// adminHandler admin权限拦截器
func adminHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		jwttoken := c.Cookie(service.Config.JwtTokenKey)
		userId, err := userIdByToken(string(jwttoken))
		if err != nil || userId == "" {
			c.Redirect(http.StatusOK, []byte("/admin/login"))
			c.Abort() // 终止后续调用
			return
		}
		// 传递从jwttoken获取的userId
		c.Set(constant.TOKEN_USER_ID, userId)
	}
}
func userIdByToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("token不能为空")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(service.Config.JwtSecret), nil
	})
	if !token.Valid {
		return "", errors.New("token is not valid")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return "", fmt.Errorf("that's not even a token:%w", err)
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return "", fmt.Errorf("timing is everything:%w", err)
	} else if err != nil {
		return "", fmt.Errorf("couldn't handle this token:%w", err)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId := mapClaims[constant.TOKEN_USER_ID].(string)
		return userId, nil
	}
	return "", errors.New("token错误或过期")
}
