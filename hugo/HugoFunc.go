package hugo

import "html/template"

// T 多语言i18n适配,例如 {{ T "nextPage" }}
func T(key string) (string, error) {

	return key, nil
}

//safeHTML 转义html字符串
func safeHTML(html string) (string, error) {
	ss := template.HTMLEscapeString(html)
	return ss, nil
}

//relURL 真实的url
func relURL(url string) (string, error) {
	return safeHTML(url)
}
