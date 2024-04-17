// Copyright (c) 2023 gpress Authors.
//
// This file is part of gpress.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"gitee.com/chunanyong/zorm"
	"golang.org/x/crypto/sha3"
)

// onlyOnce控制并发
// var onlyOnce = make(chan struct{}, 1)
var searchDataLock = &sync.Mutex{}
var genStaticHtmlLock = &sync.Mutex{}

// genSearchDataJson 生成flexSearch需要的json文件,默认2000条数据
func genSearchDataJson() error {
	//onlyOnce <- struct{}{}
	//defer func() { <-onlyOnce }()
	searchDataLock.Lock()
	defer searchDataLock.Unlock()

	finder := zorm.NewSelectFinder(tableContentName, "id,title,hrefURL,summary,createTime,tag,categoryName,content,description").Append("WHERE status=1 order by sortNo desc")
	page := zorm.NewPage()
	page.PageSize = 2000
	datas := make([]Content, 0)
	err := zorm.Query(context.Background(), finder, &datas, page)
	if err != nil {
		return err
	}
	for i := 0; i < len(datas); i++ {
		if datas[i].HrefURL == "" {
			datas[i].HrefURL = funcBasePath() + "post/" + datas[i].Id
		}
	}
	dataBytes, err := json.Marshal(datas)
	if err != nil {
		return err
	}
	err = os.WriteFile(searchDataJsonFile, dataBytes, os.ModePerm)
	if err != nil {
		return err
	}

	//压缩文件
	err = doGzipFile(searchDataJsonFile+compressedFileSuffix, bytes.NewReader(dataBytes))

	return err
}

