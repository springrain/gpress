<!-- 查询 content -->
{{ $contentSQL := "* FROM content WHERE id=? order by sortno desc" }}

{{$contentID := printf "/single/%s" (lastURI .UrlPathParam) }}

{{ $content := selectOne "content" $contentSQL $contentID }}

title: {{ $content.Title }}  
CreateTime: {{ $content.CreateTime }}  
UpdateTime: {{ $content.UpdateTime }}  
CategoryName: {{ $content.CategoryName }}  
--- 

{{ $content.Markdown }}