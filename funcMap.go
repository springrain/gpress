package main

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"regexp"
	"strings"

	"gitee.com/chunanyong/zorm"
)

var funcMap = template.FuncMap{

	"basePath":     funcBasePath,
	"addInt":       funcAddInt,
	"addFloat":     funcAddFloat,
	"T":            funcT,
	"safeHTML":     funcSafeHTML,
	"safeURL":      funcSafeURL,
	"hrefURL":      funcHrefURL,
	"relURL":       funcRelURL,
	"site":         funcSite,
	"category":     funcCategory,
	"pageTemplate": funcPageTemplate,
	"selectList":   funcSelectList,
	"selectOne":    funcSelectOne,
	//"md5":      funcMD5,
	//"sass":       funcSass,
	//"themePath":  funcThemePath,
	//"themeFile":  funcThemeFile,

}

// funcBasePath 基础路径,前端所有的资源请求必须带上 {{basePath}}
func funcBasePath() string {
	return config.BasePath
}

// funcT 多语言i18n适配,例如 {{ T "nextPage" }}
func funcT(key string) (string, error) {
	return key, nil
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
	finder := zorm.NewSelectFinder(tableSiteName).Append(" WHERE id=?", appName)
	site := Site{}
	_, err := zorm.QueryRow(context.Background(), finder, &site)
	return site, err
}

// 菜单信息
func funcCategory() ([]Category, error) {
	finder := zorm.NewSelectFinder(tableCategoryName)
	finder.Append(" WHERE status=1 order by sortNo desc")
	page := zorm.NewPage()
	page.PageSize = 200
	list := make([]Category, 0)
	err := zorm.Query(context.Background(), finder, &list, page)
	return list, err
}

// 页面模板
func funcPageTemplate() ([]PageTemplate, error) {
	finder := zorm.NewSelectFinder(tablePageTemplateName)
	finder.Append(" order by sortNo desc")
	page := zorm.NewPage()
	page.PageSize = 200

	list := make([]PageTemplate, 0)
	err := zorm.Query(context.Background(), finder, &list, page)
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
func funcSelectList(urlPathParam string, q string, pageNo int, sql string, values ...interface{}) (ResponseData, error) {
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
	switch urlPathParam {
	case tableConfigName:
		data := make([]Config, 0)
		zorm.Query(context.Background(), finder, &data, page)
		responseData.Data = data
	case tableUserName:
		data := make([]User, 0)
		zorm.Query(context.Background(), finder, &data, page)
		responseData.Data = data
	case tableSiteName:
		data := make([]Site, 0)
		zorm.Query(context.Background(), finder, &data, page)
		responseData.Data = data
	case tablePageTemplateName:
		data := make([]PageTemplate, 0)
		zorm.Query(context.Background(), finder, &data, page)
		responseData.Data = data
	case tableCategoryName:
		data := make([]Category, 0)
		zorm.Query(context.Background(), finder, &data, page)
		responseData.Data = data
	case tableContentName:
		data := make([]Content, 0)
		zorm.Query(context.Background(), finder, &data, page)
		responseData.Data = data
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

func funcSelectOne(urlPathParam string, sql string, values ...interface{}) (ResponseData, error) {
	responseData := ResponseData{StatusCode: 0}
	sql = strings.TrimSpace(sql)
	if sql == "" || strings.Contains(sql, ";") {
		err := errors.New("sql语句错误")
		responseData.ERR = err
		return responseData, err
	}
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
			responseData.Data = data[0]
		} else {
			responseData.Data = Config{}
		}
	case tableUserName:
		data := make([]User, 0)
		zorm.Query(ctx, finder, &data, page)
		if len(data) > 0 {
			responseData.Data = data[0]
		} else {
			responseData.Data = User{}
		}
	case tableSiteName:
		data := make([]Site, 0)
		zorm.Query(ctx, finder, &data, page)
		if len(data) > 0 {
			responseData.Data = data[0]
		} else {
			responseData.Data = Site{}
		}
	case tablePageTemplateName:
		data := make([]PageTemplate, 0)
		zorm.Query(ctx, finder, &data, page)
		if len(data) > 0 {
			responseData.Data = data[0]
		} else {
			responseData.Data = PageTemplate{}
		}
	case tableCategoryName:
		data := make([]Category, 0)
		zorm.Query(ctx, finder, &data, page)
		if len(data) > 0 {
			responseData.Data = data[0]
		} else {
			responseData.Data = Category{}
		}
	case tableContentName:
		data := make([]Content, 0)
		zorm.Query(ctx, finder, &data, page)
		if len(data) > 0 {
			responseData.Data = data[0]
		} else {
			responseData.Data = Content{}
		}
	default:
		err := errors.New(urlPathParam + "表不存在!")
		responseData.ERR = err
		responseData.StatusCode = 0
		return responseData, err
	}
	//resultMap := map[string]interface{}{"statusCode": 1, "data": data, "urlPathParam": tableName}
	responseData.StatusCode = 1
	responseData.UrlPathParam = urlPathParam
	return responseData, nil
}

func funcJsonMarshal(obj interface{}) (string, error) {
	jsonByte, err := json.Marshal(obj)
	return string(jsonByte), err
}
