# gpress

## 介绍
云原生高性能内容平台,基于Hertz + Go template + Bleve全文检索实现,仅需100M运行内存   
默认端口660  
开发时需要先解压gpressdatadir/dict.zip      

## 贡献者授权说明
gpress使用GPL-3.0开源协议,特授权项目贡献者商业化时无需开源,说明如下:  
- 所有贡献者,根据代码提交记录(pr以合并为准),每一条提交记录,授权一个商业化站点
- 项目开发者角色,不限制商业化站点数量
- 项目管理员角色,不限制商业化站点数量,可二次授权客户商业化用途

## 软件架构
使用 Hertz + Go template + Bleve全文检索,不使用数据库    

不使用struct对象,全部使用map保存数据,可以随时增加属性字段.记录所有字段的名称,类型,中文名,code

模型的字段属性也是map,应用的文章回把模型的map属性全部取出,和自己的map覆盖合并.  

## 数据结构
所有的数据结构都使用Map实现,不再使用struct.因使用Bleve做NoSQL数据库,所以map可以任意添加字段.  
所有不需要分词的字符串,Mapping.Analyzer = keyword.Name 指定为keyword分词器.这样就可以类似数据库 name=value 作为精确的查询条件了.  
值使用逗号(,)隔开,使用 commaAnalyzerName 分词器,实现类似sql in 的效果.  

在TableInfo表里设置tableFiled='module',记录所有的Module.只是记录,并不创建index,全部保存到context里,用于全局检索    


ID默认使用时间戳(23位)+随机数(9位),全局唯一      
### 表信息(表名:tableInfo)

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| id          | string      | 主键         | 否      |    -  |
| code        | string      | 表Code     | 否      |    -  |
| name        | string      | 表名称     | 否      |    -  |
| tableFiled   | string      | 表类型     | 否      |  index:表.  module:模型  |
| createTime  |string  | 创建时间     | -       |  2006-01-02 15:04:05  |
| ipdateTime  |string  | 更新时间     | -       |  2006-01-02 15:04:05  |
| createUser  | string      | 创建人       | -       |  初始化 system  |
| sortNo      | int         | 排序         | -       |  正序  |
| status      | int         | 是否有效     | -       |  无效(0),正常显示(1),界面不显示(3)  |

### 表字段(表名:tableField)
记录所有表字段code和中文说明.  
理论上所有的表字段都可以放到这个表里,因为都是Map,就不需要再单独指定表的字段了,可以动态创建Table(目前建议这样做)  

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| ID          | string      | 主键        | 否       |    -  |
| TableCode   | string      | 表代码     | 否       |  类似表名 user,site,pageTemplate,navMenu,module,content  |
| TableName   | string      | 表名称     | 否       |  类似表名中文说明  |
| BusinessID  | string      | 业务ID       | 否       | 处理业务记录临时增加的字段,意外情况  |
| FieldCode   | string      | 字段代码     |否       |    -  |
| FieldName   | string      | 字段中文名称 | 否       |    -  |
| FieldType   | int         | 字段类型     | -       | 数字(1),日期(2),文本框(3),文本域(4),富文本(5),下拉框(6),单选(7),多选(8),上传图片(9),上传附件(10),轮播图(11),音频(12),视频(13)  |
| FieldFormat | string      | 数据格式,用于日期或者数字| 否 |  -  |    
| DefaultValue| string      | 默认值       | 否      |       -  |
| AnalyzerName| string      | 分词器名称    | -       | 为 '' 不设置  |
| CreateTime  |string  | 创建时间     | -       |  2006-01-02 15:04:05  |
| UpdateTime  |string  | 更新时间     | -       |  2006-01-02 15:04:05  |
| CreateUser  | string      | 创建人       | -       |  初始化 system  |
| SortNo      | int         | 排序         | -       |  正序  |
| Status      | int         | 是否有效     | -       |  无效(0),正常显示(1),界面不显示(3)  |

### 用户(表名:user)
后台只有一个用户,账号admin 密码默认admin 可以自己修改.

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| id          | string      | 主键        | 否       |    -  |
| account     | string      | 登录名称     | 否       |  默认admin  |
| passWord    | string      | 密码        | 否       |    -  |
| userName    | string      | 中文名称     | 否       |    -  |

