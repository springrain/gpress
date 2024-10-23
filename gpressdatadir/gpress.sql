CREATE TABLE IF NOT EXISTS config (
		id TEXT PRIMARY KEY   NOT NULL,
		basePath         TEXT NOT NULL,
		jwtSecret        TEXT NOT NULL,
		jwttokenKey      TEXT NOT NULL,
		serverPort       TEXT NOT NULL,
		timeout          INT  NOT NULL,
		maxRequestBodySize INT,
		proxy            TEXT NULL,
		createTime       TEXT,
		updateTime       TEXT,
		createUser       TEXT,
		sortNo           int,
		status           int  
	 ) strict ;

CREATE TABLE IF NOT EXISTS user (
		id TEXT PRIMARY KEY     NOT NULL,
		account         TEXT  NOT NULL,
		password         TEXT   NOT NULL,
		userName         TEXT NOT NULL,
		chainType        TEXT,
		chainAddress     TEXT,
		createTime       TEXT,
		updateTime       TEXT,
		createUser       TEXT,
		sortNo           int,
		status           int  
	 ) strict ;

CREATE TABLE IF NOT EXISTS category (
		id TEXT PRIMARY KEY     NOT NULL,
		name          TEXT  NOT NULL,
		hrefURL           TEXT,
		hrefTarget        TEXT,
		pid        TEXT,
		themePC        TEXT,
		moduleID        TEXT,
		comCode        TEXT,
		templateFile        TEXT,
		childTemplateFile        TEXT,
		keyword           TEXT,
		description           TEXT,
		createTime        TEXT,
		updateTime        TEXT,
		createUser        TEXT,
		sortNo            int NOT NULL,
		status            int NOT NULL
	 ) strict ;
INSERT INTO category (status,sortNo,createUser,updateTime,createTime,childTemplateFile,templateFile,comCode,moduleID,themePC,pid,hrefTarget,hrefURL,name,id) VALUES (1,2,NULL,'2023-06-27 22:41:20','2023-06-27 22:41:20',NULL,NULL,',web,',NULL,NULL,NULL,'',NULL,'Web','web');
INSERT INTO category (status,sortNo,createUser,updateTime,createTime,childTemplateFile,templateFile,comCode,moduleID,themePC,pid,hrefTarget,hrefURL,name,id) VALUES (1,1,NULL,'2023-06-27 22:41:20','2023-06-27 22:41:20',NULL,NULL,',about,',NULL,NULL,NULL,'','/post/about','About','about');

CREATE TABLE IF NOT EXISTS content (
		id TEXT PRIMARY KEY     NOT NULL,
		moduleID         TEXT  ,
		title         TEXT   NOT NULL,
		keyword           TEXT,
		description           TEXT,
		hrefURL           TEXT,
		subtitle           TEXT,
		categoryID           TEXT,
		categoryName           TEXT,
		comCode        TEXT,
		templateFile           TEXT,
		author           TEXT,
		tag           TEXT,
		toc           TEXT,
		summary           TEXT,
		content           TEXT,
		markdown          TEXT,
		thumbnail         TEXT,
		createTime        TEXT,
		updateTime        TEXT,
		createUser        TEXT,
		sortNo            int NOT NULL,
		status            int NOT NULL
	 ) strict ;

CREATE TABLE IF NOT EXISTS site (
		id TEXT PRIMARY KEY     NOT NULL,
		title         TEXT  NOT NULL,
		name         TEXT   NOT NULL,
		domain         TEXT,
		keyword         TEXT,
		description         TEXT,
		theme         TEXT NOT NULL,
		themePC         TEXT,
		themeWAP         TEXT,
		siteThemeWEIXIN         TEXT,
		logo         TEXT,
		favicon         TEXT,
		footer         TEXT,
		createTime        TEXT,
		updateTime        TEXT,
		createUser        TEXT,
		sortNo            int NOT NULL,
		status            int NOT NULL
	 ) strict ;
