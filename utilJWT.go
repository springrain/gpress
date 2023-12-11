// Copyright (c) 2023 gpress Authors.
//
// This file is part of gpress.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	// 定义您需要加入的有效载荷字段
	// 这里只作为示例，可以根据实际需求进行定义
	UserId  string `json:"userId"`
	Expires int64  `json:"exp"`
}

// 定义头部和有效载荷
var defaultHeader = Header{
	Alg: "HS512",
	Typ: "JWT",
}

// newJWTToken 创建一个jwtToken
func newJWTToken(userId string) (string, error) {
	if userId == "" {
		return "", errors.New("userId is nil ")
	}

	payload := Payload{
		UserId:  userId,
		Expires: time.Now().Add(time.Duration(config.Timeout) * time.Second).Unix(), // 设置过期时间
	}

	// 序列化头部和有效载荷
	headerBytes, err := json.Marshal(defaultHeader)
	if err != nil {
		return "", err
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Base64编码头部和有效载荷
	headerString := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadString := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// 计算签名
	hasher := hmac.New(sha512.New, []byte(config.JwtSecret))
	hasher.Write([]byte(headerString + "." + payloadString))
	signature := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))

	// 拼接JWT字符串
	jwtString := fmt.Sprintf("%s.%s.%s", headerString, payloadString, signature)

	return jwtString, err
}

func userIdByToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("token is nil")
	}

	// 分割JWT字符串
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid JWT format")
	}

	// 解码头部和有效载荷部分
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	// 解析头部
	var header Header
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		return "", err
	}

	// 验证算法
	if header.Alg != "HS512" {
		return "", err
	}

	// 计算签名
	hasher := hmac.New(sha512.New, []byte(config.JwtSecret))
	hasher.Write([]byte(parts[0] + "." + parts[1]))
	expectedSignature := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))

	// 比较签名
	if parts[2] != expectedSignature {
		return "", errors.New("JWT signature is invalid")
	}

	// 解析有效载荷
	var payload Payload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return "", err
	}

	// 检查有效载荷中的过期时间
	// 根据您的需求进行适当的时间验证
	if payload.Expires < time.Now().Unix() {
		return "", errors.New("JWT signature is expires")
	}
	return payload.UserId, nil
}
