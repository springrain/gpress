hugo源码有个depts依赖包,几乎依赖了hugo的所有包,应该是为了处理go不允许互相依赖的中间包,这样就造成源码修改hugo极其复杂,放弃!!!!

https://gohugo.io/functions/

1.函数引入模板修改为原生引用模板的方式,主要原因是需要重写template,支持上下文,比较复杂

{{ partial "head.html" . }}     修改为 {{ template "theme/partial/head.html" }}
{{ partial "header.html" . }}   修改为 {{ template "theme/partial/header.html" }}
{{ partial "slideout.html" . }} 修改为 {{ template "theme/partial/slideout.html" }}
{{ partial "comments.html" . }} 修改为 {{ template "theme/partial/comments.html" }}
{{ partial "footer.html" . }}   修改为 {{ template "theme/partial/footer.html" }}
{{ partial "scripts.html" . }}  修改为 {{ template "theme/partial/scripts.html" }}

可以批量替换

2.取消sass/scss的编译功能,hugo要么使用go-libsass的c++库或者操作系统安装dart-sass-embedded扩展,都比较复杂
使用dart-sass不打入二进制包,调用命令进行编译
  下载 https://github.com/sass/dart-sass    
  文档 https://sass-lang.com/documentation/cli/dart-sass  

  把hugo的
  ```go
<!-- Styles -->
{{ $style := resources.Get "sass/main.scss" | toCSS | minify | fingerprint }}
<link href="{{ $style.RelPermalink }}" rel="stylesheet">
  ```
  修改为

```go
<link href="{{sass "sass/main.scss" }}" rel="stylesheet">
```


  ```bat
   ### 编译 assets\sass 下所有的 sass/scss文件 到 resources\_gen\assets\scss\sass 目录下
   dart-sass\windows\sass.bat --style=compressed --charset --no-source-map assets\sass:resources\_gen\assets\scss\sass

   ### --style=compressed 压缩
   ### --charset 使用编码
   ### --no-source-map 不生成源码map文件
   ### : 路径分隔符
  ```