INSERT INTO site (status,sortNo,createUser,updateTime,createTime,footer,favicon,logo,siteThemeWEIXIN,themeWAP,themePC,theme,description,keyword,domain,name,title,id)VALUES (1,1,NULL,NULL,NULL,'<div class="copyright"><span class="copyright-year">&copy; 2008 - 2024<span class="author">jiagou.com 版权所有 <a href=''https://beian.miit.gov.cn'' target=''_blank''>豫ICP备xxxxx号</a>   <a href=''http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=xxxx''  target=''_blank''><img src=''/public/gongan.png''>豫公网安备xxxxx号</a></span></span></div>','public/favicon.png','public/logo.png','default','default','default','default','Web3内容平台,Hertz + Go template + FTS5全文检索,支持以太坊和百度超级链,兼容Hugo、WordPress生态,使用Wasm扩展插件,只需200M内存','gpress,web3,Hugo,WordPress,以太坊,百度超级链','jiagou.com','架构','jiagou','gpress_site');

CREATE VIRTUAL TABLE IF NOT EXISTS fts_content USING fts5(
		title, 
		keyword, 
		description,
		subtitle,
		categoryName,
		summary,
		toc,
		tag,
		author, 

	    tokenize = 'simple 0',
		content='content', 
		content_rowid='rowid'
	);
CREATE TRIGGER IF NOT EXISTS trigger_content_insert AFTER INSERT ON content
		BEGIN
			INSERT INTO fts_content (rowid, title, keyword, description,subtitle,categoryName,summary,toc,tag,author)
			VALUES (new.rowid,  new.title, new.keyword, new.description,new.subtitle,new.categoryName,new.summary,new.toc,new.tag,new.author);
		END;
	
	CREATE TRIGGER IF NOT EXISTS trigger_content_delete AFTER DELETE ON content
		BEGIN
			INSERT INTO fts_content (fts_content,  title, keyword, description,subtitle,categoryName,summary,toc,tag,author)
			VALUES ('delete',  old.title, old.keyword, old.description,old.subtitle,old.categoryName,old.summary,old.toc,old.tag,old.author);
		END;
	
	CREATE TRIGGER IF NOT EXISTS trigger_content_update AFTER UPDATE ON content
		BEGIN
			INSERT INTO fts_content (fts_content, rowid, title, keyword, description,subtitle,categoryName,summary,toc,tag,author)
			VALUES ('delete', old.rowid,  old.title, old.keyword, old.description,old.subtitle,old.categoryName,old.summary,old.toc,old.tag,old.author);
			INSERT INTO fts_content (rowid, title, keyword, description,subtitle,categoryName,summary,toc,tag,author)
			VALUES (new.rowid, new.title, new.keyword, new.description,new.subtitle,new.categoryName,new.summary,new.toc,new.tag,new.author);
		END;

