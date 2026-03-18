<!-- 查询 content -->
{{ $contentSQL := "* FROM content WHERE id=? and status<3 order by status desc, sortno desc" }}

{{ if eq .userType 1}}
  {{ $contentSQL ="* FROM content WHERE id=? order by sortno desc" }}
{{end}}

{{ $content := selectOne "content" $contentSQL .UrlPathParam }}

Title: {{ $content.Title }}  
CreateTime: {{ $content.CreateTime }}  
UpdateTime: {{ $content.UpdateTime }}  
CategoryName: {{ $content.CategoryName }}  
---

{{ safeHTML $content.Markdown }}
