<!-- 查询 content -->
{{ $contentSQL := "* FROM content WHERE id=? and status<3 order by status desc, sortno desc" }}

{{ if eq .userType 1}}
  {{ $contentSQL ="* FROM content WHERE id=? order by sortno desc" }}
{{end}}

{{ $content := selectOne "content" $contentSQL .UrlPathParam }}

title: {{ $content.Title }}  
CreateTime: {{ $content.CreateTime }}  
UpdateTime: {{ slice $content.UpdateTime 0 10 }}  
CategoryName: {{ $content.CategoryName }}  
TOC: {{ $content.Toc }}  
---

{{ $content.Markdown }}
