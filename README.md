# zcms

## 介绍
云原生高性能CMS,基于gin + golang template + Bleve全文检索实现,128M内存即可完美运行  

功能说明:  
- 后台就一种用户,登录就是管理员,菜单路由写死,不需要权限分配  
- 记录登录失败的次数,密码连续错误5次屏蔽账号10分钟,防止暴力登录  
- 全文检索产生的文件,用户上传的文件,站点的模板文件放到外部目录
- 后台提供网站导航菜单管理,广告位管理,底部footer条管理,新闻资讯管理和表单提交管理. 模板暂不提供在线编辑功能,开发人员编辑好上传到指定目录.
- 因为使用nosql(全文检索),所以文章的属性字段可以随意添加,code用于检索,中文name用于显示.


## 软件架构
使用 gin + golang template + Bleve全文检索,不使用数据库  
使用golang 1.16 的新特性 Go embed,打包静态资源文件  

不使用struct对象,全部使用map保存数据,可以随时增加属性字段.记录所有字段的名称,类型,中文名,code 用于界面展示  

模型的字段属性也是map,应用的文章回把模型的map属性全部取出,和自己的map覆盖合并.  


## 数据结构
所有的数据结构都使用Map实现,不再使用struct.因使用Bleve做NoSQL数据库,所以map可以任意添加字段.  
所有不需要分词的字符串,Mapping.Analyzer = keyword.Name 指定为keyword分词器.这样就可以类似数据库 name=value 作为精确的查询条件了.  
还需要实现 分号(,) 的分词器,实现类似sql in 的效果.  
在IndexField表里设置IndexCode='Module',然后在IndexField中插入每个module的字段,每个module实例的ModuleCode都是不同的,使用Module_+后缀的方式命名,只是记录,并不创建index  


### 索引和字段(索引名:IndexField)
记录所有索引字段code和中文说明,用于前台界面渲染展示.  
理论上所有的索引字段都可以放到这个表里,因为都是Map,就不需要再单独指定索引的字段了,统一从这里查询出来.(目前建议这样做,因为全是Map)

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| ID          | string      | 主键        | 否       |    -  |
| IndexCode   | string      | 索引代码     | 否       |  类似表名 User,SiteInfo,PageTemplate,NavMenu,Module,Content  |
| IndexName   | string      | 索引名称     | 否       |  类似表名中文说明  |
| BusinessID  | string      | 业务ID       | 否       | 处理业务记录临时增加的字段,意外情况  |
| FieldCode   | string      | 字段代码     |否       |    -  |
| FieldName   | string      | 字段中文名称 | 否       |    -  |
| FieldType   | int         | 字段类型     | -       | 数字(1),日期(2),文本(3),下拉框(4),单选(5),多选(6),上传图片(7),上传附件(8),轮播图(9)  |
| AnalyzerName| string      | 分词器名称    | -       | 为 '' 不分词  |
| CreateTime  | time.Time   | 创建时间     | -       |  2006-01-02 15:04:05  |
| UpdateTime  | time.Time   | 更新时间     | -       |  2006-01-02 15:04:05  |
| CreateUser  | string      | 创建人       | -       |  默认 admin  |
| SortNo      | int         | 排序         | -       |  正序  |
| Active      | int         | 是否有效     | -       |  无效(0),有效(1)  |

### 用户(索引名:User)
后台只有一个用户,账号admin 密码默认admin 可以自己修改.

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| ID          | string      | 主键        | 否       |    -  |
| Account     | string      | 登录名称     | 否       |  默认admin  |
| PassWord    | string      | 密码        | 否       |    -  |
| UserName    | string      | 中文名称     | 否       |    -  |

### 站点信息(SiteInfo)
站点的信息,例如 title,logo,keywords,description等

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| ID          | string      | 主键         | 否      |    -  |
| Title       | string      | 站点名称     | 否      |     -  |
| KeyWords    | string      | 关键字       | 否      |     -  |
| Description | string      | 站点描述     | 否      |     -  |
| Logo        | string      | logo        | 否      |     -  |
| Favicon     | string      | Favicon     | 否      |     -  |