// genStaticFile 生成全站静态文件和gzip文件,包括静态的html和search-data.json
func genStaticFile() error {
	genStaticHtmlLock.Lock()
	defer genStaticHtmlLock.Unlock()
	ctx := context.Background()
	contents := make([]Content, 0)
	domain := ""
	if site.Domain != "" {
		if strings.HasPrefix(site.Domain, "http://") || strings.HasPrefix(site.Domain, "https://") {
			domain = site.Domain
		} else { //默认使用https协议
			domain = "https://" + site.Domain
		}
	}
	f_post := zorm.NewSelectFinder(tableContentName, "id,tag").Append(" WHERE status<2 order by sortNo desc")
	err := zorm.Query(ctx, f_post, &contents, nil)
	if err != nil {
		return err
	}

	tagsMap := make(map[string]bool, 0)

	//删除整个目录
	os.RemoveAll(staticHtmlDir)
	//生成首页index网页
	fileHash, _, err := writeStaticHtml("", "")
	if fileHash == "" || err != nil {
		return err
	}
	//创建sitemap.xml
	sitemapFile, err := os.OpenFile(staticHtmlDir+"sitemap.xml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer sitemapFile.Close() // 确保在函数结束时关闭文件
	sitemapFile.WriteString(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	sitemapFile.WriteString("<url><loc>" + domain + "/index.html" + "</loc></url>")
	//上一个分页
	prvePageFileHash := ""
	//生成文章的静态网页
	for i := 0; i < len(contents); i++ {
		tag := contents[i].Tag
		if tag != "" {
			tagsMap[tag] = true
		}
		postId := contents[i].Id
		//postURL := httpServerPath + "post/" + postId
		fileHash, success, err := writeStaticHtml("post/"+postId, "")
		if fileHash == "" || err != nil {
			continue
		}
		if success {
			sitemapFile.WriteString("<url><loc>" + domain + "/post/" + postId + "/index.html</loc></url>")
		}

		fileHash, success, err = writeStaticHtml("page/"+strconv.Itoa(i+1), prvePageFileHash)
		if fileHash == "" || err != nil {
			continue
		}
		if success {
			sitemapFile.WriteString("<url><loc>" + domain + "/page/" + strconv.Itoa(i+1) + "/index.html</loc></url>")
		}
		//如果hash完全一致,认为是最后一页
		prvePageFileHash = fileHash
	}
	//生成栏目的静态网页
	categoryIds := make([]string, 0)
	f_category := zorm.NewSelectFinder(tableCategoryName, "id").Append(" WHERE status<2 order by sortNo desc")
	err = zorm.Query(ctx, f_category, &categoryIds, nil)
	if err != nil {
		return err
	}
	for i := 0; i < len(categoryIds); i++ {
		categoryId := categoryIds[i]
		//生成栏目首页index
		fileHash, success, err := writeStaticHtml("category/"+categoryId, "")
		if fileHash == "" || err != nil {
			return err
		}
		if success {
			sitemapFile.WriteString("<url><loc>" + domain + "/category/" + categoryId + "/index.html</loc></url>")
		}
		for j := 0; j < len(contents); j++ {
			fileHash, success, err := writeStaticHtml("category/"+categoryId+"/page/"+strconv.Itoa(j+1), prvePageFileHash)
			if fileHash == "" || err != nil {
				continue
			}
			if success {
				sitemapFile.WriteString("<url><loc>" + domain + "/category/" + categoryId + "/page/" + strconv.Itoa(j+1) + "/index.html</loc></url>")
			}
			//如果hash完全一致,认为是最后一页
			prvePageFileHash = fileHash
		}
	}

	//生成tag的静态页
	for tag := range tagsMap {
		//生成栏目首页index
		fileHash, success, err := writeStaticHtml("tag/"+tag, "")
		if fileHash == "" || err != nil {
			return err
		}
		if success {
			sitemapFile.WriteString("<url><loc>" + domain + "/tag/" + tag + "/index.html</loc></url>")
		}
		for j := 0; j < len(contents); j++ {
			fileHash, success, err := writeStaticHtml("tag/"+tag+"/page/"+strconv.Itoa(j+1), prvePageFileHash)
			if fileHash == "" || err != nil {
				continue
			}
			if success {
				sitemapFile.WriteString("<url><loc>" + domain + "/tag/" + tag + "/page/" + strconv.Itoa(j+1) + "/index.html</loc></url>")
			}
			//如果hash完全一致,认为是最后一页
			prvePageFileHash = fileHash
		}
	}
	//结束写入sitemap文件
	sitemapFile.WriteString("</urlset>")

	//遍历当前使用的模板文件夹,压缩文本格式的文件
	err = filepath.Walk(templateDir+"theme/"+site.Theme+"/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 分隔符统一为 / 斜杠
		path = filepath.ToSlash(path)

		// 只处理 js 和 css 文件夹
		if !(strings.Contains(path, "/js/") || strings.Contains(path, "/css/")) {
			return nil
		}

		//获取文件后缀
		suffix := filepath.Ext(path)

		// 压缩 js,mjs,json,css,html
		// 压缩字体文件 ttf,otf,svg  gzip_types font/ttf font/otf image/svg+xml
		if !(suffix == ".js" || suffix == ".mjs" || suffix == ".json" || suffix == ".css" || suffix == ".html" || suffix == ".ttf" || suffix == ".otf" || suffix == ".svg") {
			return nil
		}

		// 获取要打包的文件信息
		readFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer readFile.Close()
		reader := bufio.NewReader(readFile)
		//压缩文件
		err = doGzipFile(path+compressedFileSuffix, reader)

		return err
	})
	if err != nil {
		FuncLogError(err)
		return err
	}
	// 重新生成 search-data.json
	err = genSearchDataJson()
	if err != nil {
		FuncLogError(err)
	}
	// TODO 复制主题里的css,js,image 和公共的public文件夹到statichtml根目录

	return err
}

// writeStaticHtml 写入静态html
func writeStaticHtml(urlFilePath string, fileHash string) (string, bool, error) {
	httpurl := httpServerPath + urlFilePath
	filePath := staticHtmlDir + urlFilePath
	if urlFilePath != "" {
		filePath = filePath + "/"
	}
	response, err := http.Get(httpurl)
	if err != nil {
		return "", false, err
	}
	// 读取资源数据 body: []byte
	body, err := io.ReadAll(response.Body)
	// 关闭资源流
	response.Body.Close()
	if err != nil {
		return "", false, err
	}
	//计算hash
	bytehex := sha3.Sum256(body)
	bodyHash := hex.EncodeToString(bytehex[:])
	if bodyHash == fileHash { //如果hash一致,不再生成文件
		return bodyHash, false, nil
	}
	// 写入文件
	os.MkdirAll(filePath, os.ModePerm)
	err = os.WriteFile(filePath+"index.html", body, os.ModePerm)
	if err != nil {
		return bodyHash, false, err
	}
	// 压缩gzip文件
	err = doGzipFile(filePath+"index.html"+compressedFileSuffix, bytes.NewReader(body))
	if err != nil {
		return bodyHash, false, err
	}
	return bodyHash, true, nil
}

// doGzipFile 压缩gzip文件
func doGzipFile(gzipFilePath string, reader io.Reader) error {

	//如果文件存在就删除
	if pathExist(gzipFilePath) {
		os.Remove(gzipFilePath)
	}
	//创建文件
	gzipFile, err := os.Create(gzipFilePath)
	if err != nil {
		return err
	}
	defer gzipFile.Close()

	gzipWrite, err := gzip.NewWriterLevel(gzipFile, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer gzipWrite.Close()
	_, err = io.Copy(gzipWrite, reader)
	return err
}
