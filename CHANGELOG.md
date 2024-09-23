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