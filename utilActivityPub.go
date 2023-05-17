package main

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Signature struct {
	KeyID     string `json:"keyId"`
	Algorithm string `json:"algorithm"`
	Headers   string `json:"headers"`
	Value     string `json:"signature"`
}

func parseSignature(signatureString string) (*Signature, error) {
	//逗号分割签名的字符串
	s1 := strings.Split(signatureString, ",")
	sigMap := make(map[string]string, 0)
	for _, s2 := range s1 {
		dIndex := strings.Index(s2, "=")
		if dIndex < 0 {
			continue
		}
		sigMap[s2[:dIndex]] = strings.Trim(strings.TrimSpace(s2[dIndex+1:]), `"`)
	}
	sigByte, _ := json.Marshal(sigMap)
	// 解析签名字符串，提取相关信息
	signature := &Signature{}
	err := json.Unmarshal(sigByte, signature)
	if err != nil {
		return nil, err
	}
	return signature, nil
}
func getPublicKey(publicKeyID string) (*rsa.PublicKey, error) {
	// 根据公钥 ID 获取对应的公钥
	// 这里使用假数据，实际使用时需要替换为真实的公钥获取逻辑
	//publicKeyPEM := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAr9HicDyHYlpGVYVHrm7j\nU7Nq4z9SeynK8UUi+JoBWuotChg2oSDQtWuj+zdQSKM3g27+sqNNw/BuZp85BVT6\n8PRyamTHjVrZPj6JIC+A/EGeJTqycODoMTDTTdz3evxBUbPAH7By91VrMNE5i8zl\nJ40IqAYYNLjmUdvQliGmGpX/xmPAfIeJ/mMQ3kCq/2uSICrL1ORicAB/qqXgyPsB\nWZCTYOOdJsV9bbbhAQUqRjevZrRIdaVcrIObxTDY0VgtBJgsElGNxbnb/g4vfPgy\nWdi/E0qLSRyayml8lGZhPccgY3PnqGO765X/j0tra/I4JIjLC0AOV0nLs0fLmH72\nEwIDAQAB\n-----END PUBLIC KEY-----\n"
	publicKeyPEM, err := requestJsonValue(publicKeyID, "publicKey.publicKeyPem")
	if err != nil {
		return nil, err
	}

	// 解析公钥 PEM 格式
	block, _ := pem.Decode([]byte(publicKeyPEM.(string)))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("公钥解析失败")
	}
	// 解析公钥 DER 格式
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 转换为 RSA 公钥
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("公钥类型错误")
	}

	return publicKey, nil
}

func buildSignatureData(c *app.RequestContext, headers string) string {
	// 构建签名字符串
	var comparisonStrings []string
	signedHeaders := strings.Split(headers, " ")
	for _, header := range signedHeaders {
		value := ""
		header = strings.TrimSpace(header)
		if header == "(request-target)" {
			method := string(c.Method())
			method = strings.ToLower(method)
			uri := string(c.Request.URI().Path())
			value = fmt.Sprintf("%s %s", method, uri)
		} else {
			value = string(c.GetHeader(header))
		}
		comparisonStrings = append(comparisonStrings, header+": "+value)
	}
	return strings.Join(comparisonStrings, "\n")
}

func verifySignature(publicKey *rsa.PublicKey, data string, signatureValue string) bool {
	// 验证签名
	hashed := sha256.Sum256([]byte(data))
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureValue)
	if err != nil {
		return false
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signatureBytes)
	if err != nil {
		return false
	}

	return true
}

func requestJsonValue(httpurl string, key string) (interface{}, error) {

	// 解析 URL 字符串
	parsedURL, err := url.Parse(httpurl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL:%w", err)
	}
	// 获取域名或 IP 地址
	host := parsedURL.Hostname()
	c, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	//设置翻墙代理
	c.SetProxy(protocol.ProxyURI(protocol.ParseURI("http://127.0.0.1:49864/")))

	response := &protocol.Response{}
	request := &protocol.Request{}
	request.SetMethod(consts.MethodGet)
	request.Header.SetContentTypeBytes([]byte(activityPubContentType))
	request.SetHeader("Accept", activityPubAccept)
	request.SetHeader("Host", host)
	request.SetRequestURI(httpurl)
	err = c.Do(context.Background(), request, response)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() == 301 || response.StatusCode() == 302 { //重定向
		location := response.Header.Get("location")
		request.SetRequestURI(location)
		err = c.Do(context.Background(), request, response)
		if err != nil {
			return nil, err
		}
	}

	bodyMap := make(map[string]interface{})
	json.Unmarshal(response.Body(), &bodyMap)

	keys := strings.Split(key, ".")

	var value interface{}
	for i, k := range keys {
		value = bodyMap[k]
		if i+1 < len(keys) {
			bodyMap = bodyMap[k].(map[string]interface{})
		}

	}
	return value, nil
}
