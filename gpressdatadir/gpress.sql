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
INSERT INTO site (status,sortNo,createUser,updateTime,createTime,footer,favicon,logo,siteThemeWEIXIN,themeWAP,themePC,theme,description,keyword,domain,name,title,id)VALUES (1,1,NULL,NULL,NULL,'<div class="copyright"><span class="copyright-year">&copy; 2008 - 2024<span class="author">jiagou.com 版权所有 <a href=''https://beian.miit.gov.cn'' target=''_blank''>豫ICP备xxxxx号</a>   <a href=''http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=xxxx''  target=''_blank''><img src=''/public/gongan.png''>豫公网安备xxxxx号</a></span></span></div>','/public/favicon.png','/public/logo.png','default','default','default','default','Web3内容平台,Hertz + Go template + FTS5全文检索,支持以太坊和百度超级链,兼容hugo生态,使用Wasm扩展插件,只需200M内存','gpress,web3,hugo,wordpress,以太坊,百度超级链','jiagou.com','架构','jiagou','gpress_site');

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

INSERT INTO content (status,sortNo,createUser,updateTime,createTime,thumbnail,markdown,content,summary,toc,tag,author,templateFile,comCode,categoryName,categoryID,subtitle,hrefURL,description,keyword,title,moduleID,id)
            VALUES (0,100,NULL,'2024-04-06 20:28:38','2023-06-27 22:43:53',NULL,'本站服务器配置:ecs.t6-c4m1.large,2核CPU,512M内存,20G高效云盘,RockyLinux 9 .  
使用gpress迁移了hugo的even主题和markdown文件.  

我所见识过的一切都将消失一空,就如眼泪消逝在雨中......  
不妨大胆一些,大胆一些......  

小项目:  
* [springrain](https://gitee.com/chunanyong/springrain)
* [zorm](https://gitee.com/chunanyong/zorm)
* [gpress](https://gitee.com/gpress/gpress)
* [gowe](https://gitee.com/chunanyong/gowe)','<p>本站服务器配置:ecs.t6-c4m1.large,2核CPU,512M内存,20G高效云盘,RockyLinux 9 .<br>
使用gpress迁移了hugo的even主题和markdown文件.</p>
<p>我所见识过的一切都将消失一空,就如眼泪消逝在雨中......<br>
不妨大胆一些,大胆一些......</p>
<p>小项目:</p>
<ul>
<li><a href="https://gitee.com/chunanyong/springrain">springrain</a></li>
<li><a href="https://gitee.com/chunanyong/zorm">zorm</a></li>
<li><a href="https://gitee.com/gpress/gpress">gpress</a></li>
<li><a href="https://gitee.com/chunanyong/gowe">gowe</a></li>
</ul>
','本站服务器配置:1核CPU,512M内存,20G硬盘,AnolisOS(ANCK).使用hugo和even模板,编译成静态文件,Nginx作为WEB服务器.我所见识过的一切都将消失一空,就如眼泪消逝在雨中......	不妨大胆一些,大胆一些......','',NULL,'springrain',NULL,',about,','About','about',NULL,NULL,NULL,NULL,'about',NULL,'about'
                    ),
                    (1,101,'','2024-04-06 20:39:20','2024-04-06 20:31:32','','## 介绍  
Web3内容平台,Hertz + Go template + FTS5全文检索,支持以太坊和百度超级链,兼容hugo生态,使用Wasm扩展插件,只需200M内存  
    
**作为静态站点：** gpress生成的静态文件和Hugo一致,也可以简单认为gpress是Hugo的后台管理,兼容Hugo主题生态,已迁移多款Hugo主题:[even](gitee.com/gpress/gpress/tree/master/gpressdatadir/template/theme/default)、[doks](gitee.com/gpress/gpress-doks)、[book](gitee.com/gpress/gpress-book)、[geekdoc](gitee.com/gpress/gpress-geekdoc)......   
**作为动态站点：** gpress功能简单,只有7个菜单,5张表,5000行代码,使用SQLite,一键启动,只需200M内存,支持全文检索......  
**作为Web3：** gpress已支持以太坊和百度超级链账户体系,会基于Wasm持续迭代去中心功能,让数据自由一点点......  
**作为后浪：** 相对于Hugo、WordPress等优秀的内容平台,gpress还有很多不足,功能简单而又稚嫩......  
**帮助文档：** [点击查看帮助文档](/public/doc/index.md)   
 
个人博客 [jiagou.com](jiagou.com/post/about) 使用gpress搭建,搜索和后台管理是动态,其他是静态页面.  
![](/public/index.png "")

## 开发环境  
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

## 静态化
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

    #### gzip 静态压缩配置 开始####
    #gzip_static on;
    ## 请求的是个目录,302重定向到 目录下的 index.html
    #if ( -d $request_filename ) {
        ## 不是 / 结尾
    #    rewrite [^\/]$ $uri/index.html redirect;
        ##以 / 结尾的
    #    rewrite ^(.*) ${uri}index.html redirect;      
    #}
    #### gzip 静态压缩配置 结束####

    
    root   /data/gpress/gpressdatadir/statichtml;
    index  index.html index.htm;
}

```  
## 阿里云计算巢
[点击部署gpress到阿里云计算巢](https://computenest.console.aliyun.com/service/instance/create/cn-hangzhou?type=user&ServiceId=service-d4000c9b22c54e5cbffe),也可以独立购买阿里云的服务器,进行部署.选择```张家口机房```最低配置的 ```ecs.t6-c4m1.large``` 规格```2核CPU 0.5G内存 20G高效云盘 RockyLinux9 按使用流量-带宽峰值80M```就够用了,一年100元左右,性价比高.  

## 表结构  
ID默认使用时间戳(23位)+随机数(9位),全局唯一.  
建表语句```gpressdatadir/gpress.sql```          

### 配置(表名:config)
安装时会读取```gpressdatadir/install_config.json```

| 列名  | 类型        | 说明         |  备注       | 
| ----------- | ----------- | ----------- | ----------- |
| id          | string      | 主键        |gpress_config |
| basePath    | string      | 基础路径    |  默认 /      |
| jwtSecret   | string      | jwt密钥     | 随机生成     |
| jwttokenKey | string      | jwt的key    |  默认 jwttoken  |
| serverPort  | string      | IP:端口     |  默认 :660  |
| timeout     | int         | jwt超时时间秒|  默认 7200  |
| maxRequestBodySize| int   | 最大请求|  默认 20M  |
| proxy       | string      | http代理地址 |             |
| createTime  | string      | 创建时间     |  2006-01-02 15:04:05  |
| updateTime  | string      | 更新时间     |  2006-01-02 15:04:05  |
| createUser  | string      | 创建人       |  初始化 system  |
| sortNo      | int         | 排序         |  正序  |
| status      | int         | 状态     |  链接访问(0),公开(1),私密(2)  |

### 用户(表名:user)
后台只有一个用户.

| 列名  | 类型         | 说明        |  备注       | 
| ----------- | ----------- | ----------- | ----------- |
| id          | string      | 主键        | gpress_admin |
| account     | string      | 登录名称    |  默认admin  |
| passWord    | string      | 密码        |    -  |
| userName    | string      | 说明        |    -  |
| createTime  | string      | 创建时间     |  2006-01-02 15:04:05  |
| updateTime  | string      | 更新时间     |  2006-01-02 15:04:05  |
| createUser  | string      | 创建人       |  初始化 system  |
| sortNo      | int         | 排序         |  正序  |
| status      | int         | 状态     |  链接访问(0),公开(1),私密(2)  |

### 站点信息(site)
站点的信息,例如 title,logo,keywords,description等

| 列名    | 类型         | 说明    |  备注       | 
| ----------- | ----------- | ----------- | ----------- |
| id          | string      | 主键        |gpress_site  |
| title       | string      | 站点名称     |     -  |
| keyword     | string      | 关键字       |     -  |
| description | string      | 站点描述    |     -  |
| theme       | string      | 默认主题     | 默认使用default  |
| themePC     | string      | PC主题      | 先从cookie获取,如果没有从Header头取值,写入cookie,默认使用default  |
| themeWAP    | string      | 手机主题    | 先从cookie获取,如果没有从Header头取值,写入cookie,默认使用default  |
| themeWEIXIN | string      | 微信主题    | 先从cookie获取,如果没有从Header头取值,写入cookie,默认使用default  |
| logo        | string      | logo       |     -  |
| favicon     | string      | Favicon    |     -  |
| createTime  | string      | 创建时间     |  2006-01-02 15:04:05  |
| updateTime  | string      | 更新时间     |  2006-01-02 15:04:05  |
| createUser  | string      | 创建人       |  初始化 system  |
| sortNo      | int         | 排序         |  正序  |
| status      | int         | 状态     |  链接访问(0),公开(1),私密(2)  |

### 导航菜单(表名:category)
| 列名    | 类型         | 说明    |  备注       | 
| ----------- | ----------- | ----------- | ----------- |
| id          | string      | 主键        |    -  |
| name        | string      | 导航名称     |    -  |
| hrefURL     | string      | 跳转路径     |    -  |
| hrefTarget  | string      | 跳转方式     | _self,_blank,_parent,_top|
| pid         | string      | 父导航ID     | 父导航ID  |
| moduleID    | string      | module表ID   |  导航菜单下的文章默认使用的模型字段 |
| comCode     | string      | 逗号隔开的全路径 | 逗号隔开的全路径  |
| templateFile  | string      | 模板文件       | 当前导航页的模板  |
| childTemplateFile  | string | 子主题模板文件  | 子页面默认使用的模板,子页面如果不设置,默认使用这个模板 |
| createTime  | string      | 创建时间     |  2006-01-02 15:04:05  |
| updateTime  | string      | 更新时间     |  2006-01-02 15:04:05  |
| createUser  | string      | 创建人       |  初始化 system  |
| sortNo      | int         | 排序         |  正序  |
| status      | int         | 状态     |  链接访问(0),公开(1),私密(2)  |

### 文章内容(表名:content)
| 列名  | 类型        | 说明        | 是否分词 |  备注                  | 
| ----------- | ----------- | ----------- | ------- | ---------------------- |
| id          | string      | 主键         | 否      |    -                   |
| moduleID    | string      | 模型ID       | 否      |  文章使用的模型字段     |
| title       | string      | 文章标题     | 是      |    使用 jieba 分词器    |
| keyword     | string      | 关键字       | 是      |    使用 jieba 分词器    |
| description | string      | 站点描述     | 是      |    使用 jieba 分词器    |
| hrefURL     | string      | 自身页面路径 | 否      |    -                    |
| subtitle    | string      | 副标题       | 是      |      使用 jieba 分词器  |
| author      | string      | 作者         | 是      |      使用 jieba 分词器  |
| tag         | string      | 标签         | 是      |      使用 jieba 分词器  |
| toc         | string      | 目录         | 是      |      使用 jieba 分词器  |
| summary     | string      | 摘要         | 是      |      使用 jieba 分词器  |
| categoryName| string      | 导航菜单,逗号(,)隔开| 是| 使用 jieba 分词器.      |
| categoryID  | string      | 导航ID       | 否      | -                       |
| comCode     | string      | 逗号隔开的全路径 | 逗号隔开的全路径  |
| templateFile  | string      | 模板文件       | 否      | 模板                    |
| content     | string      | 文章内容     | 否      |                         |
| markdown    | string      | Markdown内容 | 否      |                         |
| thumbnail   | string      | 封面图       | 否      |                         |
| createTime  | string      | 创建时间     | -       |  2006-01-02 15:04:05    |
| updateTime  | string      | 更新时间     | -       |  2006-01-02 15:04:05    |
| createUser  | string      | 创建人       | -       |  初始化 system          |
| sortNo      | int         | 排序         | -       |  正序                   |
| status      | int         | 状态     | -       |  链接访问(0),公开(1),私密(2)  |','<h2 id="介绍">介绍</h2>
<p>Web3内容平台,Hertz + Go template + FTS5全文检索,支持以太坊和百度超级链,兼容hugo生态,使用Wasm扩展插件,只需200M内存</p>
<p><strong>作为静态站点：</strong> gpress生成的静态文件和Hugo一致,也可以简单认为gpress是Hugo的后台管理,兼容Hugo主题生态,已迁移多款Hugo主题:<a href="gitee.com/gpress/gpress/tree/master/gpressdatadir/template/theme/default">even</a>、<a href="gitee.com/gpress/gpress-doks">doks</a>、<a href="gitee.com/gpress/gpress-book">book</a>、<a href="gitee.com/gpress/gpress-geekdoc">geekdoc</a>......<br>
<strong>作为动态站点：</strong> gpress功能简单,只有7个菜单,5张表,5000行代码,使用SQLite,一键启动,只需200M内存,支持全文检索......<br>
<strong>作为Web3：</strong> gpress已支持以太坊和百度超级链账户体系,会基于Wasm持续迭代去中心功能,让数据自由一点点......<br>
<strong>作为后浪：</strong> 相对于Hugo、WordPress等优秀的内容平台,gpress还有很多不足,功能简单而又稚嫩......<br>
<strong>帮助文档：</strong> <a href="/public/doc/index.md">点击查看帮助文档</a></p>
<p>个人博客 <a href="jiagou.com/post/about">jiagou.com</a> 使用gpress搭建,搜索和后台管理是动态,其他是静态页面.<br>
<img src="/public/index.png" alt="" title=""></p>
<h2 id="开发环境">开发环境</h2>
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
</span></span></code></pre></div><h2 id="静态化">静态化</h2>
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
</span></span><span class="line"><span class="cl">    #### gzip 静态压缩配置 开始####
</span></span><span class="line"><span class="cl">    #gzip_static on;
</span></span><span class="line"><span class="cl">    ## 请求的是个目录,302重定向到 目录下的 index.html
</span></span><span class="line"><span class="cl">    #if ( -d $request_filename ) {
</span></span><span class="line"><span class="cl">        ## 不是 / 结尾
</span></span><span class="line"><span class="cl">    #    rewrite [^\/]$ $uri/index.html redirect;
</span></span><span class="line"><span class="cl">        ##以 / 结尾的
</span></span><span class="line"><span class="cl">    #    rewrite ^(.*) ${uri}index.html redirect;      
</span></span><span class="line"><span class="cl">    #}
</span></span><span class="line"><span class="cl">    #### gzip 静态压缩配置 结束####
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl">    
</span></span><span class="line"><span class="cl">    root   /data/gpress/gpressdatadir/statichtml;
</span></span><span class="line"><span class="cl">    index  index.html index.htm;
</span></span><span class="line"><span class="cl">}
</span></span><span class="line"><span class="cl">
</span></span></code></pre></div><h2 id="阿里云计算巢">阿里云计算巢</h2>
<p><a href="https://computenest.console.aliyun.com/service/instance/create/cn-hangzhou?type=user&amp;ServiceId=service-d4000c9b22c54e5cbffe">点击部署gpress到阿里云计算巢</a>,也可以独立购买阿里云的服务器,进行部署.选择<code>张家口机房</code>最低配置的 <code>ecs.t6-c4m1.large</code> 规格<code>2核CPU 0.5G内存 20G高效云盘 RockyLinux9 按使用流量-带宽峰值80M</code>就够用了,一年100元左右,性价比高.</p>
<h2 id="表结构">表结构</h2>
<p>ID默认使用时间戳(23位)+随机数(9位),全局唯一.<br>
建表语句<code>gpressdatadir/gpress.sql</code></p>
<h3 id="配置(表名:config)">配置(表名:config)</h3>
<p>安装时会读取<code>gpressdatadir/install_config.json</code></p>
<table>
<thead>
<tr>
<th>列名</th>
<th>类型</th>
<th>说明</th>
<th>备注</th>
</tr>
</thead>
<tbody>
<tr>
<td>id</td>
<td>string</td>
<td>主键</td>
<td>gpress_config</td>
</tr>
<tr>
<td>basePath</td>
<td>string</td>
<td>基础路径</td>
<td>默认 /</td>
</tr>
<tr>
<td>jwtSecret</td>
<td>string</td>
<td>jwt密钥</td>
<td>随机生成</td>
</tr>
<tr>
<td>jwttokenKey</td>
<td>string</td>
<td>jwt的key</td>
<td>默认 jwttoken</td>
</tr>
<tr>
<td>serverPort</td>
<td>string</td>
<td>IP:端口</td>
<td>默认 :660</td>
</tr>
<tr>
<td>timeout</td>
<td>int</td>
<td>jwt超时时间秒</td>
<td>默认 7200</td>
</tr>
<tr>
<td>maxRequestBodySize</td>
<td>int</td>
<td>最大请求</td>
<td>默认 20M</td>
</tr>
<tr>
<td>proxy</td>
<td>string</td>
<td>http代理地址</td>
<td></td>
</tr>
<tr>
<td>createTime</td>
<td>string</td>
<td>创建时间</td>
<td>2006-01-02 15:04:05</td>
</tr>
<tr>
<td>updateTime</td>
<td>string</td>
<td>更新时间</td>
<td>2006-01-02 15:04:05</td>
</tr>
<tr>
<td>createUser</td>
<td>string</td>
<td>创建人</td>
<td>初始化 system</td>
</tr>
<tr>
<td>sortNo</td>
<td>int</td>
<td>排序</td>
<td>正序</td>
</tr>
<tr>
<td>status</td>
<td>int</td>
<td>状态</td>
<td>链接访问(0),公开(1),私密(2)</td>
</tr>
</tbody>
</table>
<h3 id="用户(表名:user)">用户(表名:user)</h3>
<p>后台只有一个用户.</p>
<table>
<thead>
<tr>
<th>列名</th>
<th>类型</th>
<th>说明</th>
<th>备注</th>
</tr>
</thead>
<tbody>
<tr>
<td>id</td>
<td>string</td>
<td>主键</td>
<td>gpress_admin</td>
</tr>
<tr>
<td>account</td>
<td>string</td>
<td>登录名称</td>
<td>默认admin</td>
</tr>
<tr>
<td>passWord</td>
<td>string</td>
<td>密码</td>
<td>-</td>
</tr>
<tr>
<td>userName</td>
<td>string</td>
<td>说明</td>
<td>-</td>
</tr>
<tr>
<td>createTime</td>
<td>string</td>
<td>创建时间</td>
<td>2006-01-02 15:04:05</td>
</tr>
<tr>
<td>updateTime</td>
<td>string</td>
<td>更新时间</td>
<td>2006-01-02 15:04:05</td>
</tr>
<tr>
<td>createUser</td>
<td>string</td>
<td>创建人</td>
<td>初始化 system</td>
</tr>
<tr>
<td>sortNo</td>
<td>int</td>
<td>排序</td>
<td>正序</td>
</tr>
<tr>
<td>status</td>
<td>int</td>
<td>状态</td>
<td>链接访问(0),公开(1),私密(2)</td>
</tr>
</tbody>
</table>
<h3 id="站点信息(site)">站点信息(site)</h3>
<p>站点的信息,例如 title,logo,keywords,description等</p>
<table>
<thead>
<tr>
<th>列名</th>
<th>类型</th>
<th>说明</th>
<th>备注</th>
</tr>
</thead>
<tbody>
<tr>
<td>id</td>
<td>string</td>
<td>主键</td>
<td>gpress_site</td>
</tr>
<tr>
<td>title</td>
<td>string</td>
<td>站点名称</td>
<td>-</td>
</tr>
<tr>
<td>keyword</td>
<td>string</td>
<td>关键字</td>
<td>-</td>
</tr>
<tr>
<td>description</td>
<td>string</td>
<td>站点描述</td>
<td>-</td>
</tr>
<tr>
<td>theme</td>
<td>string</td>
<td>默认主题</td>
<td>默认使用default</td>
</tr>
<tr>
<td>themePC</td>
<td>string</td>
<td>PC主题</td>
<td>先从cookie获取,如果没有从Header头取值,写入cookie,默认使用default</td>
</tr>
<tr>
<td>themeWAP</td>
<td>string</td>
<td>手机主题</td>
<td>先从cookie获取,如果没有从Header头取值,写入cookie,默认使用default</td>
</tr>
<tr>
<td>themeWEIXIN</td>
<td>string</td>
<td>微信主题</td>
<td>先从cookie获取,如果没有从Header头取值,写入cookie,默认使用default</td>
</tr>
<tr>
<td>logo</td>
<td>string</td>
<td>logo</td>
<td>-</td>
</tr>
<tr>
<td>favicon</td>
<td>string</td>
<td>Favicon</td>
<td>-</td>
</tr>
<tr>
<td>createTime</td>
<td>string</td>
<td>创建时间</td>
<td>2006-01-02 15:04:05</td>
</tr>
<tr>
<td>updateTime</td>
<td>string</td>
<td>更新时间</td>
<td>2006-01-02 15:04:05</td>
</tr>
<tr>
<td>createUser</td>
<td>string</td>
<td>创建人</td>
<td>初始化 system</td>
</tr>
<tr>
<td>sortNo</td>
<td>int</td>
<td>排序</td>
<td>正序</td>
</tr>
<tr>
<td>status</td>
<td>int</td>
<td>状态</td>
<td>链接访问(0),公开(1),私密(2)</td>
</tr>
</tbody>
</table>
<h3 id="导航菜单(表名:category)">导航菜单(表名:category)</h3>
<table>
<thead>
<tr>
<th>列名</th>
<th>类型</th>
<th>说明</th>
<th>备注</th>
</tr>
</thead>
<tbody>
<tr>
<td>id</td>
<td>string</td>
<td>主键</td>
<td>-</td>
</tr>
<tr>
<td>name</td>
<td>string</td>
<td>导航名称</td>
<td>-</td>
</tr>
<tr>
<td>hrefURL</td>
<td>string</td>
<td>跳转路径</td>
<td>-</td>
</tr>
<tr>
<td>hrefTarget</td>
<td>string</td>
<td>跳转方式</td>
<td>_self,_blank,_parent,_top</td>
</tr>
<tr>
<td>pid</td>
<td>string</td>
<td>父导航ID</td>
<td>父导航ID</td>
</tr>
<tr>
<td>moduleID</td>
<td>string</td>
<td>module表ID</td>
<td>导航菜单下的文章默认使用的模型字段</td>
</tr>
<tr>
<td>comCode</td>
<td>string</td>
<td>逗号隔开的全路径</td>
<td>逗号隔开的全路径</td>
</tr>
<tr>
<td>templateFile</td>
<td>string</td>
<td>模板文件</td>
<td>当前导航页的模板</td>
</tr>
<tr>
<td>childTemplateFile</td>
<td>string</td>
<td>子主题模板文件</td>
<td>子页面默认使用的模板,子页面如果不设置,默认使用这个模板</td>
</tr>
<tr>
<td>createTime</td>
<td>string</td>
<td>创建时间</td>
<td>2006-01-02 15:04:05</td>
</tr>
<tr>
<td>updateTime</td>
<td>string</td>
<td>更新时间</td>
<td>2006-01-02 15:04:05</td>
</tr>
<tr>
<td>createUser</td>
<td>string</td>
<td>创建人</td>
<td>初始化 system</td>
</tr>
<tr>
<td>sortNo</td>
<td>int</td>
<td>排序</td>
<td>正序</td>
</tr>
<tr>
<td>status</td>
<td>int</td>
<td>状态</td>
<td>链接访问(0),公开(1),私密(2)</td>
</tr>
</tbody>
</table>
<h3 id="文章内容(表名:content)">文章内容(表名:content)</h3>
<table>
<thead>
<tr>
<th>列名</th>
<th>类型</th>
<th>说明</th>
<th>是否分词</th>
<th>备注</th>
</tr>
</thead>
<tbody>
<tr>
<td>id</td>
<td>string</td>
<td>主键</td>
<td>否</td>
<td>-</td>
</tr>
<tr>
<td>moduleID</td>
<td>string</td>
<td>模型ID</td>
<td>否</td>
<td>文章使用的模型字段</td>
</tr>
<tr>
<td>title</td>
<td>string</td>
<td>文章标题</td>
<td>是</td>
<td>使用 jieba 分词器</td>
</tr>
<tr>
<td>keyword</td>
<td>string</td>
<td>关键字</td>
<td>是</td>
<td>使用 jieba 分词器</td>
</tr>
<tr>
<td>description</td>
<td>string</td>
<td>站点描述</td>
<td>是</td>
<td>使用 jieba 分词器</td>
</tr>
<tr>
<td>hrefURL</td>
<td>string</td>
<td>自身页面路径</td>
<td>否</td>
<td>-</td>
</tr>
<tr>
<td>subtitle</td>
<td>string</td>
<td>副标题</td>
<td>是</td>
<td>使用 jieba 分词器</td>
</tr>
<tr>
<td>author</td>
<td>string</td>
<td>作者</td>
<td>是</td>
<td>使用 jieba 分词器</td>
</tr>
<tr>
<td>tag</td>
<td>string</td>
<td>标签</td>
<td>是</td>
<td>使用 jieba 分词器</td>
</tr>
<tr>
<td>toc</td>
<td>string</td>
<td>目录</td>
<td>是</td>
<td>使用 jieba 分词器</td>
</tr>
<tr>
<td>summary</td>
<td>string</td>
<td>摘要</td>
<td>是</td>
<td>使用 jieba 分词器</td>
</tr>
<tr>
<td>categoryName</td>
<td>string</td>
<td>导航菜单,逗号(,)隔开</td>
<td>是</td>
<td>使用 jieba 分词器.</td>
</tr>
<tr>
<td>categoryID</td>
<td>string</td>
<td>导航ID</td>
<td>否</td>
<td>-</td>
</tr>
<tr>
<td>comCode</td>
<td>string</td>
<td>逗号隔开的全路径</td>
<td>逗号隔开的全路径</td>
<td></td>
</tr>
<tr>
<td>templateFile</td>
<td>string</td>
<td>模板文件</td>
<td>否</td>
<td>模板</td>
</tr>
<tr>
<td>content</td>
<td>string</td>
<td>文章内容</td>
<td>否</td>
<td></td>
</tr>
<tr>
<td>markdown</td>
<td>string</td>
<td>Markdown内容</td>
<td>否</td>
<td></td>
</tr>
<tr>
<td>thumbnail</td>
<td>string</td>
<td>封面图</td>
<td>否</td>
<td></td>
</tr>
<tr>
<td>createTime</td>
<td>string</td>
<td>创建时间</td>
<td>-</td>
<td>2006-01-02 15:04:05</td>
</tr>
<tr>
<td>updateTime</td>
<td>string</td>
<td>更新时间</td>
<td>-</td>
<td>2006-01-02 15:04:05</td>
</tr>
<tr>
<td>createUser</td>
<td>string</td>
<td>创建人</td>
<td>-</td>
<td>初始化 system</td>
</tr>
<tr>
<td>sortNo</td>
<td>int</td>
<td>排序</td>
<td>-</td>
<td>正序</td>
</tr>
<tr>
<td>status</td>
<td>int</td>
<td>状态</td>
<td>-</td>
<td>链接访问(0),公开(1),私密(2)</td>
</tr>
</tbody>
</table>
','','<ul>
<li>
<ul>
<li>
<a href="#%E4%BB%8B%E7%BB%8D">介绍</a></li>
<li>
<a href="#%E5%BC%80%E5%8F%91%E7%8E%AF%E5%A2%83">开发环境</a></li>
<li>
<a href="#%E9%9D%99%E6%80%81%E5%8C%96">静态化</a></li>
<li>
<a href="#%E9%98%BF%E9%87%8C%E4%BA%91%E8%AE%A1%E7%AE%97%E5%B7%A2">阿里云计算巢</a></li>
<li>
<a href="#%E8%A1%A8%E7%BB%93%E6%9E%84">表结构</a><ul>
<li>
<a href="#%E9%85%8D%E7%BD%AE(%E8%A1%A8%E5%90%8D:config)">配置(表名:config)</a></li>
<li>
<a href="#%E7%94%A8%E6%88%B7(%E8%A1%A8%E5%90%8D:user)">用户(表名:user)</a></li>
<li>
<a href="#%E7%AB%99%E7%82%B9%E4%BF%A1%E6%81%AF(site)">站点信息(site)</a></li>
<li>
<a href="#%E5%AF%BC%E8%88%AA%E8%8F%9C%E5%8D%95(%E8%A1%A8%E5%90%8D:category)">导航菜单(表名:category)</a></li>
<li>
<a href="#%E6%96%87%E7%AB%A0%E5%86%85%E5%AE%B9(%E8%A1%A8%E5%90%8D:content)">文章内容(表名:content)</a></li>
</ul>
</li>
</ul>
</li>
</ul>
','','','',',web,','Web','web','','','','','gpress','','gpress'
);
