<!-- 查询 content -->
{{ $contentSQL := "* FROM content WHERE id=? order by sortno desc" }}

{{$contentID := printf "/single/%s" (lastURI .UrlPathParam) }}

{{ $content := selectOne "content" $contentSQL $contentID }}

title: {{ $content.Title }}  
CreateTime: {{ $content.CreateTime }}  
UpdateTime: {{ slice $content.UpdateTime 0 10 }}  
CategoryName: {{ $content.CategoryName }}  
TOC: {{ $content.Toc }}  
--- 

{{ $content.Markdown }}