### 站点信息(site)
站点的信息,例如 title,logo,keywords,description等

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| id          | string      | 主键         | 否      |    -  |
| title       | string      | 站点名称     | 否      |     -  |
| keyword     | string      | 关键字       | 否      |     -  |
| description | string      | 站点描述     | 否      |     -  |
| theme       | string      | 默认主题        | 否      | 默认使用default  |
| themePC     | string      | PC主题      | 否      | 先从cookie获取,如果没有从Header头取值,写入cookie,默认使用default  |
| themeWAP    | string      | 手机主题     | 否      | 先从cookie获取,如果没有从Header头取值,写入cookie,默认使用default  |
| themeWEIXIN | string      | 微信主题     | 否      | 先从cookie获取,如果没有从Header头取值,写入cookie,默认使用default  |
| logo        | string      | logo        | 否      |     -  |
| favicon     | string      | Favicon     | 否      |     -  |


### 页面模板(表名:pageTemplate)
后台只有一个用户,账号admin 密码默认admin 可以自己修改.

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| id          | string      | 主键        | 否       |    -  |
| templateName| string      | 模板名称     | 否       |    -  |
| templatePath| string      | 模板路径     | 否       |    -  |
| sortNo      | int         | 排序        | -       |  正序  |
| status      | int         | 是否有效     | -       |  无效(0),正常显示(1),界面不显示(3)  |

### 导航菜单(表名:navMenu)
| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| id          | string      | 主键        | 否       |    -  |
| menuName    | string      | 菜单名称     | 否       |    -  |
| menuName    | string      | 菜单名称     | 否       |    -  |
| hrefURL     | string      | 跳转路径     | 否       |    -  |
| hrefTarget  | string      | 跳转方式     | 否       | _self,_blank,_parent,_top|
| pid         | string      | 父菜单ID     | 否       | 父菜单ID  |
| moduleID    | string      | module表ID | 否       |  导航菜单下的文章默认使用的模型字段 |
| comCode     | string      | 逗号隔开的全路径 | 否    | 逗号隔开的全路径  |
| templateID  | string      | 模板Id       | 否       | 当前导航页的模板  |
| childTemplateID  | string | 子页面模板Id  | 否      | 子页面默认使用的模板,子页面如果不设置,默认使用这个模板 |
| sortNo      | int         | 排序        | -       |  正序  |
| status      | int         | 是否有效     | -       |  无效(0),正常显示(1),界面不显示(3)  |

### 模型数据(表名:module_default)
在tableInfo表里设置tableFiled='module',记录所有的Module.只是记录,并不创建index,全部保存到context里,用于全局检索   

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| id          | string      | 主键         | 否      |    -  |
| title       | string      | 文章标题     | 是      |     -  |
| keyword     | string      | 关键字       | 是      |    使用 commaAnalyzerName 分词器,实现类似sql in 的效果.    |
| description | string      | 站点描述     | 否      |     -  |
| pageURL     | string      | 自身页面路径 | 否       |    -  |
| subtitle    | string      | 副标题       | 是      |     -  |
| content     | string      | 文章内容     | 是      |       |
| createTime  |string  | 创建时间     | -       |  2006-01-02 15:04:05  |
| updateTime  |string  | 更新时间     | -       |  2006-01-02 15:04:05  |
| createUser  | string      | 创建人       | -       |  初始化 system  |
| sortNo      | int         | 排序        | -       |  正序  |
| status      | int         | 是否有效     | -       |  无效(0),正常显示(1),界面不显示(3)  |


### 文章内容(表名:content)
文章内容表,默认使用 module_default 的模型字段

| codeName    | 类型         | 中文名称    | 是否分词 |  备注       | 
| ----------- | ----------- | ----------- | ------- | ----------- |
| id          | string      | 主键         | 否      |    -  |
| moduleID    | string      | 模型ID       | 否      |  文章使用的模型字段 |
| title       | string      | 文章标题     | 是      |     -  |
| keyword     | string      | 关键字       | 是      |    使用 commaAnalyzerName 分词器,实现类似sql in 的效果.    |
| description | string      | 站点描述     | 否      |     -  |
| pageURL     | string      | 自身页面路径 | 否       |    -  |
| subtitle    | string      | 副标题       | 是      |     -  |
| navMenuID   | string      | 导航ID       | 否      | -    |
| navMenuNames| string      | 导航名称,逗号(,)隔开     | 是      | 使用 gseAnalyzerMapping 分词器.  |
| templateID  | string      | 模板Id       | 否      | 模板  |
| content     | string      | 文章内容     | 是      |       |
| createTime  |string  | 创建时间     | -       |  2006-01-02 15:04:05  |
| updateTime  |string  | 更新时间     | -       |  2006-01-02 15:04:05  |
| createUser  | string      | 创建人       | -       |  初始化 system  |
| sortNo      | int         | 排序        | -       |  正序  |
| status      | int         | 是否有效     | -       |  无效(0),正常显示(1),界面不显示(3)  |

