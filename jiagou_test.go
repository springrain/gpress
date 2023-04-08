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

}

func TestReadmks(t *testing.T) {
	deleteAll(context.Background(), "content")
	files, err := os.ReadDir("D:/post")
	if err != nil {
		t.Error("读取错误")
	}
	for i, file := range files {
		fileName := file.Name()
		source, err := os.ReadFile("D:/post/" + fileName)
		if err != nil {
			continue
		}
		id := fileName[:strings.LastIndex(fileName, ".")]
		fmt.Println(id)
		markdown := strings.TrimSpace(string(source))
		smk := markdown[strings.Index(markdown, "---")+3:]
		start := strings.Index(smk, "---") + 6
		end := start + 200
		if len(markdown) < end {
			end = len(markdown)
		}
		summary := markdown[start:end]
		summary = strings.TrimSpace(summary)
		cMap := make(map[string]interface{}, 0)
		cMap["id"] = id
		cMap["summary"] = summary
		cMap["markdown"] = markdown
		cMap["sortNo"] = i + 1

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

		saveNewIndex(context.Background(), "content", cMap)

	}

	fmt.Println("------------")

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
