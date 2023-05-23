package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

/**
https://lawrenceli.me/blog/activitypub
https://s3lph.me/activitypub-static-site.html
https://wangqiao.me/posts/activitypub-from-decentralized-to-distributed-social-networks/
https://blog.joinmastodon.org/2018/06/how-to-implement-a-basic-activitypub-server/
https://tinysubversions.com/notes/reading-activitypub/#a-note-on-json-ld
https://emptystack.top/activitypub-for-static-blog/

https://docs.joinmastodon.org/
https://www.w3.org/TR/activitypub/
https://www.w3.org/TR/activitystreams-vocabulary/

**/
/*
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
*/
func TestSendActivity(t *testing.T) {
	// 构建Activity对象
	activity := map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
		"type":     "Create",
		"to": []string{
			"https://mastodon.social/users/9iuorg",
		},
		"actor": "https://" + activityPubDefaultDomain + "/activitypub/api/user/test11",
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

func TestPublicKey(t *testing.T) {
	// 获取公钥
	publicKey, err := getRSAPublicKeyPem("https://mastodon.social/users/9iuorg#main-key")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(publicKey)

}

func TestMakeSign(t *testing.T) {
	s, err := generateRSASignature("123")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(s)
}

func TestRFC2616(t *testing.T) {
	// mastodon 的日期格式为RFC2616 -- time.RFC1123
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	fmt.Println(date)
}
func TestCrateNote(t *testing.T) {
	//date := time.Now().Unix()
	id := "https://" + activityPubDefaultDomain + "/note/90"
	data := map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
		//操作自身的ID. 和actor同一个域名下的uri !!!
		"id":   id + "#Create",
		"txId": id + "#Create",
		"type": "Delete",
		//"type": "Delete", //不同的操作类型
		"to":    []string{"https://www.w3.org/ns/activitystreams#Public"},
		"actor": "https://" + activityPubDefaultDomain + "/activitypub/api/user/test11",
		"object": map[string]interface{}{
			// 业务Object 的ID,Update和Delete依据这个ID. 和actor同一个域名下的uri !!!
			"id":      id,
			"type":    "Note",
			"content": "发表一个Note",
		},
	}
	headerMap, _ := wrapRequestHeader("https://mastodon.social/inbox", "POST", data, "https://"+activityPubDefaultDomain+"/activitypub/api/user/test11#main-key", true)
	reponseMap, err := sendActivityPubRequest("https://mastodon.social/inbox", "POST", data, headerMap)
	if err != nil {
		t.Error(err)
	}
	body, _ := json.Marshal(reponseMap)
	//这里需要debug暂停一下,等待mastodon服务请求账户的公钥,一般很快
	fmt.Println(string(body))

}
func TestSend11(t *testing.T) {
	//date := time.Now().Unix()
	id := "https://" + activityPubDefaultDomain + "/post/88"
	data := map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
		//操作自身的ID. 和actor同一个域名下的uri !!!
		"id": id + "#Create",
		//"id":   date,
		"type": "Delete",
		//"type": "Delete", //不同的操作类型
		//"to":   []string{"https://mastodon.social/users/9iuorg"},
		"to":    []string{"https://www.w3.org/ns/activitystreams#Public"},
		"actor": "https://" + activityPubDefaultDomain + "/activitypub/api/user/test11",
		"object": map[string]interface{}{
			// 业务Object 的ID,Update和Delete依据这个ID. 和actor同一个域名下的uri !!!
			"id":   id,
			"type": "Note",
			//"url":     "https://"+activityPubDefaultDomain+"/post/78-k8snodocker",
			"content": "一个简单的测试",
			//"cc":      []string{"https://mastodon.social/users/9iuorg"},
		},
	}
	headerMap, _ := wrapRequestHeader("https://mastodon.social/inbox", "POST", data, "https://"+activityPubDefaultDomain+"/activitypub/api/user/test11#main-key", true)
	reponseMap, err := sendActivityPubRequest("https://mastodon.social/inbox", "POST", data, headerMap)
	//reponseMap, err := sendRequest("https://mastodon.social/users/9iuorg/inbox", "POST", data, true)
	if err != nil {
		t.Error(err)
	}
	body, _ := json.Marshal(reponseMap)
	//这里需要debug暂停一下,等待mastodon服务请求账户的公钥,一般很快
	fmt.Println(string(body))

}

// 测试发表评论
func TestSendReply(t *testing.T) {
	//date := time.Now().Unix()
	id := "https://" + activityPubDefaultDomain + "/reply/89"
	data := map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
		//操作自身的ID. 和actor同一个域名下的uri !!!
		"id":   id + "#Create",
		"type": "Create",
		//"type": "Delete", //不同的操作类型
		//"to":   []string{"https://mastodon.social/users/9iuorg"},
		"to":    []string{"https://www.w3.org/ns/activitystreams#Public"},
		"actor": "https://" + activityPubDefaultDomain + "/activitypub/api/user/test11",
		//转发关注
		"cc": []string{"https://mastodon.social/users/9iuorg/followers"},
		"object": map[string]interface{}{
			// 业务Object 的ID,Update和Delete依据这个ID. 和actor同一个域名下的uri !!!
			"id":   id,
			"type": "Note",
			//给谁回复,一般是资源URL
			"inReplyTo": "https://mastodon.social/@9iuorg/110354226839707746",
			//回复的内容
			"content": "简单的回复",
		},
	}

	headerMap, _ := wrapRequestHeader("https://mastodon.social/inbox", "POST", data, "https://"+activityPubDefaultDomain+"/activitypub/api/user/test11#main-key", true)
	reponseMap, err := sendActivityPubRequest("https://mastodon.social/inbox", "POST", data, headerMap)
	if err != nil {
		t.Error(err)
	}
	body, _ := json.Marshal(reponseMap)
	//这里需要debug暂停一下,等待mastodon服务请求账户的公钥,一般很快
	fmt.Println(string(body))

}
