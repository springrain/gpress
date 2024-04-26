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
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gitee.com/chunanyong/zorm"
)

var funcMap = template.FuncMap{

	"basePath":    funcBasePath,
	"addInt":      funcAddInt,
	"addFloat":    funcAddFloat,
	"convertType": funcConvertType,
	"safeHTML":    funcSafeHTML,
	"safeURL":     funcSafeURL,
	"hrefURL":     funcHrefURL,
	"relURL":      funcRelURL,
	"site":        funcSite,
	//"category":     funcCategory,
	"selectList": funcSelectList,
	"selectOne":  funcSelectOne,
	//"md5":      funcMD5,
	//"sass":       funcSass,
	//"themePath":  funcThemePath,
	//"themeFile":  funcThemeFile,
	//"convertType":      funcConvertJson,
	//"convertMap":       funcConvertMap,
	"hasPrefix":        hasPrefix,
	"hasSuffix":        hasSuffix,
	"contains":         contains,
	"generateStringID": FuncGenerateStringID,
	"treeCategory":     funcTreeCategory,
	"themeName":        funcThemeName,
	"themeTemplate":    funcThemeTemplate,
	"version":          funcVersion,
	"seq":              funcSeq,
}

// funcBasePath 基础路径,前端所有的资源请求必须带上 {{basePath}}
func funcBasePath() string {
	return config.BasePath
}

// funcConvertType 多语言i18n适配,例如 {{ T "nextPage" }}
func funcConvertType(value interface{}, sourceType string, targetType string) (interface{}, error) {
	// json字符串转成Map
	if sourceType == "json" && targetType == "object" {
		obj := make(map[string]interface{})
		jsonStr := value.(string)
		json.Unmarshal([]byte(jsonStr), &obj)
		return obj, nil
	} else if sourceType == "object" && targetType == "json" { //对象转成json字符串
		jsonData, err := json.Marshal(value)
		if err != nil {
			return "{}", nil
		}
		s := string(jsonData)
		return s, nil
	} else if sourceType == "string" && targetType == "int" { //字符串转int
		valueStr := value.(string)
		valueInt, _ := strconv.Atoi(valueStr)
		return valueInt, nil
	} else if sourceType == "int" && targetType == "string" { //int转字符串
		valueInt := value.(int)
		valueStr := strconv.Itoa(valueInt)
		return valueStr, nil
	}
	return nil, nil
}

// funcSafeHTML 转义html字符串
func funcSafeHTML(html string) (template.HTML, error) {
	ss := template.HTML(html)
	return ss, nil
}

// funcSafeURL 转义URL字符串
func funcSafeURL(html string) (template.URL, error) {
	ss := template.URL(html)
	return ss, nil
}

func funcHrefURL(href string) (string, error) {
	href = strings.TrimSpace(href)
	if strings.HasPrefix(href, "http") { //http协议开头
		return href, nil
	}
	href = strings.TrimPrefix(href, "/") //斜杠开头就删除掉
	return funcBasePath() + href, nil
}

// funcRelURL 拼接url路径的
func funcRelURL(url string) (template.HTML, error) {
	return funcSafeHTML(themePath + url)
}

// 站点信息
func funcSite() (Site, error) {
	finder := zorm.NewSelectFinder(tableSiteName).Append(" WHERE id=?", "gpress_site")
	site := Site{}
	_, err := zorm.QueryRow(context.Background(), finder, &site)
	return site, err
}

/*
// 导航信息
func funcCategory() ([]Category, error) {
	finder := zorm.NewSelectFinder(tableCategoryName)
	finder.Append(" WHERE status=1 order by sortNo desc")
	page := zorm.NewPage()
	page.PageSize = 200
	list := make([]Category, 0)
	err := zorm.Query(context.Background(), finder, &list, page)
	return list, err
}
*/

// funcThemeTemplate 主题模板,prefix 过滤文件前缀
func funcThemeTemplate(prefix string) ([]ThemeTemplate, error) {
	list := make([]ThemeTemplate, 0)

	matches, err := filepath.Glob(themeDir + site.Theme + "/" + prefix + "*.html")
	if err != nil {
		return list, err
	}

	// 遍历匹配到的文件路径,并打印
	for _, match := range matches {
		path := filepath.ToSlash(match)
		path = path[strings.Index(path, themeDir)+len(themeDir):]
		path = path[strings.Index(path, "/")+1:]
		themeTemplate := ThemeTemplate{}
		themeTemplate.Name = path
		themeTemplate.FilePath = path
		list = append(list, themeTemplate)
	}

	return list, err
}

