<img src="gpressdatadir/public/gpress-logo.png" height="150px" />    

<a href="./README.md">English</a> | <a href="./README.zh-CN.md">简体中文</a>         
## Introduction  
Web3 content platform, Hertz + Go template + FTS5 full-text search, supports Ethereum and Baidu Super Chain, compatible with Hugo and WordPress ecosystems, uses Wasm extension plugins, and requires only 200M memory.  

**As a static site:** The static files generated by gpress are consistent with Hugo, and gpress can also be simply considered as the backend management of Hugo, compatible with Hugo theme ecosystems. Multiple Hugo themes have been migrated: [even](https://gitee.com/gpress/gpress/tree/master/gpressdatadir/template/theme/default), [doks](https://gitee.com/gpress/gpress-doks), [book](https://gitee.com/gpress/gpress-book), [geekdoc](https://gitee.com/gpress/gpress-geekdoc)...  
**As a dynamic site:** gpress is simple in functionality, with only 7 menus, 5 tables, and 5000 lines of code. It uses SQLite, starts with one click, and requires only 200M memory, supporting full-text search. It is compatible with WordPress theme ecosystems, and multiple WordPress themes have been migrated: [generatepress](https://gitee.com/gpress/wp-generatepress), [astra](https://gitee.com/gpress/wp-astra)...  
**As Web3:** gpress already supports Ethereum and Baidu Super Chain account systems and will continue to iterate decentralized features based on Wasm, allowing data to be a bit freer...  
**As a newcomer:** Compared to excellent content platforms like Hugo and WordPress, gpress still has many shortcomings, being simple and immature in functionality...  
**Documentation:** [Click to view the documentation](./gpressdatadir/public/doc/index.md)  

The personal blog [jiagou.com](https://jiagou.com) is built using gpress, with dynamic search and backend management, while the rest are static pages.  
<img src="gpressdatadir/public/index.png" width="600px">

## Development Environment  
gpress uses ```https://github.com/wangfenjin/simple``` as the FTS5 full-text search extension. The compiled libsimple file is placed in the ```gpressdatadir/fts5``` directory. If gpress fails to start and reports an error connecting to the database, please check if the libsimple file is correct. If you need to recompile libsimple, please refer to https://github.com/wangfenjin/simple.  

The default port is 660, and the backend management address is http://127.0.0.1:660/admin/login.  
First, unzip ```gpressdatadir/dict.zip```.  
Run ```go run --tags "fts5" .```.  
Package: ```go build --tags "fts5" -ldflags "-w -s"```.  

The development environment requires CGO compilation configuration. Set ```set CGO_ENABLED=1```, download [mingw64](https://github.com/niXman/mingw-builds-binaries/releases) and [cmake](https://cmake.org/download/), and configure the bin to the environment variables. Note to rename ```mingw64/bin/mingw32-make.exe``` to ```make.exe```.  
Modify vscode's launch.json to add ``` ,"buildFlags": "--tags=fts5" ``` for debugging fts5.  
Test needs to be done manually: ```go test -timeout 30s --tags "fts5" -run ^TestReadmks$ gitee.com/gpress/gpress```.  
Package: ```go build --tags "fts5" -ldflags "-w -s"```.  
When recompiling simple, it is recommended to use the precompiled version from ```https://github.com/wangfenjin/simple```.  
Note to modify the Windows compilation script, remove the ```libgcc_s_seh-1.dll``` and ```libstdc++-6.dll``` dependencies for mingw64 compilation, and turn off ```BUILD_TEST_EXAMPLE``` as there are conflicts.  
```bat
rmdir /q /s build
mkdir build && cd build
cmake .. -G "Unix Makefiles" -DBUILD_TEST_EXAMPLE=OFF -DCMAKE_INSTALL_PREFIX=release -DCMAKE_CXX_FLAGS="-static-libgcc -static-libstdc++" -DCMAKE_EXE_LINKER_FLAGS="-Wl,-Bstatic -lstdc++ -lpthread -Wl,-Bdynamic"
make && make install
```

## Staticization  
The backend ```Refresh Site``` function will generate static HTML files to the ```statichtml``` directory, along with ```gzip_static``` files. You need to copy the ```css, js, image``` of the currently used theme and the ```gpressdatadir/public``` directory to the ```statichtml``` directory, or use Nginx reverse proxy to specify the directory without copying files.  
Nginx configuration example:
```conf
### CSS files of the current theme (default)
location ~ ^/css/ {
    #gzip_static on;
    root /data/gpress/gpressdatadir/template/theme/default;  
}
### JS files of the current theme (default)
location ~ ^/js/ {
    #gzip_static on;
    root /data/gpress/gpressdatadir/template/theme/default;  
}
### Image files of the current theme (default)
location ~ ^/image/ {
    root /data/gpress/gpressdatadir/template/theme/default;  
}
### search-data.json FlexSearch JSON data
location ~ ^/public/search-data.json {
    #gzip_static on;
    root /data/gpress/gpressdatadir;  
}
### Public files
location ~ ^/public/ {
    root /data/gpress/gpressdatadir;  
}
    
### Admin backend management, request dynamic service
location ~ ^/admin/ {
    proxy_redirect     off;
    proxy_set_header   Host      $host;
    proxy_set_header   X-Real-IP $remote_addr;
    proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
    proxy_set_header   X-Forwarded-Proto $scheme;
    proxy_pass  http://127.0.0.1:660;  
}
### Static HTML directory
location / {
    proxy_redirect     off;
    proxy_set_header   Host      $host;
    proxy_set_header   X-Real-IP $remote_addr;
    proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
    proxy_set_header   X-Forwarded-Proto $scheme; 
    ## If there is a q query parameter, use the dynamic service. Also supports FlexSearch parsing public/search-data.json
    if ($arg_q) { 
       proxy_pass  http://127.0.0.1:660;  
       break;
    }

    ### Enable gzip static compression
    #gzip_static on;

    ### Nginx 1.26+ does not need to 302 redirect to the index.html under the directory, gzip_static will also take effect. This configuration is kept for record.
    ##if ( -d $request_filename ) {
        ## Not ending with /
    ##    rewrite [^\/]$ $uri/index.html redirect;
        ## Ending with /
    ##    rewrite ^(.*) ${uri}index.html redirect;      
    ##}
    
    ### Static file directory of the current theme (default)
    root   /data/gpress/gpressdatadir/statichtml/default;
    
    ### if directive may conflict with try_files directive, causing try_files to be invalid
    ## Avoid directory 301 redirect, e.g., /about will 301 to /about/           
    try_files $uri $uri/index.html;
    
    index  index.html index.htm;
}

```  

## Backend Management Supports English
The gpress backend management currently supports both Chinese and English, with the capability to extend to other languages. Language files are located in ```gpressdatadir/locales```. By default, the system uses Chinese (```zh-CN```) upon initial installation. If English is preferred, you can modify the ```"locale":"zh-CN"``` to ```"locale":"en-US"``` in the ```gpressdatadir/install_config.json``` file before installation. Alternatively, after successful installation, you can change the ```Language``` setting to ```English``` in the ```Settings``` and restart the system to apply the changes.

## Alibaba Cloud Computing Nest  
[Click to deploy gpress to Alibaba Cloud Computing Nest](https://computenest.console.aliyun.com/service/instance/create/cn-hangzhou?type=user&ServiceId=service-d4000c9b22c54e5cbffe), or you can purchase the lowest configuration Alibaba Cloud server separately for deployment. Choose ```Zhangjiakou Data Center```, specification ```ecs.t6-c4m1.large```, configuration ```2-core CPU 0.5G memory 20G efficient cloud disk RockyLinux9 pay-by-traffic-bandwidth peak 80M```, costing 100 yuan per year, around 200 yuan for five years.     

## Table Structure  
ID defaults to timestamp (23 digits) + random number (9 digits), globally unique.  
Table creation statement ```gpressdatadir/gpress.sql```          

### Configuration (Table Name: config)
Reads ```gpressdatadir/install_config.json``` during installation.

| columnName  | Type        | Description         |  Remarks       | 
| ----------- | ----------- | ----------- | ----------- |
| id          | string      | Primary Key        |gpress_config |
| basePath    | string      | Base Path    |  Default /      |
| jwtSecret   | string      | JWT Secret     | Randomly generated     |
| jwttokenKey | string      | JWT Key    |  Default jwttoken  |
| serverPort  | string      | IP:Port     |  Default :660  |
| timeout     | int         | JWT Timeout Seconds|  Default 7200  |
| maxRequestBodySize | int  | Max Request Size     |  Default 20M  |
| locale      | string      | Language Pack       |  Default zh-CN,en-US |
| proxy       | string      | HTTP Proxy Address |             |
| createTime  | string      | Creation Time     |  2006-01-02 15:04:05  |
| updateTime  | string      | Update Time     |  2006-01-02 15:04:05  |
| createUser  | string      | Creator       |  Initialization system  |
| sortNo      | int         | Sort Order         |  Ascending  |
| status      | int         | Status     |  Link Access (0), Public (1), Top (2), Private (3)  |

### User (Table Name: user)
There is only one user in the backend.

| columnName  | Type         | Description        |  Remarks       | 
| ----------- | ----------- | ----------- | ----------- |
| id          | string      | Primary Key        | gpress_admin |
| account     | string      | Login Name    |  Default admin  |
| passWord    | string      | Password        |    -  |
| userName    | string      | Description        |    -  |
| createTime  | string      | Creation Time     |  2006-01-02 15:04:05  |
| updateTime  | string      | Update Time     |  2006-01-02 15:04:05  |
| createUser  | string      | Creator       |  Initialization system  |
| sortNo      | int         | Sort Order         |  Ascending  |
| status      | int         | Status     |  Link Access (0), Public (1), Top (2), Private (3)  |

### Site Information (Table Name:site)
Site information, such as title, logo, keywords, description, etc.

| columnName    | Type         | Description    |  Remarks       | 
| ----------- | ----------- | ----------- | ----------- |
| id          | string      | Primary Key        |gpress_site  |
| title       | string      | Site Name     |     -  |
| keyword     | string      | Keywords       |     -  |
| description | string      | Site Description    |     -  |
| theme       | string      | Default Theme     | Default uses default  |
| themePC     | string      | PC Theme      | First get from cookie, if not, get from Header, write to cookie, default uses default  |
| themeWAP    | string      | Mobile Theme    | First get from cookie, if not, get from Header, write to cookie, default uses default  |
| themeWX     | string      | WeChat Theme    | First get from cookie, if not, get from Header, write to cookie, default uses default  |
| logo        | string      | Logo       |     -  |
| favicon     | string      | Favicon    |     -  |
| createTime  | string      | Creation Time     |  2006-01-02 15:04:05  |
| updateTime  | string      | Update Time     |  2006-01-02 15:04:05  |
| createUser  | string      | Creator       |  Initialization system  |
| sortNo      | int         | Sort Order         |  Ascending  |
| status      | int         | Status     |  Link Access (0), Public (1), Top (2), Private (3)  |

### Navigation Menu (Table Name: category)
| columnName    | Type         | Description    |  Remarks       | 
| ----------- | ----------- | ----------- | ----------- |
| id          | string      | Primary Key         | URL path, separated by /, e.g., /web/ |
| name        | string      | Navigation Name     |    -  |
| hrefURL     | string      | Redirect Path     |    -  |
| hrefTarget  | string      | Redirect Method     | _self,_blank,_parent,_top|
| pid         | string      | Parent Navigation ID     | Parent Navigation ID  |
| templateFile  | string      | Template File       | Current navigation page template  |
| childTemplateFile  | string | Child Theme Template File  | Default template for child pages, if not set, default uses this template |
| keyword     | string      | Navigation Keywords   | Yes      |        |
| description | string      | Navigation Description     | Yes      |        |
| createTime  | string      | Creation Time     |  2006-01-02 15:04:05  |
| updateTime  | string      | Update Time     |  2006-01-02 15:04:05  |
| createUser  | string      | Creator       |  Initialization system  |
| sortNo      | int         | Sort Order         |  Ascending  |
| status      | int         | Status     |  Link Access (0), Public (1), Top (2), Private (3)  |

### Content (Table Name: content)
| columnName  | Type        | Description        | Whether to Tokenize |  Remarks                  | 
| ----------- | ----------- | ----------- | ------- | ---------------------- |
| id          | string      | Primary Key         |   No    | URL path, separated by /, e.g., /web/nginx-use-hsts |
| title       | string      | Title     | Yes      |     Uses jieba tokenizer    |
| keyword     | string      | Content Keywords   | Yes      |     Uses jieba tokenizer    |
| description | string      | Content Description     | Yes      |     Uses jieba tokenizer    |
| hrefURL     | string      | Self Page Path | No      |    -                    |
| subtitle    | string      | Subtitle       | Yes      |      Uses jieba tokenizer  |
| author      | string      | Author         | Yes      |      Uses jieba tokenizer  |
| tag         | string      | Tags         | Yes      |      Uses jieba tokenizer  |
| toc         | string      | Table of Contents         | Yes      |      Uses jieba tokenizer  |
| summary     | string      | Summary         | Yes      |      Uses jieba tokenizer  |
| categoryName| string      | Navigation Menu, separated by comma (,)| Yes| Uses jieba tokenizer.      |
| categoryID  | string      | Navigation ID       | No      | -                       |
| templateFile| string      | Template File     | No      | Template                    |
| content     | string      | Content     | No      |                         |
| markdown    | string      | Markdown Content | No      |                         |
| thumbnail   | string      | Cover Image       | No      |                         |
| signature   | string      | Private Key Signature of Content | No   |                         |
| signAddress | string      | Signature Address   | No   |                         |
| signChain   | string      | Chain of Address | No   |                         |
| txID        | string      | On-chain Transaction Hash  | No   |                         |
| createTime  | string      | Creation Time     | -       |  2006-01-02 15:04:05    |
| updateTime  | string      | Update Time     | -       |  2006-01-02 15:04:05    |
| createUser  | string      | Creator       | -       |  Initialization system          |
| sortNo      | int         | Sort Order         | -       |  Ascending                   |
| status      | int         | Status     | -       |  Link Access (0), Public (1), Top (2), Private (3)  |

#### Copyright and Software Copyright Description
* The software copyright registration number of this gpress is 2025SR0120223
* The software copyright of this gpress is owned by us, and secondary software copyright applications are prohibited. Infringement will be prosecuted
* The copyright of programs developed by developers using gpress belongs to the developers
* Please retain the copyright without any other restrictions. That is to say, you must include the original license agreement statement in your distribution, whether you distribute it in binary or source code form
* The open-source version is released under the AGPL-3.0 open-source license and is provided for free use, but it is not allowed to release and sell modified and derivative code as closed-source commercial software!