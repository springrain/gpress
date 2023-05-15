package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
)

// activityJSONRender JSON contains the given interface object.
type activityJSONRender struct {
	data        interface{}
	contentType string
}

var activityJsonContentType = "application/activity+json; charset=utf-8"

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
		writeContentType(resp, activityJsonContentType)
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
			"https://" + domain + "/@" + username,
			"https://" + domain + "/users/" + username,
		},
		"links": []map[string]interface{}{
			{
				"rel":  "http://webfinger.net/rel/profile-page",
				"type": "text/html",
				"href": "https://" + domain + "/@" + username,
			},
			{
				"rel":  "self",
				"type": "application/activity+json",
				"href": "https://" + domain + "/users/" + username,
			},
			{
				"rel":      "http://ostatus.org/schema/1.0/subscribe",
				"template": "https://" + domain + "/authorize_interaction?uri={uri}",
			},
		},
	}
	c.Render(http.StatusOK, activityJSONRender{data: data, contentType: "application/jrd+json; charset=utf-8"})
}

func funcActivityPubUserInfo(ctx context.Context, c *app.RequestContext) {
	accept := string(c.GetHeader("Accept"))
	fmt.Println(accept)
	data := funcActivityPubUserInfoJson()
	if accept == "application/activity+json" { //json类型
		c.Render(http.StatusOK, activityJSONRender{data: data})
		c.Abort() // 终止后续调用
		return
	}
	//返回页面
	c.HTML(http.StatusOK, "activitypub/user.html", data)
}
func funcActivityPubUserInfoJson() map[string]interface{} {
	// 构造 activityPubUser JSON 对象
	data := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
		},
		"id":                "https://lawrenceli.me/api/activitypub/actor",
		"type":              "Person",
		"name":              "Lawrence Li",
		"preferredUsername": "lawrence",
		"summary":           "Blog",
		"inbox":             "https://lawrenceli.me/api/activitypub/inbox",
		"outbox":            "https://lawrenceli.me/api/activitypub/outbox",
		"followers":         "https://lawrenceli.me/api/activitypub/followers",
		"icon": map[string]string{
			"type":      "Image",
			"mediaType": "image/png",
			"url":       "https://lawrenceli.me/icon.png",
		},
		"publicKey": map[string]string{
			"id":           "https://lawrenceli.me/api/activitypub/actor#main-key",
			"owner":        "https://lawrenceli.me/api/activitypub/actor",
			"publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0RHqCKo3Zl+ZmwsyJUFe\niUBYdiWQe6C3W+d89DEzAEtigH8bI5lDWW0Q7rT60eppaSnoN3ykaWFFOrtUiVJT\nNqyMBz3aPbs6BpAE5lId9aPu6s9MFyZrK5QtuWfAGwv9VZPwUHrEJCFiY1G5IgK/\n+ZErSKYUTUYw2xSAZnLkalMFTRmLbmj8SlWp/5fryQd4jyRX/tBlsyFs/qvuwBtw\nuGSkWgTIMAYV71Wny9ns+Nwr4HYfF5eo2zInpwIYTCEbil79HcikUUTTO/vMMoqx\n46IiHcMj0SPlzDXxelZgqm0ojK2Z7BGudjvwSbWq/GtLoaXHeMUVpcOCtpyvtLr2\nYwIDAQAB\n-----END PUBLIC KEY-----",
		},
	}

	return data
}
