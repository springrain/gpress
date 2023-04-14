package main

import (
	"context"
	"errors"
	"html/template"
	"strings"

	"gitee.com/chunanyong/zorm"
)

var funcMap = template.FuncMap{

	"basePath":   funcBasePath,
	"addInt":     funcAddInt,
	"addFloat":   funcAddFloat,
	"T":          funcT,
	"safeHTML":   funcSafeHTML,
	"relURL":     funcRelURL,
	"tableFiled": funcTableFiled,
	"site":       funcSite,
	"navMenu":    funcNavMenu,
	"selectList": funcSelectList,
	"selectOne":  funcSelectOne,
	"analyzer":   funcAnalyzer,
	"fieldType":  funcFieldType,
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

// funcRelURL 拼接url路径的
func funcRelURL(url string) (template.HTML, error) {
	return funcSafeHTML(themePath + url)
}

// funcTableFiled 根据indexName查找字段
func funcTableFiled(tableName string) ([]TableFieldStruct, error) {
	ctx := context.Background()
	tableField, err := findTableFieldStruct(ctx, tableName, 0)
	return tableField, err
}

// 站点信息
func funcSite() (map[string]interface{}, error) {
	finder := zorm.NewSelectFinder(tableSiteName, "*").Append(" WHERE id=?", "gpress")
	rowMap, err := zorm.QueryRowMap(context.Background(), finder)
	if err != nil {
		return make(map[string]interface{}, 0), err
	}
	return rowMap, err
}

// 菜单信息
func funcNavMenu() ([]map[string]interface{}, error) {
	finder := zorm.NewSelectFinder(tableNavMenuName)
	finder.Append(" order by sortNo desc")
	page := zorm.NewPage()
	page.PageSize = 200
	return zorm.QueryMap(context.Background(), finder, page)
}

var analyzerMap = map[string]string{commaAnalyzerName: "逗号分词器", gseAnalyzerName: "默认分词器", keywordAnalyzerName: "不分词", numericAnalyzerName: "数字分词器", datetimeAnalyzerName: "日期分词器"}

func funcAnalyzer() map[string]string {
	return analyzerMap
}

var fieldTypeMap = map[int]string{1: "数字", 2: "日期", 3: "文本框", 4: "文本域", 5: "富文本", 6: "下拉框", 7: "单选", 8: "多选", 9: "上传图片", 10: "上传附件", 11: "轮播图", 12: "音频", 13: "视频"}

func funcFieldType() map[int]string {
	return fieldTypeMap
}

func funcAddInt(x, y int) int {
	return x + y
}
func funcAddFloat(x, y float64) float64 {
	return x + y
}

func funcSelectList(q string, pageNo int, sql string, values ...interface{}) (map[string]interface{}, error) {
	errMap := map[string]interface{}{"statusCode": 0}
	sql = strings.TrimSpace(sql)
	if sql == "" || strings.Contains(sql, ";") {
		err := errors.New("sql语句错误")
		errMap["err"] = err
		return errMap, err
	}

	finder := zorm.NewFinder().Append("SELECT")
	if q != "" { // 如果有搜索关键字
		whereSQL := strings.ToLower(sql)
		i := strings.Index(whereSQL, " where ")
		if i < 0 { // 没有where
			finder.Append(sql, values...)
			finder.Append(" where id in (select id from fts_content where fts_content match jieba_query(?) ) ", q)
		} else {
			finder.Append(sql[:i+7]+" id in (select id from fts_content where fts_content match jieba_query(?) ) and ", q)
			finder.Append(sql[i+7:], values...)
		}

	} else {
		finder.Append(sql, values...)
	}

	//finder.Append("order by sortNo desc")
	page := zorm.NewPage()
	page.PageNo = pageNo
	data, err := zorm.QueryMap(context.Background(), finder, page)
	if err != nil {
		errMap["err"] = err
		return errMap, err
	}

	resultMap := map[string]interface{}{"statusCode": 1, "data": data, "page": page}
	return resultMap, err
}

func funcSelectOne(sql string, values ...interface{}) (map[string]interface{}, error) {
	errMap := map[string]interface{}{"statusCode": 0}
	sql = strings.TrimSpace(sql)
	if sql == "" || strings.Contains(sql, ";") {
		err := errors.New("sql语句错误")
		errMap["err"] = err
		return errMap, err
	}
	finder := zorm.NewFinder().Append("SELECT")
	finder.Append(sql, values...)

	page := zorm.NewPage()
	page.PageSize = 1
	page.PageNo = 1
	resultMaps, err := zorm.QueryMap(context.Background(), finder, page)

	if err != nil {
		errMap["err"] = err
		return errMap, err
	}
	if len(resultMaps) < 1 {
		return errMap, err
	}
	resultMap := resultMaps[0]
	//resultMap := map[string]interface{}{"statusCode": 1, "data": data, "urlPathParam": tableName}
	resultMap["statusCode"] = 1
	return resultMap, err
}

/*
func funcThemePath() string {
	return themePath
}

func funcThemeFile(file string) string {
	return "/theme/" + config.Theme + "/" + file
}

// 测试自定义函数
func funcMD5(in string) ([]string, error) {
	list := make([]string, 2)

	hash := md5.Sum([]byte(in))
	list[0] = in
	list[1] = hex.EncodeToString(hash[:])
	return list, nil
}

// funcSass 编译sass,生成css
func funcSass(sassFile string) (string, error) {
	sassPath := templateDir + "theme/" + config.Theme + "/assets/" + sassFile
	pathHash := hashSha256(sassPath)

	//生成的css路径
	filePath := templateDir + "theme/" + config.Theme + "/css/" + pathHash + ".css"
	//filePath = filepath.FromSlash(filePath)

	//url 访问路径
	fileUrl := themePath + "css/" + pathHash + ".css"

	if pathExists(filePath) { //如果文件已经存在了,直接返回
		return funcSafeHTML(fileUrl)
	}
	var cmd *exec.Cmd
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	if goos == "windows" {
		cmdStr := datadir + "dart-sass/" + goos + "-" + goarch + "/sass.bat --style=compressed --charset --no-source-map " + sassPath + ":" + filePath
		// 分隔符统一系统符号
		cmdStr = filepath.FromSlash(cmdStr)
		cmd = exec.Command("cmd", "/C", cmdStr) // windows
	} else if goos == "linux" || goos == "darwin" {
		cmdStr := datadir + "dart-sass/" + goos + "-" + goarch + "/sass --style=compressed --charset --no-source-map " + sassPath + ":" + filePath
		cmd = exec.Command("bash", "-c", cmdStr) // mac or linux
	}
	err := cmd.Run()
	if err != nil {
		FuncLogError(err)
		return fileUrl, err
	}
	fileUrl, err = funcSafeHTML(fileUrl)
	//增加静态资源映射,使用了static文件夹,不需要再映射了
	//h.StaticFile(fileUrl, filePath)
	return fileUrl, err
}
*/
