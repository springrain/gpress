package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const (
	activityPubAccept        = "application/activity+json"
	activityPubContentType   = "application/activity+json; charset=utf-8"
	activityPubDefaultDomain = "activitypub.gpress.cn"
	keyId                    = "https://" + activityPubDefaultDomain + "/activitypub/api/user/test11#main-key"
)

// responseJsonValue 发送GET请求,把返回值的json获取指定的key,例如
func responseJsonValue(httpurl string, jsonKey string) (interface{}, error) {

	bodyMap, err := sendRequest(httpurl, consts.MethodGet, nil, false)
	if err != nil {
		return nil, err
	}
	if jsonKey == "" {
		return bodyMap, nil
	}

	keys := strings.Split(jsonKey, ".")

	var value interface{}
	for i, k := range keys {
		value = bodyMap[k]
		if i+1 < len(keys) {
			if bodyMap[k] == nil {
				return value, nil
			}
			bodyMap = bodyMap[k].(map[string]interface{})
		}

	}
	return value, nil
}

// funcSendPostAndSignature 返送post请求,并签名
func sendRequest(httpurl string, method string, body map[string]interface{}, isSign bool) (map[string]interface{}, error) {
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

	// mastodon 的日期格式为RFC2616 -- time.RFC1123
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

	bodyByte := []byte("")
	if body != nil {
		bodyByte, _ = json.Marshal(body)
	}
	hash := sha256.Sum256(bodyByte)
	digest := "SHA-256=" + base64.StdEncoding.EncodeToString(hash[:])

	request := &protocol.Request{}
	request.SetRequestURI(httpurl)
	request.SetMethod(method)
	request.SetBody(bodyByte)
	request.Header.SetContentTypeBytes([]byte(activityPubContentType))
	request.SetHeader("Accept", activityPubAccept)
	request.SetHeader("Host", host)
	request.SetHeader("Digest", digest)
	request.SetHeader("Date", date)

	if isSign {
		headers := "(request-target) host date digest content-type"
		comparisonStrings := make([]string, 0)
		comparisonStrings = append(comparisonStrings, "(request-target): "+strings.ToLower(method)+" "+parsedURL.Path)
		comparisonStrings = append(comparisonStrings, "host: "+host)
		comparisonStrings = append(comparisonStrings, "date: "+date)
		comparisonStrings = append(comparisonStrings, "digest: "+digest)
		comparisonStrings = append(comparisonStrings, "content-type: "+activityPubContentType)
		signatureString := strings.Join(comparisonStrings, "\n")

		//生成签名数据
		signatureBase64, err := makeSignature(signatureString)
		if err != nil {
			return nil, err
		}
		// 构建 signature
		signature := fmt.Sprintf(`keyId="%s",algorithm="rsa-sha256",headers="%s",signature="%s"`, keyId, headers, signatureBase64)
		//fmt.Println("-----signature:" + signature)
		request.SetHeader("Signature", signature)
	}
	response := &protocol.Response{}
	err = c.Do(context.Background(), request, response)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() == 301 || response.StatusCode() == 302 { //一次重定向
		location := response.Header.Get("location")
		request.SetRequestURI(location)
		err = c.Do(context.Background(), request, response)
		if err != nil {
			return nil, err
		}
	}

	bodyMap := make(map[string]interface{})
	json.Unmarshal(response.Body(), &bodyMap)

	return bodyMap, nil

}
