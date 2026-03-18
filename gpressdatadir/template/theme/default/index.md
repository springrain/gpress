{{ $site:=site }}
<!-- 查询 content 列表 -->
{{ $selectList := selectList "content" .q .pageNo 20 "* FROM content WHERE status in (1,2) order by  status desc,sortno desc"  }}

title: {{ $site.Title }}  
description: {{ $site.Description }}  
keyword: {{ $site.Keyword }}  
---

{{ range $k,$v := $selectList.Data }}
# [{{ .Title }}]({{basePath}}{{ trimSlash $v.Id }}) 
- UpdateTime: {{ slice .UpdateTime 0 10 }}  
- Summary: {{ safeHTML .Summary }}   
{{ end }}