package main

import "html/template"

// funcT 多语言i18n适配,例如 {{ T "nextPage" }}
func funcT(key string) (string, error) {
	return key, nil
}

// funcSafeHTML 转义html字符串
func funcSafeHTML(html string) (string, error) {
	ss := template.HTMLEscapeString(html)
	return ss, nil
}

// funcRelURL 真实的url
func funcRelURL(url string) (string, error) {
	return funcSafeHTML(url)
}
