v1.2.1
- 升级go 1.26
- 修复ico格式的上传错误
- 完善文档,注释

v1.2.0
- 升级zorm v1.8.2
- 完善文档,注释

v1.1.9
- 全文检索markdown字段
- 支持PostgreSQL数据库,使用pg_search实现全文检索
- 数据库字段命名有驼峰修改为下划线
- 完善文档,注释

v1.1.8
- 更新fts分词字典
- 完善文档,注释

v1.1.7
- 优化日志输出
- 完善文档,注释

v1.1.6
- 增加Dockerfile
- 完善文档,注释

v1.1.5
- 修复刷新站点并发崩溃的Bug
- json无法序列化error类型,使用Message返回错误信息
- 修复日志输出Bug
- 完善文档,注释

v1.1.4
- 感谢@soldier_of_love,完善标签文档
- 修复sitemap生成bug
- 完善文档,注释

v1.1.3
- 修复静态资源压缩bug
- 完善文档,注释

v1.1.2
- 静态化时删除过期文件
- jodit 替换 wangEditor 
- 完善文档,注释

v1.1.1
- 感谢@soldier_of_love,拼音生成内容的路径标识
- 去掉路径里的空格
- 完善文档,注释

v1.1.0
- 增加自定义标签,支持解析视频/音频标签: !video[描述](视频地址) , !audio[描述](音频地址)
- 完善文档,注释

v1.0.9
- 升级FTS5分词组件
- 依赖Go 1.24
- 增强windows系统的兼容性
- 优化上传文件名
- 优化加载fts5分词扩展
- 完善文档,注释

v1.0.8
- 后台管理支持多语言
- 上传文件单独目录隔离,避免互相影响
- 后台管理页面增加 更新SQL 功能
- 内容表增加 txID 字段,记录上链交易的Hash;配置表增加 locale 字段,设置语言
- 修改错误的categorys拼写  
- 主题管理过滤掉.gz压缩文件
- 动态增加路由映射,去掉routeCategoryMap的处理逻辑
- 完善文档,注释

v1.0.7
- URI作为导航和内容的ID,例如:/web/是导航ID,/web/nginx-use-hsts是内容ID
- 去掉comCode,moduleID字段,增加signature,signAddress和signChain字段
- 统一映射静态文件,兼容项目前缀路径   
- 完善文档,注释

v1.0.6
- 使用本地的js资源文件
- Nginx 1.26+ 不需要再进行302重定向到目录下的index.html,gzip_static也会生效
- 根据Cookie和User-Agent请求头,为pc,wap,weixin配置不同的主题,并支持HarmonyOS
- site表中siteThemeWEIXIN字段名修改为themeWX
- 修复上传文件路径异常的问题
- 调整后台管理功能实现
- 完善文档,注释

v1.0.5
- 修复修改项目前缀造成访问异常的bug  
- 感谢 @soldier_of_love 的[pr](https://gitee.com/gpress/gpress/pulls/8),eth验签从go-ethereum库切换到dcrd库
- 感谢 @soldier_of_love 的[pr](https://gitee.com/gpress/gpress/pulls/7),优化静态路由配置，不用动态在更新路由  
- 感谢 @soldier_of_love 的[pr](https://gitee.com/gpress/gpress/pulls/6),解压失败的情况 合理删除目录  
- 感谢 @soldier_of_love 的[pr](https://gitee.com/gpress/gpress/pulls/5),解决sliceCategory2Tree三层及以上层数异常问题
- 感谢 @soldier_of_love 的[pr](https://gitee.com/gpress/gpress/pulls/4),解决AES-CBC模式加解密 部分文本异常
- 感谢 @soldier_of_love 的[pr](https://gitee.com/gpress/gpress/pulls/3),解决非法压缩包上传未删除问题
- 感谢 @soldier_of_love 的[pr](https://gitee.com/gpress/gpress/pulls/2),解决初始化DB副作用
- 完善文档,注释

v1.0.4
- 统一类型转换方法 convertType, 方便扩展
- 感谢 @lifj22 的[issue](https://gitee.com/gpress/gpress/issues/I9J1RH),导航模板category开头,内容模板content开头
- status字段增加置顶(2),原私密(2)修改为私密(3)
- 完善文档,注释

v1.0.3
- 栏目页自定义keyword,description,静态化生成sitemap.xml
- 增加seq标签,用于循环数字
- 迁移WordPress主题:[generatepress](https://gitee.com/gpress/wp-generatepress)、[astra](https://wpastra.com)
- 完善文档,注释

v1.0.2
- templateID 更名为 templateFile
- 修改示例数据
- 优化后台管理页面
- 完善文档,注释

v1.0.1
- 主题模板支持上传和文件修改
- 增加GitHub Actions文件,用于编译MacOS版本,感谢 @CDCDDCDC 编译libsimple.dylib
- 部署阿里云计算巢服务,方便云平台使用
- 完善文档,注释

v1.0.0
- 初始化版本