/*
var analyzerMap = map[string]string{commaAnalyzerName: "逗号分词器", gseAnalyzerName: "默认分词器", keywordAnalyzerName: "不分词", numericAnalyzerName: "数字分词器", datetimeAnalyzerName: "日期分词器"}

	func funcAnalyzer() map[string]string {
		return analyzerMap
	}
*/

func funcAddInt(x, y int) int {
	return x + y
}
func funcAddFloat(x, y float64) float64 {
	return x + y
}

// 查询'order by'在sql中出现的开始位置和结束位置
// Query the start position and end position of'order by' in SQL
var (
	orderByExpr      = "(?i)\\s(order)\\s+by\\s"
	orderByRegexp, _ = regexp.Compile(orderByExpr)
)

// findOrderByIndex 查询order by在sql中出现的开始位置和结束位置
// findOrderByIndex Query the start position and end position of'order by' in SQL
func findOrderByIndex(strsql *string) []int {
	loc := orderByRegexp.FindStringIndex(*strsql)
	return loc
}
func funcSelectList(urlPathParam string, q string, pageNo int, pageSize int, sql string, values ...interface{}) (ResponseData, error) {
	responseData := ResponseData{StatusCode: 0}
	sql = strings.TrimSpace(sql)
	if sql == "" || strings.Contains(sql, ";") {
		err := errors.New("sql语句错误")
		responseData.ERR = err
		responseData.StatusCode = 0
		return responseData, err
	}

	finder := zorm.NewFinder().Append("SELECT")
	if q != "" { // 如果有搜索关键字
		whereSQL := strings.ToLower(sql)
		locOrderBy := findOrderByIndex(&sql)
		orderBy := ""
		if len(locOrderBy) > 0 {
			orderBy = sql[locOrderBy[0]:]
			sql = sql[:locOrderBy[0]]
		}

		i := strings.Index(whereSQL, " where ")
		// fst5 搜索相关性排序 ORDER BY rank; 后期再进行修改调整,先按照sortNo排序
		if i < 0 { // 没有where
			finder.Append(sql, values...)
			finder.Append(" where rowid in (select rowid from fts_content where fts_content match jieba_query(?) ) ", q)
		} else {
			finder.Append(sql[:i+7]+" rowid in (select rowid from fts_content where fts_content match jieba_query(?) ) and ", q)
			finder.Append(sql[i+7:], values...)
		}
		finder.Append(orderBy)
	} else {
		finder.Append(sql, values...)
	}

	//finder.Append("order by sortNo desc")
	page := zorm.NewPage()
	page.PageNo = pageNo
	if pageSize > 1000 {
		pageSize = 1000
	}
	page.PageSize = pageSize
	ctx := context.Background()
	switch urlPathParam {
	case tableConfigName:
		data := make([]Config, 0)
		zorm.Query(ctx, finder, &data, page)
		responseData.Data = data
	case tableUserName:
		data := make([]User, 0)
		zorm.Query(ctx, finder, &data, page)
		responseData.Data = data
	case tableSiteName:
		data := make([]Site, 0)
		zorm.Query(ctx, finder, &data, page)
		responseData.Data = data
	case tableCategoryName:
		page.PageSize = 100
		data := make([]Category, 0)
		zorm.Query(ctx, finder, &data, page)
		responseData.Data = data
	case tableContentName:
		data := make([]Content, 0)
		zorm.Query(ctx, finder, &data, page)
		responseData.Data = data
	case "": // 对象为空查询map
		responseData.Data, responseData.ERR = zorm.QueryMap(ctx, finder, page)
	default:
		err := errors.New(urlPathParam + "表不存在!")
		responseData.ERR = err
		responseData.StatusCode = 0
		return responseData, err
	}
	responseData.Page = page
	responseData.StatusCode = 1
	return responseData, nil
}

