package hugo

import (
	"text/template"
)

// FuncT 多语言i18n适配,例如 {{ T "nextPage" }}
func FuncT(key string) (string, error) {
	return key, nil
}

// FuncSafeHTML 转义html字符串
func FuncSafeHTML(html string) (string, error) {
	ss := template.HTMLEscapeString(html)
	return ss, nil
}

// FuncRelURL 真实的url
func FuncRelURL(url string) (string, error) {
	return FuncSafeHTML(url)
}
