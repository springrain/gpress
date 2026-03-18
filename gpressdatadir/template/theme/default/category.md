<!-- 查询 content 列表 -->
{{ $contentSQL :="* FROM content WHERE status in (1,2) and category_id=? order by  status desc,sortno desc" }}
{{ $categorySQL := "* FROM category WHERE id=? and status<3 order by status desc, sortno desc"}}


{{ if eq .userType 1}}
  {{ $contentSQL ="* FROM content WHERE category_id=? order by sortno desc" }}
  {{ $categorySQL = "* FROM category WHERE id=? order by sortno desc"}}
{{end}}

{{ $selectList := selectList "content" .q .pageNo 20 $contentSQL .UrlPathParam }}

{{ $nav := selectOne "category" $categorySQL .UrlPathParam }}

title: {{ $nav.Name }}  
description: {{ $nav.Description }}  
keyword: {{ $nav.Keyword }}  
---

{{ range $k,$v := $selectList.Data }}
# [{{ .Title }}]({{basePath}}{{ trimSlash $v.Id }}) 
- UpdateTime: {{ slice .UpdateTime 0 10 }}  
- Summary: {{ safeHTML .Summary }}   
{{ end }}