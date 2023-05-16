package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
)

// activityJSONRender JSON contains the given interface object.
type activityJSONRender struct {
	data        interface{}
	contentType string
}

const activityPubAccept = "application/activity+json"

var activityPubContentType = "application/activity+json; charset=utf-8"

// Render (JSON) writes data with custom ContentType.
func (r activityJSONRender) Render(resp *protocol.Response) error {
	r.WriteContentType(resp)
	jsonBytes, err := json.Marshal(r.data)
	if err != nil {
		return err
	}
	resp.AppendBody(jsonBytes)
	return nil
}

// WriteContentType (JSON) writes JSON ContentType.
func (r activityJSONRender) WriteContentType(resp *protocol.Response) {
	if r.contentType == "" {
		writeContentType(resp, activityPubContentType)
	} else {
		writeContentType(resp, r.contentType)
	}
}

func writeContentType(resp *protocol.Response, value string) {
	resp.Header.SetContentType(value)
}

func funcWebFinger(ctx context.Context, c *app.RequestContext) {
	resource := c.Query("resource")
	parts := strings.Split(resource, ":")
	if len(parts) != 2 || parts[0] != "acct" {
		c.Abort() // 终止后续调用
		return
	}
	acctParts := strings.Split(parts[1], "@")
	if len(acctParts) != 2 {
		// resource 参数格式不正确，返回错误响应
		c.Abort() // 终止后续调用
		return
	}
	username := acctParts[0]
	domain := acctParts[1]

	// todo
	// 需要查一下数据库,是否有这个用户的数据

	//如果没有 就返回nil

	// 构造 WebFinger JSON 对象
	data := map[string]interface{}{
		"subject": "acct:" + username + "@" + domain,
		"aliases": []string{
			"https://" + domain + "/acititypub/user/@" + username,
			"https://" + domain + "/acititypub/api/user/" + username,
		},
		"links": []map[string]interface{}{
			{
				"rel":  "http://webfinger.net/rel/profile-page",
				"type": "text/html",
				"href": "https://" + domain + "/acititypub/user/@" + username,
			},
			{
				"rel":  "self",
				"type": "application/activity+json",
				"href": "https://" + domain + "/acititypub/api/user/" + username,
			},
			{
				"rel":      "http://ostatus.org/schema/1.0/subscribe",
				"template": "https://" + domain + "/acititypub/authorize_interaction?uri={uri}",
			},
		},
	}
	c.Render(http.StatusOK, activityJSONRender{data: data, contentType: "application/jrd+json; charset=utf-8"})
}

func funcActivityPubUsers(ctx context.Context, c *app.RequestContext) {
	accept := string(c.GetHeader("Accept"))
	host := string(c.Host())
	userName := c.Param("userName")
	// 构造 activityPubUser JSON 对象
	data := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
		},
		"id":                "https://" + host + "/acititypub/api/user/" + userName,
		"type":              "Person",
		"name":              userName,
		"preferredUsername": userName,
		"summary":           "Blog",
		"inbox":             "https://" + host + "/acititypub/api/inbox/" + userName,
		"outbox":            "https://" + host + "/acititypub/api/outbox/" + userName,
		"followers":         "https://" + host + "/acititypub/api/followers/" + userName,
		"icon": map[string]string{
			"type":      "Image",
			"mediaType": "image/png",
			"url":       "https://" + host + "/acititypub/images/" + userName + "/icon.png",
		},
		/*
			"publicKey": map[string]string{
				"id":            "https://" + host + "/acititypub/api/user/" + userName + "/actor#main-key",
				"owner":       "https://" + host + "/acititypub/api/user/" + userName + "/actor",
				"publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0RHqCKo3Zl+ZmwsyJUFe\niUBYdiWQe6C3W+d89DEzAEtigH8bI5lDWW0Q7rT60eppaSnoN3ykaWFFOrtUiVJT\nNqyMBz3aPbs6BpAE5lId9aPu6s9MFyZrK5QtuWfAGwv9VZPwUHrEJCFiY1G5IgK/\n+ZErSKYUTUYw2xSAZnLkalMFTRmLbmj8SlWp/5fryQd4jyRX/tBlsyFs/qvuwBtw\nuGSkWgTIMAYV71Wny9ns+Nwr4HYfF5eo2zInpwIYTCEbil79HcikUUTTO/vMMoqx\n46IiHcMj0SPlzDXxelZgqm0ojK2Z7BGudjvwSbWq/GtLoaXHeMUVpcOCtpyvtLr2\nYwIDAQAB\n-----END PUBLIC KEY-----",
			},
		*/
	}

	if strings.Contains(accept, activityPubAccept) { //json类型
		c.Render(http.StatusOK, activityJSONRender{data: data})
		c.Abort() // 终止后续调用
		return
	}
	//返回页面
	c.HTML(http.StatusOK, "activitypub/user.html", data)
}

func funcActivityPubOutBox(ctx context.Context, c *app.RequestContext) {
	accept := string(c.GetHeader("Accept"))
	host := string(c.Host())
	userName := c.Param("userName")
	// 构造 activityPubUser JSON 对象
	data := map[string]interface{}{

		"@context":   "https://www.w3.org/ns/activitystreams",
		"id":         "https://" + host + "/acititypub/api/outbox/" + userName,
		"summary":    "一个简单的测试",
		"type":       "OrderedCollection",
		"totalItems": 100,
		"orderedItems": []map[string]interface{}{
			{
				"@context":     "https://www.w3.org/ns/activitystreams",
				"id":           "https://" + host + "/post/78-k8snodocker",
				"type":         "Note",
				"published":    time.Now(),
				"attributedTo": "https://" + host + "/acititypub/api/user/" + userName,
				"content":      "<a href=\"https://" + host + "/post/78-k8snodocker\">K8S不使用Docker</a>",
				"url":          "https://" + host + "/post/78-k8snodocker",
				"to":           []string{"https://www.w3.org/ns/activitystreams#Public"},
				//"cc":           []string{"https://" + host + "/acititypub/api/followers/" + userName},
			}, {
				"@context":     "https://www.w3.org/ns/activitystreams",
				"id":           "https://" + host + "/post/77-nftonxuperchain",
				"type":         "Note",
				"published":    time.Now(),
				"attributedTo": "https://" + host + "/acititypub/api/user/" + userName,
				"content":      "<a href=\"https://" + host + "/post/77-nftonxuperchain\">百度开放网络发行数字藏品</a>",
				"url":          "https://" + host + "/post/77-nftonxuperchain",
				"to":           []string{"https://www.w3.org/ns/activitystreams#Public"},
				//"cc":           []string{"https://" + host + "/acititypub/api/followers/" + userName},
			},
		},
	}
	if strings.Contains(accept, activityPubAccept) { //json类型
		c.Render(http.StatusOK, activityJSONRender{data: data})
		c.Abort() // 终止后续调用
		return
	}
	//返回页面
	c.HTML(http.StatusOK, "activitypub/outbox.html", data)
}
