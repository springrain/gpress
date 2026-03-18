<!-- 查询 content 列表 -->
{{ $selectList := selectList "content" .q .pageNo 20  "* FROM content WHERE status in (1,2) and tag=? order by status desc,sortno desc" .UrlPathParam }}

Title: {{ .UrlPathParam }}
---

{{ range $k,$v := $selectList.Data }}
# [{{ .Title }}]({{basePath}}{{ trimSlash $v.Id }}) 
- UpdateTime: {{ .UpdateTime }}  
- Summary: {{ safeHTML .Summary }}   
{{ end }}