INSERT INTO content (
                        status,
                        sortNo,
                        createUser,
                        updateTime,
                        createTime,
                        thumbnail,
                        markdown,
                        content,
                        summary,
                        toc,
                        tag,
                        author,
                        templateFile,
                        comCode,
                        categoryName,
                        categoryID,
                        subtitle,
                        hrefURL,
                        description,
                        keyword,
                        title,
                        moduleID,
                        id
                    )
                    VALUES (
                        0,
                        100,
                        NULL,
                        '2024-04-11 11:02:40',
                        '2023-06-27 22:43:53',
                        NULL,
                        '本站服务器配置:阿里云张家口机房,ecs.t6-c4m1.large,2核CPU,512M内存,20G高效云盘,RockyLinux 9 .  
使用[gpress](https://gitee.com/gpress/gpress)迁移了[Hugo](https://github.com/gohugoio/hugo)的[even](https://github.com/olOwOlo/hugo-theme-even)主题和markdown文件.  

我所见识过的一切都将消失一空,就如眼泪消逝在雨中......  
不妨大胆一些,大胆一些......  

小项目:  
* [springrain](https://gitee.com/chunanyong/springrain)
* [zorm](https://gitee.com/chunanyong/zorm)
* [gpress](https://gitee.com/gpress/gpress)
* [gowe](https://gitee.com/chunanyong/gowe)',
                        '<p>本站服务器配置:阿里云张家口机房,ecs.t6-c4m1.large,2核CPU,512M内存,20G高效云盘,RockyLinux 9 .<br>
使用<a href="https://gitee.com/gpress/gpress">gpress</a>迁移了<a href="https://github.com/gohugoio/hugo">hugo</a>的<a href="https://github.com/olOwOlo/hugo-theme-even">even</a>主题和markdown文件.</p>
<p>我所见识过的一切都将消失一空,就如眼泪消逝在雨中......<br>
不妨大胆一些,大胆一些......</p>
<p>小项目:</p>
<ul>
<li><a href="https://gitee.com/chunanyong/springrain">springrain</a></li>
<li><a href="https://gitee.com/chunanyong/zorm">zorm</a></li>
<li><a href="https://gitee.com/gpress/gpress">gpress</a></li>
<li><a href="https://gitee.com/chunanyong/gowe">gowe</a></li>
</ul>
',
                        '本站服务器配置:阿里云张家口机房,ecs.t6-c4m1.large,2核CPU,512M内存,20G高效云盘,RockyLinux 9.使用Hugo和even模板,编译成静态文件,Nginx作为WEB服务器.我所见识过的一切都将消失一空,就如眼泪消逝在雨中......	不妨大胆一些,大胆一些......',
                        '',
                        NULL,
                        'springrain',
                        NULL,
                        ',about,',
                        'About',
                        'about',
                        NULL,
                        NULL,
                        NULL,
                        NULL,
                        'about',
                        NULL,
                        'about'
                    ),
                    (
                        1,
                        101,
                        '',
                        '2024-04-11 11:42:36',
                        '2024-04-06 20:31:32',
                        '',
                        '# 介绍  
Web3内容平台,Hertz + Go template + FTS5全文检索,支持以太坊和百度超级链,兼容Hugo、WordPress生态,使用Wasm扩展插件,只需200M内存  
    
**作为静态站点：** gpress生成的静态文件和Hugo一致,也可以简单认为gpress是Hugo的后台管理,兼容Hugo主题生态,已迁移多款Hugo主题:[even](gitee.com/gpress/gpress/tree/master/gpressdatadir/template/theme/default)、[doks](gitee.com/gpress/gpress-doks)、[book](gitee.com/gpress/gpress-book)、[geekdoc](gitee.com/gpress/gpress-geekdoc)......   
**作为动态站点：** gpress功能简单,只有7个菜单,5张表,5000行代码,使用SQLite,一键启动,只需200M内存,支持全文检索.兼容WordPress主题生态,已迁移多款WordPress主题:[generatepress](https://gitee.com/gpress/wp-generatepress)、[astra](https://gitee.com/gpress/wp-astra)......   
**作为Web3：** gpress已支持以太坊和百度超级链账户体系,会基于Wasm持续迭代去中心功能,让数据自由一点点......  
**作为后浪：** 相对于Hugo、WordPress等优秀的内容平台,gpress还有很多不足,功能简单而又稚嫩......  
**帮助文档：** [点击查看帮助文档](https://gitee.com/gpress/gpress/blob/master/gpressdatadir/public/doc/index.md)    
 
个人博客 [jiagou.com](jiagou.com/post/about) 使用gpress搭建,搜索和后台管理是动态,其他是静态页面.  
![](/public/index.png "")
# 开发环境  
gprss使用了 ```https://github.com/wangfenjin/simple``` 作为FTS5的全文检索扩展,编译好的libsimple文件放到 ```gpressdatadir/fts5``` 目录下,如果gpress启动报错连不上数据库,请检查libsimple文件是否正确,如果需要重新编译libsimple,请参考 https://github.com/wangfenjin/simple.  

默认端口660,后台管理地址 http://127.0.0.1:660/admin/login    
需要先解压```gpressdatadir/dict.zip```      
运行 ```go run --tags "fts5" .```     
打包: ```go build --tags "fts5" -ldflags "-w -s"```  

开发环境需要配置CGO编译,设置```set CGO_ENABLED=1```,下载[mingw64](https://github.com/niXman/mingw-builds-binaries/releases)和[cmake](https://cmake.org/download/),并把bin配置到环境变量,注意把```mingw64/bin/mingw32-make.exe``` 改名为 ```make.exe```  
注意修改vscode的launch.json,增加 ``` ,"buildFlags": "--tags=fts5" ``` 用于调试fts5    
test需要手动测试:``` go test -timeout 30s --tags "fts5"  -run ^TestReadmks$ gitee.com/gpress/gpress ```  
打包: ``` go build --tags "fts5" -ldflags "-w -s" ```   
重新编译simple时,建议使用```https://github.com/wangfenjin/simple```编译好的.  
注意修改widnows编译脚本,去掉 mingw64 编译依赖的```libgcc_s_seh-1.dll```和```libstdc++-6.dll```,同时关闭```BUILD_TEST_EXAMPLE```,有冲突
```bat
rmdir /q /s build
mkdir build && cd build
cmake .. -G "Unix Makefiles" -DBUILD_TEST_EXAMPLE=OFF -DCMAKE_INSTALL_PREFIX=release -DCMAKE_CXX_FLAGS="-static-libgcc -static-libstdc++" -DCMAKE_EXE_LINKER_FLAGS="-Wl,-Bstatic -lstdc++ -lpthread -Wl,-Bdynamic"
make && make install
```

# 静态化
后台 ```刷新站点``` 功能会生成静态html文件到 ```statichtml``` 目录,同时生成```gzip_static```文件,需要把正在使用的主题的 ```css,js,image```和```gpressdatadir/public```目录复制到 ```statichtml```目录下,也可以用Nginx反向代理指定目录.    
nginx 配置示例如下:
```conf
### 当前在用主题(default)的css文件
location ~ ^/css/ {
    #gzip_static on;
    root /data/gpress/gpressdatadir/template/theme/default;  
}
### 当前在用主题(default)的js文件
location ~ ^/js/ {
    #gzip_static on;
    root /data/gpress/gpressdatadir/template/theme/default;  
}
### 当前在用主题(default)的image文件
location ~ ^/image/ {
    root /data/gpress/gpressdatadir/template/theme/default;  
}
### search-data.json FlexSearch搜索的JSON数据
location ~ ^/public/search-data.json {
    #gzip_static on;
    root /data/gpress/gpressdatadir;  
}
### public 公共文件
location ~ ^/public/ {
    root /data/gpress/gpressdatadir;  
}
    
### admin 后台管理,请求动态服务
location ~ ^/admin/ {
    proxy_redirect     off;
    proxy_set_header   Host      $host;
    proxy_set_header   X-Real-IP $remote_addr;
    proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
    proxy_set_header   X-Forwarded-Proto $scheme;
    proxy_pass  http://127.0.0.1:660;  
}
###  静态html目录
location / {
    proxy_redirect     off;
    proxy_set_header   Host      $host;
    proxy_set_header   X-Real-IP $remote_addr;
    proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
    proxy_set_header   X-Forwarded-Proto $scheme; 
    ## 存在q查询参数,使用动态服务.也支持FlexSearch解析public/search-data.json
    if ($arg_q) { 
       proxy_pass  http://127.0.0.1:660;  
       break;
    }

    ### 开启gzip静态压缩
    #gzip_static on;
    ### Nginx 1.26+ 不需要再进行302重定向到目录下的index.html,gzip_static也会生效.这段配置留作记录.
    ##if ( -d $request_filename ) {
        ## 不是 / 结尾
    ##    rewrite [^\/]$ $uri/index.html redirect;
        ##以 / 结尾的
    ##    rewrite ^(.*) ${uri}index.html redirect;      
    ##}


    
    root   /data/gpress/gpressdatadir/statichtml;
    index  index.html index.htm;
}

```  
# 阿里云计算巢
[点击部署gpress到阿里云计算巢](https://computenest.console.aliyun.com/service/instance/create/cn-hangzhou?type=user&ServiceId=service-d4000c9b22c54e5cbffe),也可以单独购买阿里云最低配服务器,进行部署.选择```张家口机房```,规格```ecs.t6-c4m1.large```,配置```2核CPU 0.5G内存 20G高效云盘 RockyLinux9 按使用流量-带宽峰值80M```,一年100元,五年200元左右.  ',
                        '<h1 id="介绍">介绍</h1>
<p>Web3内容平台,Hertz + Go template + FTS5全文检索,支持以太坊和百度超级链,兼容Hugo、WordPress生态,使用Wasm扩展插件,只需200M内存</p>
<p><strong>作为静态站点：</strong> gpress生成的静态文件和Hugo一致,也可以简单认为gpress是Hugo的后台管理,兼容Hugo主题生态,已迁移多款Hugo主题:<a href="gitee.com/gpress/gpress/tree/master/gpressdatadir/template/theme/default">even</a>、<a href="gitee.com/gpress/gpress-doks">doks</a>、<a href="gitee.com/gpress/gpress-book">book</a>、<a href="gitee.com/gpress/gpress-geekdoc">geekdoc</a>......<br>
<strong>作为动态站点：</strong> gpress功能简单,只有7个菜单,5张表,5000行代码,使用SQLite,一键启动,只需200M内存,支持全文检索.兼容WordPress主题生态,已迁移多款WordPress主题:<a href="https://gitee.com/gpress/wp-generatepress">generatepress</a>、<a href="https://gitee.com/gpress/wp-astra">astra</a>......<br>
<strong>作为Web3：</strong> gpress已支持以太坊和百度超级链账户体系,会基于Wasm持续迭代去中心功能,让数据自由一点点......<br>
<strong>作为后浪：</strong> 相对于Hugo、WordPress等优秀的内容平台,gpress还有很多不足,功能简单而又稚嫩......<br>
<strong>帮助文档：</strong> <a href="https://gitee.com/gpress/gpress/blob/master/gpressdatadir/public/doc/index.md">点击查看帮助文档</a></p>
<p>个人博客 <a href="jiagou.com/post/about">jiagou.com</a> 使用gpress搭建,搜索和后台管理是动态,其他是静态页面.<br>
<img src="/public/index.png" alt="" title=""></p>
<h1 id="开发环境">开发环境</h1>
<p>gprss使用了 <code>https://github.com/wangfenjin/simple</code> 作为FTS5的全文检索扩展,编译好的libsimple文件放到 <code>gpressdatadir/fts5</code> 目录下,如果gpress启动报错连不上数据库,请检查libsimple文件是否正确,如果需要重新编译libsimple,请参考 <a href="https://github.com/wangfenjin/simple">https://github.com/wangfenjin/simple</a>.</p>
<p>默认端口660,后台管理地址 http://127.0.0.1:660/admin/login <br>
需要先解压<code>gpressdatadir/dict.zip</code>   <br>
运行 <code>go run --tags &quot;fts5&quot; .</code>  <br>
打包: <code>go build --tags &quot;fts5&quot; -ldflags &quot;-w -s&quot;</code></p>
<p>开发环境需要配置CGO编译,设置<code>set CGO_ENABLED=1</code>,下载<a href="https://github.com/niXman/mingw-builds-binaries/releases">mingw64</a>和<a href="https://cmake.org/download/">cmake</a>,并把bin配置到环境变量,注意把<code>mingw64/bin/mingw32-make.exe</code> 改名为 <code>make.exe</code><br>
注意修改vscode的launch.json,增加 <code>,&quot;buildFlags&quot;: &quot;--tags=fts5&quot;</code> 用于调试fts5 <br>
test需要手动测试:<code>go test -timeout 30s --tags &quot;fts5&quot;  -run ^TestReadmks$ gitee.com/gpress/gpress</code><br>
打包: <code>go build --tags &quot;fts5&quot; -ldflags &quot;-w -s&quot;</code><br>
重新编译simple时,建议使用<code>https://github.com/wangfenjin/simple</code>编译好的.<br>
注意修改widnows编译脚本,去掉 mingw64 编译依赖的<code>libgcc_s_seh-1.dll</code>和<code>libstdc++-6.dll</code>,同时关闭<code>BUILD_TEST_EXAMPLE</code>,有冲突</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bat" data-lang="bat"><span class="line"><span class="cl"><span class="k">rmdir</span> /q /s build
</span></span><span class="line"><span class="cl"><span class="k">mkdir</span> build <span class="p">&amp;&amp;</span> <span class="k">cd</span> build
</span></span><span class="line"><span class="cl">cmake .. -G <span class="s2">&#34;Unix Makefiles&#34;</span> -DBUILD_TEST_EXAMPLE=OFF -DCMAKE_INSTALL_PREFIX=release -DCMAKE_CXX_FLAGS=<span class="s2">&#34;-static-libgcc -static-libstdc++&#34;</span> -DCMAKE_EXE_LINKER_FLAGS=<span class="s2">&#34;-Wl,-Bstatic -lstdc++ -lpthread -Wl,-Bdynamic&#34;</span>
</span></span><span class="line"><span class="cl">make <span class="p">&amp;&amp;</span> make install
</span></span></code></pre></div><h1 id="静态化">静态化</h1>
<p>后台 <code>刷新站点</code> 功能会生成静态html文件到 <code>statichtml</code> 目录,同时生成<code>gzip_static</code>文件,需要把正在使用的主题的 <code>css,js,image</code>和<code>gpressdatadir/public</code>目录复制到 <code>statichtml</code>目录下,也可以用Nginx反向代理指定目录. <br>
nginx 配置示例如下:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-fallback" data-lang="fallback"><span class="line"><span class="cl">### 当前在用主题(default)的css文件
</span></span><span class="line"><span class="cl">location ~ ^/css/ {
</span></span><span class="line"><span class="cl">    #gzip_static on;
</span></span><span class="line"><span class="cl">    root /data/gpress/gpressdatadir/template/theme/default;  
</span></span><span class="line"><span class="cl">}
</span></span><span class="line"><span class="cl">### 当前在用主题(default)的js文件
</span></span><span class="line"><span class="cl">location ~ ^/js/ {
</span></span><span class="line"><span class="cl">    #gzip_static on;
</span></span><span class="line"><span class="cl">    root /data/gpress/gpressdatadir/template/theme/default;  
</span></span><span class="line"><span class="cl">}
</span></span><span class="line"><span class="cl">### 当前在用主题(default)的image文件
</span></span><span class="line"><span class="cl">location ~ ^/image/ {
</span></span><span class="line"><span class="cl">    root /data/gpress/gpressdatadir/template/theme/default;  
</span></span><span class="line"><span class="cl">}
</span></span><span class="line"><span class="cl">### search-data.json FlexSearch搜索的JSON数据
</span></span><span class="line"><span class="cl">location ~ ^/public/search-data.json {
</span></span><span class="line"><span class="cl">    #gzip_static on;
</span></span><span class="line"><span class="cl">    root /data/gpress/gpressdatadir;  
</span></span><span class="line"><span class="cl">}
</span></span><span class="line"><span class="cl">### public 公共文件
</span></span><span class="line"><span class="cl">location ~ ^/public/ {
</span></span><span class="line"><span class="cl">    root /data/gpress/gpressdatadir;  
</span></span><span class="line"><span class="cl">}
</span></span><span class="line"><span class="cl">    
</span></span><span class="line"><span class="cl">### admin 后台管理,请求动态服务
</span></span><span class="line"><span class="cl">location ~ ^/admin/ {
</span></span><span class="line"><span class="cl">    proxy_redirect     off;
</span></span><span class="line"><span class="cl">    proxy_set_header   Host      $host;
</span></span><span class="line"><span class="cl">    proxy_set_header   X-Real-IP $remote_addr;
</span></span><span class="line"><span class="cl">    proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
</span></span><span class="line"><span class="cl">    proxy_set_header   X-Forwarded-Proto $scheme;
</span></span><span class="line"><span class="cl">    proxy_pass  http://127.0.0.1:660;  
</span></span><span class="line"><span class="cl">}
</span></span><span class="line"><span class="cl">###  静态html目录
</span></span><span class="line"><span class="cl">location / {
</span></span><span class="line"><span class="cl">    proxy_redirect     off;
</span></span><span class="line"><span class="cl">    proxy_set_header   Host      $host;
</span></span><span class="line"><span class="cl">    proxy_set_header   X-Real-IP $remote_addr;
</span></span><span class="line"><span class="cl">    proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
</span></span><span class="line"><span class="cl">    proxy_set_header   X-Forwarded-Proto $scheme; 
</span></span><span class="line"><span class="cl">    ## 存在q查询参数,使用动态服务.也支持FlexSearch解析public/search-data.json
</span></span><span class="line"><span class="cl">    if ($arg_q) { 
</span></span><span class="line"><span class="cl">       proxy_pass  http://127.0.0.1:660;  
</span></span><span class="line"><span class="cl">       break;
</span></span><span class="line"><span class="cl">    }
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl">    #### 开启gzip静态压缩
</span></span><span class="line"><span class="cl">    #gzip_static on;
</span></span><span class="line"><span class="cl">    ### Nginx 1.26+ 不需要再进行302重定向到目录下的index.html,gzip_static也会生效.这段配置留作记录.
</span></span><span class="line"><span class="cl">    ##if ( -d $request_filename ) {
</span></span><span class="line"><span class="cl">        ## 不是 / 结尾
</span></span><span class="line"><span class="cl">    ##    rewrite [^\/]$ $uri/index.html redirect;
</span></span><span class="line"><span class="cl">        ##以 / 结尾的
</span></span><span class="line"><span class="cl">    ##    rewrite ^(.*) ${uri}index.html redirect;      
</span></span><span class="line"><span class="cl">    ##}
</span></span><span class="line"><span class="cl">    
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl">    
</span></span><span class="line"><span class="cl">    root   /data/gpress/gpressdatadir/statichtml;
</span></span><span class="line"><span class="cl">    index  index.html index.htm;
</span></span><span class="line"><span class="cl">}
</span></span><span class="line"><span class="cl">
</span></span></code></pre></div><h1 id="阿里云计算巢">阿里云计算巢</h1>
<p><a href="https://computenest.console.aliyun.com/service/instance/create/cn-hangzhou?type=user&amp;ServiceId=service-d4000c9b22c54e5cbffe">点击部署gpress到阿里云计算巢</a>,也可以单独购买阿里云最低配服务器,进行部署.选择<code>张家口机房</code>,规格<code>ecs.t6-c4m1.large</code>,配置<code>2核CPU 0.5G内存 20G高效云盘 RockyLinux9 按使用流量-带宽峰值80M</code>,一年100元,五年200元左右.</p>
',
                        '',
                        '<ul>
<li>
<a href="#%E4%BB%8B%E7%BB%8D">介绍</a></li>
<li>
<a href="#%E5%BC%80%E5%8F%91%E7%8E%AF%E5%A2%83">开发环境</a></li>
<li>
<a href="#%E9%9D%99%E6%80%81%E5%8C%96">静态化</a></li>
<li>
<a href="#%E9%98%BF%E9%87%8C%E4%BA%91%E8%AE%A1%E7%AE%97%E5%B7%A2">阿里云计算巢</a></li>
</ul>
',
                        '',
                        '',
                        '',
                        ',web,',
                        'Web',
                        'web',
                        '',
                        '',
                        '',
                        '',
                        'gpress',
                        '',
                        'gpress'
                    );
