package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNavMenu(t *testing.T) {
	deleteAll(context.Background(), indexNavMenuName)
	// 获取当前时间
	now := time.Now()

	navs := []string{"About", "Web", "BlockChain", "CloudNative"}

	for i := 0; i < len(navs); i++ {
		nav := navs[i]
		menu := make(map[string]interface{}, 0)
		menu["id"] = strings.ToLower(nav)
		menu["menuName"] = nav
		menu["sortNo"] = i + 1
		menu["createTime"] = now
		menu["updateTime"] = now
		saveNewIndex(context.Background(), indexNavMenuName, menu)
	}

	siteMap := make(map[string]interface{}, 0)
	siteMap["id"] = "gpress"
	siteMap["title"] = "jiagou"
	siteMap["name"] = "架构"
	siteMap["domain"] = "jiagou.com"
	updateIndex(context.Background(), "site", "gpress", siteMap)

}

func TestReadmks(t *testing.T) {
	deleteAll(context.Background(), "content")
	files, err := os.ReadDir("D:/post")
	if err != nil {
		t.Error("读取错误")
	}
	lists := make([]map[string]interface{}, 0)
	for i, file := range files {
		fileName := file.Name()
		source, err := os.ReadFile("D:/post/" + fileName)
		if err != nil {
			continue
		}
		sortNo := i + 1
		id := fileName[:strings.LastIndex(fileName, ".")]

		fmt.Println(id)
		markdown := strings.TrimSpace(string(source))
		smk := markdown[strings.Index(markdown, "---")+3:]
		start := strings.Index(smk, "---") + 3
		smk = strings.TrimSpace(smk[start:])
		smkRune := []rune(smk)
		end := 100
		if end > len(smkRune) {
			end = len(smkRune)
		}
		summary := string(smkRune[0:end])
		summary = strings.TrimSpace(summary)
		cMap := make(map[string]interface{}, 0)
		cMap["id"] = id
		cMap["summary"] = summary
		cMap["markdown"] = markdown
		cMap["sortNo"] = sortNo

		metaData, tocHtml, html, _ := conver2Html([]byte(markdown))
		dateStr := metaData["date"].(string)
		date, _ := time.Parse("2006-01-02T15:04:05+08:00", dateStr)
		categories := metaData["categories"].([]interface{})
		navMenu := slice2string(categories)
		tags := metaData["tags"].([]interface{})
		tag := slice2string(tags)
		cMap["title"] = metaData["title"]
		cMap["author"] = metaData["author"]
		cMap["updateTime"] = date
		cMap["createTime"] = date
		cMap["navMenuName"] = navMenu
		cMap["navMenuID"] = navMenu
		cMap["tag"] = tag

		cMap["content"] = html
		cMap["toc"] = tocHtml
		lists = append(lists, cMap)
		//saveNewIndex(context.Background(), "content", cMap)

	}

	fmt.Println("------------")

	var temp map[string]interface{}     // 定义临时变量,进行数据交换
	for j := 0; j < len(lists)-1; j++ { // 外循环 循环次数
		for i := 0; i < len(lists)-1; i++ { // 内循环 数组遍历
			m1 := lists[i]
			m2 := lists[i+1]
			d1 := m1["updateTime"].(time.Time)
			d2 := m2["updateTime"].(time.Time)
			if d1.After(d2) {
				temp = lists[i]
				lists[i] = lists[i+1]
				lists[i+1] = temp
			}
		}
	}

	for i := 0; i < len(lists); i++ { // 内循环 数组遍历
		cMap := lists[i]
		cMap["sortNo"] = i + 1
		saveNewIndex(context.Background(), "content", cMap)
	}

}
func slice2string(slice []interface{}) string {
	len := len(slice)
	s := ""
	for i := 0; i < len; i++ {
		str := slice[i].(string)
		s = s + str
		if i+1 < len {
			s = s + ","
		}
	}
	return s
}
