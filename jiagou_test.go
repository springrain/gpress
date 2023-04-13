package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"gitee.com/chunanyong/zorm"
)

func TestNavMenu(t *testing.T) {
	deleteAll(context.Background(), tableNavMenuName)
	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05")

	navs := []string{"About", "Web", "BlockChain", "CloudNative"}

	for i := 0; i < len(navs); i++ {
		nav := navs[i]
		menu := zorm.NewEntityMap(tableNavMenuName)
		menu.Set("id", strings.ToLower(nav))
		menu.Set("menuName", nav)
		menu.Set("status", 1)
		menu.Set("sortNo", i+1)
		menu.Set("createTime", now)
		menu.Set("updateTime", now)

		_, err := saveEntityMap(context.Background(), menu)
		if err != nil {
			t.Error(err)
		}
	}

	siteMap := zorm.NewEntityMap(tableSiteName)
	siteMap.PkColumnName = "id"
	siteMap.Set("id", "gpress")
	siteMap.Set("title", "jiagou")
	siteMap.Set("name", "架构")
	siteMap.Set("domain", "jiagou.com")
	err := updateTable(context.Background(), siteMap)
	if err != nil {
		t.Error(err)
	}

}

func TestReadmks(t *testing.T) {
	deleteAll(context.Background(), "content")
	files, err := os.ReadDir("D:/post")
	if err != nil {
		t.Error("读取错误")
	}
	lists := make([]zorm.IEntityMap, 0)
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
		cMap := zorm.NewEntityMap(tableContentName)
		cMap.Set("id", id)
		cMap.Set("summary", summary)
		cMap.Set("markdown", markdown)
		cMap.Set("sortNo", sortNo)
		cMap.Set("status", 1)

		metaData, tocHtml, html, _ := conver2Html([]byte(markdown))
		dateStr := metaData["date"].(string)
		date, _ := time.Parse("2006-01-02T15:04:05+08:00", dateStr)
		categories := metaData["categories"].([]interface{})
		navMenu := slice2string(categories)
		tags := metaData["tags"].([]interface{})
		tag := slice2string(tags)
		cMap.Set("title", metaData["title"])
		cMap.Set("author", metaData["author"])
		cMap.Set("updateTime", date)
		cMap.Set("createTime", date)
		cMap.Set("navMenuName", navMenu)
		cMap.Set("navMenuID", navMenu)
		cMap.Set("tag", tag)

		cMap.Set("content", html)
		cMap.Set("toc", tocHtml)
		lists = append(lists, cMap)
		//saveNewTable(context.Background(), "content", cMap)

	}

	fmt.Println("------------")

	var temp zorm.IEntityMap            // 定义临时变量,进行数据交换
	for j := 0; j < len(lists)-1; j++ { // 外循环 循环次数
		for i := 0; i < len(lists)-1; i++ { // 内循环 数组遍历
			m1 := lists[i]
			m2 := lists[i+1]
			d1 := m1.GetDBFieldMap()["updateTime"].(time.Time)
			d2 := m2.GetDBFieldMap()["updateTime"].(time.Time)
			if d1.After(d2) {
				temp = lists[i]
				lists[i] = lists[i+1]
				lists[i+1] = temp
			}
		}
	}

	for i := 0; i < len(lists); i++ { // 内循环 数组遍历
		cMap := lists[i]
		cMap.Set("sortNo", i+1)
		date := cMap.GetDBFieldMap()["updateTime"].(time.Time)
		dateStr := date.Format("2006-01-02 15:04:05")
		cMap.Set("updateTime", dateStr)
		cMap.Set("createTime", dateStr)
		_, err := saveEntityMap(context.Background(), cMap)
		if err != nil {
			t.Error(err)
		}
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
