package main

import (
	"context"
	"html/template"

	"github.com/blevesearch/bleve/v2"
)

var funcMap = template.FuncMap{

	"basePath":   funcBasePath,
	"T":          funcT,
	"safeHTML":   funcSafeHTML,
	"relURL":     funcRelURL,
	"indexFiled": funcIndexFiled,
	"site":       funcSite,
	"navMenu":    funcNavMenu,
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

// funcIndexFiled 根据indexName查找字段
func funcIndexFiled(indexName string) ([]IndexFieldStruct, error) {
	ctx := context.Background()
	indexField, err := findIndexFieldStruct(ctx, indexName)
	return indexField, err
}

// 站点信息
func funcSite() (map[string]interface{}, error) {
	idQuery := bleveNewTermQuery("gpress")
	// 指定查询的字段
	idQuery.SetField("id")
	searchRequest := bleve.NewSearchRequest(idQuery)
	// 指定返回的字段
	searchRequest.Fields = []string{"*"}
	searchIndex, _, _ := openBleveIndex(indexSiteName)
	searchResult, err := searchIndex.Search(searchRequest)
	if err != nil {
		return make(map[string]interface{}, 0), err
	}
	data, err := result2Map(indexSiteName, searchResult)
	return data, err
}

// 菜单信息
func funcNavMenu() ([]map[string]interface{}, error) {
	searchIndex, _, _ := openBleveIndex(indexNavMenuName)
	query := bleve.NewQueryStringQuery("*")
	searchRequest := bleve.NewSearchRequestOptions(query, 100, 0, false)
	// 指定返回的字段
	searchRequest.Fields = []string{"*"}
	searchResult, err := searchIndex.Search(searchRequest)
	if err != nil {
		return make([]map[string]interface{}, 0), err
	}
	data, err := result2SliceMap(indexNavMenuName, searchResult)
	return data, err
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