### 页面模板(索引名:PageTemplate)
后台只有一个用户,账号admin 密码默认admin 可以自己修改.

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| ID          | string      | 主键        | 否       |    -  |
| TemplateName| string      | 模板名称     | 否       |    -  |
| TemplatePath| string      | 模板路径     | 否       |    -  |
| SortNo      | int         | 排序        | -       |  正序  |
| Active      | int         | 是否有效     | -       |  无效(0),有效(1)  |

### 导航菜单(索引名:NavMenu)
| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| ID          | string      | 主键        | 否       |    -  |
| MenuName    | string      | 菜单名称     | 否       |    -  |
| HrefURL     | string      | 跳转路径     | 否       |    -  |
| HrefTarget  | string      | 跳转方式     | 否       | _self,_blank,_parent,_top|
| PID         | string      | 父菜单ID     | 否       | 父菜单ID  |
| ModuleIndexCode| string | Module的索引名称 | 否    |  导航菜单下的文章默认使用的模型字段 |
| ComCode     | string      | 逗号隔开的全路径 | 否    | 逗号隔开的全路径  |
| TemplateID  | string      | 模板Id       | 否       | 当前导航页的模板  |
| ChildTemplateID  | string | 子页面模板Id  | 否      | 子页面默认使用的模板,子页面如果不设置,默认使用这个模板 |
| SortNo      | int         | 排序        | -       |  正序  |
| Active      | int         | 是否有效     | -       |  无效(0),有效(1)  |


### 模型数据(索引名:Module_demo)
~~文章模型,只是用来声明字段,具体信息会有Content索引全部继承~~  
暂时不使用了,这里只做参考.  
在IndexField表里设置IndexCode='Module',然后在IndexField中插入每个module的字段,每个module实例的ModuleCode都是不同的,使用Module_+后缀的方式命名,只是记录,并不创建index

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| ID          | string      | 主键         | 否      |    -  |
| ModuleIndexCode | string  | 模型Code     | 否      |    -  |
| ModuleName  | string      | 模型名称     | 否      |    -  |
| CreateTime  | time.Time   | 创建时间     | -       |  2006-01-02 15:04:05  |
| UpdateTime  | time.Time   | 更新时间     | -       |  2006-01-02 15:04:05  |
| CreateUser  | string      | 创建人       | -       |  默认 admin  |
| SortNo      | int         | 排序        | -       |  正序  |
| Active      | int         | 是否有效     | -       |  无效(0),有效(1)  |

### 模型数据(索引名:Module_demo)
在IndexField表里设置IndexCode='Module',然后在IndexField中插入每个module的字段,每个module实例的ModuleCode都是不同的,使用Module_+后缀的方式命名,只是记录,并不创建index

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| ID          | string      | 主键         | 否      |    -  |
| Title       | string      | 文章标题     | 是      |     -  |
| KeyWords    | string      | 关键字       | 否      |     -  |
| Description | string      | 站点描述     | 否      |     -  |
| Subtitle    | string      | 副标题       | 是      |     -  |
| Content     | string      | 文章内容     | 是      |       |
| CreateTime  | time.Time   | 创建时间     | -       |  2006-01-02 15:04:05  |
| UpdateTime  | time.Time   | 更新时间     | -       |  2006-01-02 15:04:05  |
| CreateUser  | string      | 创建人       | -       |  默认 admin  |
| SortNo      | int         | 排序        | -       |  正序  |
| Active      | int         | 是否有效     | -       |  无效(0),有效(1)  |


### 文章内容(索引名:Content)
文章内容表

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| ID          | string      | 主键         | 否      |    -  |
| HrefURL     | string      | 页面路径     | 否       |    -  |
| ModuleIndexCode| string   | 模型的Code   | 否   |  文章默认使用的模型字段 |
| NavMenuId   | string      | 导航ID       | 否      | 最好处理一下 分号 分词 类似in的效果  |
| NavMenuName | string      | 导航名称     | 是      | 最好处理一下 分号 分词 类似in的效果  |
| TemplateID  | string      | 模板Id       | 否      | 模板  |
| Content     | string      | 文章内容     | 是      |       |
| CreateTime  | time.Time   | 创建时间     | -       |  2006-01-02 15:04:05  |
| UpdateTime  | time.Time   | 更新时间     | -       |  2006-01-02 15:04:05  |
| CreateUser  | string      | 创建人       | -       |  默认 admin  |
| SortNo      | int         | 排序        | -       |  正序  |
| Active      | int         | 是否有效     | -       |  无效(0),有效(1)  |

