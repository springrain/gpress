package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/piprate/json-gold/ld"
)

/**
https://lawrenceli.me/blog/activitypub
https://wangqiao.me/posts/activitypub-from-decentralized-to-distributed-social-networks/
https://blog.joinmastodon.org/2018/06/how-to-implement-a-basic-activitypub-server/
https://www.w3.org/TR/activitystreams-vocabulary/
**/

func TestExpand(t *testing.T) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	// expanding in-memory document
	doc := map[string]interface{}{
		"@context":          "https://www.w3.org/ns/activitystreams",
		"id":                "https://lawrenceli.me/api/activitypub/actor",
		"type":              "Person",
		"name":              "Lawrence Li",
		"preferredUsername": "lawrence",
		"summary":           "Blog",
		"inbox":             "https://lawrenceli.me/api/activitypub/inbox",
		"outbox":            "https://lawrenceli.me/api/activitypub/outbox",
		"followers":         "https://lawrenceli.me/api/activitypub/followers",
	}

	expanded, err := proc.Expand(doc, options)
	if err != nil {
		t.Error("Error when expanding JSON-LD document:", err)
		return
	}
	fmt.Println(expanded)
}

func TestSendActivity(t *testing.T) {
	// 构建Activity对象
	activity := map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
		"type":     "Create",
		"to": []string{
			"https://mastodon.social/users/9iuorg",
		},
		"actor": "https://testmakerone.shengjian.net/activitypub/api/user/test11",
		"object": map[string]interface{}{
			"type":      "Note",
			"content":   "Hello ActivityPub!",
			"published": "2023-05-16T12:00:00Z",
		},
	}

	// 将Activity对象转换为JSON
	jsonData, err := json.Marshal(activity)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return
	}

	// 创建HTTP POST请求
	url := "https://mastodon.social/users/9iuorg/inbox" // 目标服务器的接收地址
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("HTTP request error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/activity+json")

	// 发送HTTP请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP request error:", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Activity sent successfully!")
	} else {
		fmt.Println("Failed to send activity. Status code:", resp.StatusCode)
	}
}

func TestSignActivity(t *testing.T) {
	// 获取公钥
	publicKey, err := getPublicKey("")
	if err != nil {
		fmt.Println("公钥获取失败：", err)
		return
	}
	// 构建签名字符串
	signatureData := "(request-target): post /activitypub/api/inbox/test11\nhost: testmakerone.shengjian.net\ndate: Wed, 17 May 2023 05:54:54 GMT\ndigest: SHA-256=Nr8ii7EV6Jg6is4Y+QZqPaK3COz8lRTIxOFSoM42c3Y=\ncontent-type: application/activity+json"

	signValue := "kW0n+1xIWBcW60uV6KUxFsJHO3BBF5DZcQUv70KJX6R4iWsKdjJqLluNgxaUKPjh33/1puE8Cg4GDnL5VcXp68VSpwdmcPoyaWRo5yAZXKIjC5LboI678+o2QsJHcm3+iP7jTqJJbp2Sj2LqQcfcA3FZB9Bd0U/35yXfaLOzmsm3dSEfvHr3JdgS8ZwlAXIj8/7+TYXLYUUkQ0XQodRZyHBj61spHsEz35wCIE7pWDM8N8l2qdFYN57u7tpr8+6kFIjINFoGhL+VQ7viIhoqy7rudVs5ozDcqF6/xRVCOk5Qvkd62aZb86vrm6H1AXGg5T9GTtPXGWAuiNsRxrjxFw=="
	// 验证签名
	if !verifySignature(publicKey, signatureData, signValue) {
		fmt.Println("签名验证失败")
		return
	}
}