func funcSelectOne(urlPathParam string, sql string, values ...interface{}) (interface{}, error) {
	sql = strings.TrimSpace(sql)
	if sql == "" || strings.Contains(sql, ";") {
		err := errors.New("sql语句错误")
		return nil, err
	}
	var selectOneData interface{}
	finder := zorm.NewFinder().Append("SELECT")
	finder.Append(sql, values...)

	page := zorm.NewPage()
	page.PageSize = 1
	page.PageNo = 1
	ctx := context.Background()
	switch urlPathParam {
	case tableConfigName:
		data := make([]Config, 0)
		zorm.Query(ctx, finder, &data, page)
		if len(data) > 0 {
			selectOneData = data[0]
		} else {
			selectOneData = Config{}
		}
	case tableUserName:
		data := make([]User, 0)
		zorm.Query(ctx, finder, &data, page)
		if len(data) > 0 {
			selectOneData = data[0]
		} else {
			selectOneData = User{}
		}
	case tableSiteName:
		data := make([]Site, 0)
		zorm.Query(ctx, finder, &data, page)
		if len(data) > 0 {
			selectOneData = data[0]
		} else {
			selectOneData = Site{}
		}
	case tableCategoryName:
		data := make([]Category, 0)
		zorm.Query(ctx, finder, &data, page)
		if len(data) > 0 {
			selectOneData = data[0]
		} else {
			selectOneData = Category{}
		}
	case tableContentName:
		data := make([]Content, 0)
		zorm.Query(ctx, finder, &data, page)
		if len(data) > 0 {
			selectOneData = data[0]
		} else {
			selectOneData = Content{}
		}
	case "": // 对象为空查询map
		selectOneData, _ = zorm.QueryRowMap(ctx, finder)
	default:
		err := errors.New(urlPathParam + "表不存在!")
		return selectOneData, err
	}

	return selectOneData, nil
}

func funcTreeCategory(sql string, values ...interface{}) []Category {
	ctx := context.Background()
	categorys := make([]Category, 0)
	finder := zorm.NewFinder().Append("SELECT")
	finder.Append(sql, values...)
	err := zorm.Query(ctx, finder, &categorys, nil)
	if err != nil {
		return categorys
	}

	treeCategory := sliceCategory2Tree(categorys)

	return treeCategory
}

func funcThemeName() []string {
	themeNames := make([]string, 0)
	files, err := os.ReadDir(themeDir)
	if err != nil {
		return themeNames
	}
	for _, file := range files {
		if file.IsDir() {
			themeNames = append(themeNames, file.Name())
		}
	}
	return themeNames
}

/*
	func funcConvertJson(obj interface{}) (string, error) {
		// 将对象转换为 JSON 字符串
		jsonData, err := json.Marshal(obj)
		if err != nil {
			return "{}", nil
		}
		s := string(jsonData)
		return s, nil
	}

func funcConvertMap(jsonStr string) (map[string]interface{}, error) {

		obj := make(map[string]interface{})
		err := json.Unmarshal([]byte(jsonStr), &obj)
		if err != nil {
			return obj, nil
		}
		return obj, nil
	}
*/
func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func hasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func sliceCategory2Tree(categorys []Category) []Category {
	categorysMap := make(map[string]*Category)
	for i := 0; i < len(categorys); i++ {
		c1 := categorys[i]
		var category *Category
		for j := 0; j < len(categorys); j++ {
			if c1.Pid == categorys[j].Id {
				category = &categorys[j]
				break
			}
		}

		if category != nil { //找到了上级
			delete(categorysMap, c1.Id)
			if category.Leaf == nil {
				category.Leaf = make([]Category, 0)
			}
			category.Leaf = append(category.Leaf, c1)
			categorysMap[category.Id] = category
		} else {
			categorysMap[c1.Id] = &c1
		}

	}

	//重新排序获取
	treeCategory := make([]Category, 0)
	for i := 0; i < len(categorys); i++ {
		id := categorys[i].Id
		c, has := categorysMap[id]
		if has {
			treeCategory = append(treeCategory, *c)
		}
	}

	return treeCategory

}

func funcVersion() string {
	return version
}

// seq 函数用于生成一个数字序列，这里是一个模拟的实现
func funcSeq(start, end int) []int {
	nums := make([]int, 0)
	for i := start; i <= end; i++ {
		nums = append(nums, i)
	}
	return nums
}
