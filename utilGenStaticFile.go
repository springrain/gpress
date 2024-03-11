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
	"compress/gzip"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
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
	return err
}

// genStaticHtmlFile 生成全站静态文件
func genStaticHtmlFile() error {
	genStaticHtmlLock.Lock()
	defer genStaticHtmlLock.Unlock()
	ctx := context.Background()
	contents := make([]Content, 0)
	f_post := zorm.NewSelectFinder(tableContentName, "id,tag").Append(" WHERE status<2 order by sortNo desc")
	err := zorm.Query(ctx, f_post, &contents, nil)
	if err != nil {
		return err
	}

	tagsMap := make(map[string]bool, 0)

	//删除整个目录
	os.RemoveAll(staticHtmlDir)
	//生成首页index网页
	fileHash, err := writeStaticHtml(httpServerPath, staticHtmlDir, "")
	if fileHash == "" || err != nil {
		return err
	}
	//上一个分页
	prvePageFileHash := ""
	//生成文章的静态网页
	for i := 0; i < len(contents); i++ {
		tag := contents[i].Tag
		if tag != "" {
			tagsMap[tag] = true
		}
		postId := contents[i].Id
		postURL := httpServerPath + "post/" + postId
		fileHash, err := writeStaticHtml(postURL, staticHtmlDir+"post/"+postId+"/", "")
		if fileHash == "" || err != nil {
			continue
		}
		fileHash, err = writeStaticHtml(httpServerPath+"page/"+strconv.Itoa(i+1), staticHtmlDir+"page/"+strconv.Itoa(i+1)+"/", prvePageFileHash)
		if fileHash == "" || err != nil {
			continue
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
		categoryURL := httpServerPath + "category/" + categoryId
		//生成栏目首页index
		fileHash, err := writeStaticHtml(categoryURL, staticHtmlDir+"category/"+categoryId+"/", "")
		if fileHash == "" || err != nil {
			return err
		}
		for j := 0; j < len(contents); j++ {
			fileHash, err := writeStaticHtml(httpServerPath+"category/"+categoryId+"/page/"+strconv.Itoa(j+1), staticHtmlDir+"category/"+categoryId+"/page/"+strconv.Itoa(j+1)+"/", prvePageFileHash)
			if fileHash == "" || err != nil {
				continue
			}
			//如果hash完全一致,认为是最后一页
			prvePageFileHash = fileHash
		}
	}

	//生成tag的静态页
	for tag := range tagsMap {
		tagURL := httpServerPath + "tag/" + tag
		//生成栏目首页index
		fileHash, err := writeStaticHtml(tagURL, staticHtmlDir+"tag/"+tag+"/", "")
		if fileHash == "" || err != nil {
			return err
		}
		for j := 0; j < len(contents); j++ {
			fileHash, err := writeStaticHtml(httpServerPath+"tag/"+tag+"/page/"+strconv.Itoa(j+1), staticHtmlDir+"tag/"+tag+"/page/"+strconv.Itoa(j+1)+"/", prvePageFileHash)
			if fileHash == "" || err != nil {
				continue
			}
			//如果hash完全一致,认为是最后一页
			prvePageFileHash = fileHash
		}
	}

	// TODO 复制主题里的css,js,image 和公共的public文件夹到statichtml根目录

	return nil
}

// writeStaticHtml 写入静态html
func writeStaticHtml(httpurl string, filePath string, fileHash string) (string, error) {

	response, err := http.Get(httpurl)
	if err != nil {
		FuncLogError(err)
		return "", err
	}
	// 读取资源数据 body: []byte
	body, err := io.ReadAll(response.Body)
	// 关闭资源流
	response.Body.Close()
	if err != nil {
		FuncLogError(err)
		return "", err
	}
	//计算hash
	bytehex := sha3.Sum256(body)
	bodyHash := hex.EncodeToString(bytehex[:])
	if bodyHash == fileHash { //如果hash一致,不再生成文件
		return bodyHash, nil
	}
	// 写入文件
	os.MkdirAll(filePath, os.ModePerm)
	err = os.WriteFile(filePath+"index.html", body, os.ModePerm)
	if err != nil {
		return bodyHash, err
	}
	//压缩文件
	gzipFile, err := os.Create(filePath + "index.html.gz")
	if err != nil {
		return bodyHash, err
	}

	gzipWrite := gzip.NewWriter(gzipFile)
	gzipWrite.Name = "index.html"
	//gzipWrite.Name = "index.html"
	_, err = gzipWrite.Write(body)

	defer func() {
		gzipFile.Close()
		gzipWrite.Close()
	}()
	return bodyHash, err
}
