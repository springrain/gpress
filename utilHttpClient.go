package main

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
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
	//proxyAddress             = ""
	//proxyAddress             = "http://127.0.0.1:54321"
	//keyId                    = "https://" + activityPubDefaultDomain + "/activitypub/api/user/test11#main-key"
)

// responseJsonValue 发送GET请求,把返回值的json获取指定的key,例如
func responseJsonValue(httpurl string, jsonKey string, keyId string) (interface{}, error) {
	headerMap, _ := wrapRequestHeader(httpurl, consts.MethodGet, nil, keyId, false)
	bodyMap, err := sendActivityPubRequest(httpurl, consts.MethodGet, nil, headerMap)
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

// sendActivityPubRequest 发送activityPub请求
func sendActivityPubRequest(httpurl string, method string, bodyByte []byte, headerMap map[string]string) (map[string]interface{}, error) {
	// 解析 URL 字符串
	parsedURL, err := url.Parse(httpurl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL:%w", err)
	}
	// 获取域名或 IP 地址
	host := parsedURL.Hostname()

	c, err := client.NewClient(client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return nil, err
	}

	//设置翻墙代理
	if config.Proxy != "" && !strings.HasPrefix(httpurl, "http://127.0.0.1") {
		c.SetProxy(protocol.ProxyURI(protocol.ParseURI(config.Proxy)))
	}

	request := &protocol.Request{}
	request.SetRequestURI(httpurl)
	request.SetMethod(method)
	request.SetBody(bodyByte)
	request.Header.SetContentTypeBytes([]byte(activityPubContentType))
	request.SetHeader("Accept", activityPubAccept)
	request.SetHeader("Host", host)

	for k, v := range headerMap {
		request.SetHeader(k, v)
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

func wrapRequestHeader(httpurl string, method string, bodyByte []byte, keyId string, isSign bool) (map[string]string, error) {
	headerMap := make(map[string]string)
	// 解析 URL 字符串
	parsedURL, err := url.Parse(httpurl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL:%w", err)
	}

	// mastodon 的日期格式为RFC2616 -- time.RFC1123
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

	hash := sha256.Sum256(bodyByte)
	digest := "SHA-256=" + base64.StdEncoding.EncodeToString(hash[:])

	headerMap["Digest"] = digest
	headerMap["Date"] = date

	method = strings.ToLower(method)
	if method == "post" {
		isSign = true
	}
	if isSign {
		headers := "(request-target) date digest content-type"
		comparisonStrings := make([]string, 0)
		comparisonStrings = append(comparisonStrings, "(request-target): "+method+" "+parsedURL.Path)
		//先不要host,避免转发服务器验签失败,可以转发到其他服务器
		//comparisonStrings = append(comparisonStrings, "host: "+host)
		comparisonStrings = append(comparisonStrings, "date: "+date)
		comparisonStrings = append(comparisonStrings, "digest: "+digest)
		comparisonStrings = append(comparisonStrings, "content-type: "+activityPubContentType)
		signatureString := strings.Join(comparisonStrings, "\n")

		//生成签名数据
		signatureBase64, err := generateRSASignature(signatureString)
		if err != nil {
			return nil, err
		}
		// 构建 signature
		signature := fmt.Sprintf(`keyId="%s",algorithm="rsa-sha256",headers="%s",signature="%s"`, keyId, headers, signatureBase64)
		//fmt.Println("-----signature:" + signature)
		headerMap["Signature"] = signature
	}

	return headerMap, nil

}
