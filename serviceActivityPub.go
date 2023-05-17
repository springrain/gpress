package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitee.com/chunanyong/zorm"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
)

// activityJSONRender JSON contains the given interface object.
type activityJSONRender struct {
	data        interface{}
	contentType string
}

const activityPubAccept = "application/activity+json"

const activityPubContentType = "application/activity+json; charset=utf-8"

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
			"https://" + domain + "/activitypub/user/@" + username,
			"https://" + domain + "/activitypub/api/user/" + username,
		},
		"links": []map[string]interface{}{
			{
				"rel":  "http://webfinger.net/rel/profile-page",
				"type": "text/html",
				"href": "https://" + domain + "/activitypub/user/@" + username,
			},
			{
				"rel":  "self",
				"type": "application/activity+json",
				"href": "https://" + domain + "/activitypub/api/user/" + username,
			},
			{
				"rel":      "http://ostatus.org/schema/1.0/subscribe",
				"template": "https://" + domain + "/activitypub/authorize_interaction?uri={uri}",
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
		"id":                "https://" + host + "/activitypub/api/user/" + userName,
		"type":              "Person",
		"name":              userName,
		"preferredUsername": userName,
		"summary":           "Blog",
		"inbox":             "https://" + host + "/activitypub/api/inbox/" + userName,
		"outbox":            "https://" + host + "/activitypub/api/outbox/" + userName,
		"following":         "https://" + host + "/activitypub/api/following/" + userName,
		"followers":         "https://" + host + "/activitypub/api/followers/" + userName,
		"icon": map[string]string{
			"type":      "Image",
			"mediaType": "image/png",
			"url":       "https://" + host + "/activitypub/images/" + userName + "/icon.png",
		},
		"publicKey": map[string]string{
			"id":    "https://" + host + "/activitypub/api/user/" + userName + "#main-key",
			"owner": "https://" + host + "/activitypub/api/user/" + userName,
			"publicKeyPem": `-----BEGIN PUBLIC KEY-----
			MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAot4y1T8UffW+nQwYnAhh
			fIRVaTCf92FeAOtPQ+S1/bVxAlhE9O+17Qd3C9mLOImVPq55HAV0MHzW/eByIB2B
			FDzOfeiq/arxsaCziwEL9GDOF6PiHVsD/a8kGjG0a8RiwUv/ek0n5XzA+nTIXNVZ
			bVWRikRYDHiXZYeX78ex5d2gSvuKUuQMcsMgsFYBHTVP/kL/tv5vsi1Pf5sWkaQM
			p0kiQH1Nph/vBN8Wmhl2qsjSqO3Zp7otcFQSn6L8Dvmx1dIWhpgffgagxfztTje5
			QSg6TSdRJhsBJQboMAvvlzzSsM6QdomBDB//0kiRyakPeZasNf/BkFMm+gkHHc15
			wQIDAQAB
			-----END PUBLIC KEY-----`,
		},
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
	page := zorm.NewPage()
	page.TotalCount = 2
	// 构造 activityPubUser JSON 对象
	data := map[string]interface{}{

		"@context":   "https://www.w3.org/ns/activitystreams",
		"id":         "https://" + host + "/activitypub/api/outbox/" + userName,
		"summary":    "一个简单的测试",
		"type":       "OrderedCollection",
		"totalItems": page.TotalCount,
		"first":      "https://" + host + "/activitypub/api/outbox_page/" + userName + "/1",
		"last":       "https://" + host + "/activitypub/api/outbox_page/" + userName + "/" + strconv.Itoa(page.PageCount),
	}
	if strings.Contains(accept, activityPubAccept) { //json类型
		c.Render(http.StatusOK, activityJSONRender{data: data})
		c.Abort() // 终止后续调用
		return
	}
	//返回页面
	c.HTML(http.StatusOK, "activitypub/outbox.html", data)
}

func funcActivityPubOutBoxPage(ctx context.Context, c *app.RequestContext) {
	accept := string(c.GetHeader("Accept"))
	host := string(c.Host())
	userName := c.Param("userName")
	pageNo := c.GetInt("pageNo")
	if pageNo < 1 {
		pageNo = 1
	}
	if pageNo >= 2 {
		//c.Render(http.StatusOK, activityJSONRender{data: data})
		c.Abort() // 终止后续调用
	}
	// 构造 activityPubUser JSON 对象
	data := map[string]interface{}{

		"@context":   "https://www.w3.org/ns/activitystreams",
		"id":         "https://" + host + "/activitypub/api/outbox_page/" + userName + "/" + strconv.Itoa(pageNo),
		"summary":    "一个简单的测试",
		"type":       "OrderedCollectionPage",
		"totalItems": 2,
		"prev":       "https://" + host + "/activitypub/api/outbox_page/" + userName + "/" + strconv.Itoa(pageNo-1),
		"next":       "https://" + host + "/activitypub/api/outbox_page/" + userName + "/" + strconv.Itoa(pageNo+1),
		"partOf":     "https://" + host + "/activitypub/api/outbox/" + userName,
		"orderedItems": []map[string]interface{}{
			{
				"@context":  "https://www.w3.org/ns/activitystreams",
				"type":      "Create",
				"id":        "https://" + host + "/post/78-k8snodocker",
				"actor":     "https://" + host + "/activitypub/api/user/" + userName,
				"published": time.Now(),
				"to":        []string{"https://www.w3.org/ns/activitystreams#Public"},
				"object": map[string]interface{}{
					"@context":     "https://www.w3.org/ns/activitystreams",
					"id":           "https://" + host + "/post/78-k8snodocker",
					"type":         "Note",
					"published":    time.Now(),
					"attributedTo": "https://" + host + "/activitypub/api/user/" + userName,
					"content":      "<a href=\"https://" + host + "/post/78-k8snodocker\">K8S不使用Docker</a>",
					"url":          "https://" + host + "/post/78-k8snodocker",
					"to":           []string{"https://www.w3.org/ns/activitystreams#Public"},
					//"cc":           []string{"https://" + host + "/activitypub/api/followers/" + userName},
				},
			}, {
				"@context":  "https://www.w3.org/ns/activitystreams",
				"type":      "Create",
				"id":        "https://" + host + "/post/77-nftonxuperchain",
				"actor":     "https://" + host + "/activitypub/api/user/" + userName,
				"published": time.Now(),
				"to":        []string{"https://www.w3.org/ns/activitystreams#Public"},
				"object": map[string]interface{}{
					"@context":     "https://www.w3.org/ns/activitystreams",
					"id":           "https://" + host + "/post/77-nftonxuperchain",
					"type":         "Note",
					"published":    time.Now(),
					"attributedTo": "https://" + host + "/activitypub/api/user/" + userName,
					"content":      "<a href=\"https://" + host + "/post/77-nftonxuperchain\">百度开放网络发行数字藏品</a>",
					"url":          "https://" + host + "/post/77-nftonxuperchain",
					"to":           []string{"https://www.w3.org/ns/activitystreams#Public"},
					//"cc":           []string{"https://" + host + "/activitypub/api/followers/" + userName},
				},
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

func funcActivityPubInBox(ctx context.Context, c *app.RequestContext) {
	bodyByte, _ := c.Body()
	body := make(map[string]interface{}, 0)
	json.Unmarshal(bodyByte, &body)
	c.Render(http.StatusOK, activityJSONRender{data: "success"})
	go funcSendAcceptMessage(body["id"].(string))

}

// inbox交互是通过事件异步返回给对方的inbox
func funcSendAcceptMessage(id string) {
	fmt.Println(id)
}

// activitySignatureHandler 验签拦截器
func activitySignatureHandler(ctx context.Context, c *app.RequestContext) {
	//if !strings.Contains(string(c.URI().Path()), "/test11") {
	//	return
	//}

	bodyByte, _ := c.Body()
	hash := sha256.Sum256(bodyByte)
	digest := "SHA-256=" + base64.StdEncoding.EncodeToString(hash[:])

	if digest != string(c.GetHeader("Digest")) {
		c.Abort() // 终止后续调用
		FuncLogError(errors.New("内容签名解析失败,digest=" + digest + ",Digest=" + string(c.GetHeader("Digest"))))
		return
	}

	// 从请求头部获取签名数据
	signatureString := string(c.GetHeader("Signature"))

	// 解析签名数据，获取相关信息
	signature, err := parseSignature(signatureString)
	if err != nil {
		c.Abort() // 终止后续调用
		FuncLogError(fmt.Errorf("签名解析失败：%w", err))
		return
	}

	// 获取公钥
	publicKey, err := getPublicKey(signature.KeyID)
	if err != nil {
		c.Abort() // 终止后续调用
		FuncLogError(fmt.Errorf("公钥获取失败：%w", err))
		return
	}
	// 构建签名字符串
	signatureData := buildSignatureData(c, signature.Headers)

	// 验证签名
	if !verifySignature(publicKey, signatureData, signature.Value) {
		c.Abort() // 终止后续调用
		FuncLogError(errors.New("签名验证失败"))
		return
	}

}
