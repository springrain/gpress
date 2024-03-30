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
		templateID        TEXT,
		childTemplateID        TEXT,
		createTime        TEXT,
		updateTime        TEXT,
		createUser        TEXT,
		sortNo            int NOT NULL,
		status            int NOT NULL
	 ) strict ;
INSERT INTO category (status,sortNo,createUser,updateTime,createTime,childTemplateID,templateID,comCode,moduleID,themePC,pid,hrefTarget,hrefURL,name,id) VALUES (1,2,NULL,'2023-06-27 22:41:20','2023-06-27 22:41:20',NULL,NULL,',web,',NULL,NULL,NULL,'',NULL,'Web','web');
INSERT INTO category (status,sortNo,createUser,updateTime,createTime,childTemplateID,templateID,comCode,moduleID,themePC,pid,hrefTarget,hrefURL,name,id) VALUES (1,1,NULL,'2023-06-27 22:41:20','2023-06-27 22:41:20',NULL,NULL,',about,',NULL,NULL,NULL,'','/post/about','About','about');

CREATE TABLE IF NOT EXISTS pageTemplate (
		id TEXT PRIMARY KEY     NOT NULL,
		name         TEXT  NOT NULL,
		templatePath         TEXT   NOT NULL,
		createTime        TEXT,
		updateTime        TEXT,
		createUser        TEXT,
		sortNo            int NOT NULL,
		status            int NOT NULL
	 ) strict ;

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
		templateID           TEXT,
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
INSERT INTO site (status,sortNo,createUser,updateTime,createTime,footer,favicon,logo,siteThemeWEIXIN,themeWAP,themePC,theme,description,keyword,domain,name,title,id)VALUES (1,1,NULL,NULL,NULL,'<div class="copyright"><span class="copyright-year">&copy; 2008 - 2024<span class="author">jiagou.com 版权所有 <a href=''https://beian.miit.gov.cn'' target=''_blank''>豫ICP备xxxxx号</a>   <a href=''http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=xxxx''  target=''_blank''><img src=''/public/gongan.png''>豫公网安备xxxxx号</a></span></span></div>','/public/favicon.png','/public/logo.png','default','default','default','default','gpress','gpress','jiagou.com','架构','jiagou','gpress_site');

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
INSERT INTO content (comCode,status,sortNo,createUser,updateTime,createTime,thumbnail,markdown,content,summary,toc,tag,author,templateID,categoryName,categoryID,subtitle,hrefURL,description,keyword,title,moduleID,id) 
VALUES (',about,',0,100,NULL,'2023-06-27 22:43:53','2023-06-27 22:43:53',NULL,
'---
title: "About"
date: 2017-08-20T21:38:52+08:00
lastmod: 2017-08-28T21:41:52+08:00
menu: "about"
weight: 50

---

本站服务器配置:2核CPU,1G内存,20G硬盘,AnolisOS(ANCK).  
使用gpress迁移了hugo的even主题和markdown文件.  




我所见识过的一切都将消失一空,就如眼泪消逝在雨中......  
不妨大胆一些,大胆一些......  





小项目:  
* [springrain](https://gitee.com/chunanyong/springrain)
* [zorm](https://gitee.com/chunanyong/zorm)
* [gpress](https://gitee.com/gpress/gpress)
* [gowe](https://gitee.com/chunanyong/gowe)','<p>本站服务器配置:2核CPU,1G内存,20G硬盘,AnolisOS(ANCK).<br>
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
','本站服务器配置:1核CPU,512M内存,20G硬盘,AnolisOS(ANCK).使用hugo和even模板,编译成静态文件,Nginx作为WEB服务器.我所见识过的一切都将消失一空,就如眼泪消逝在雨中......
	不妨大胆一些,大胆一些......','',NULL,'springrain',NULL,'about','about',NULL,NULL,NULL,NULL,'about',NULL,'about');
				
INSERT INTO content (comCode,status,sortNo,createUser,updateTime,createTime,thumbnail,markdown,content,summary,toc,tag,author,templateID,categoryName,categoryID,subtitle,hrefURL,description,keyword,title,moduleID,id)
                    VALUES (',web,',1,58,NULL,'2020-09-02 14:05:02','2020-09-02 14:05:02',NULL,
                    '---
title: "springrain项目说明"
date: 2020-09-02T14:05:02+08:00
draft: false
tags: ["springrain"]
categories: ["web"]
author: "springrain"
contentCopyright: ''<a rel="license noopener" href="https://en.wikipedia.org/wiki/Wikipedia:Text_of_Creative_Commons_Attribution-ShareAlike_3.0_Unported_License" target="_blank">Creative Commons Attribution-ShareAlike License</a>''
---

## 微服务
 6.0.0 项目入口是springrain-system-web,基于Istio实现微服务,正在整理文档.
## 前后分离
 6.0.0 基于VUE前后端分离,使用JWT认证.前端项目是[springrain-vue](https://gitee.com/chunanyong/springrain-vue)

## 实现了什么?
* 不增加学习成本,像单体一样开发分布式微服务.
* 不修改业务代码,可以实现单体,分层,微服务多种部署模式切换.
* 集成seata分布式事务实现. 

## 实现思路
* 启动加载springbean时,先检查本地是否有实现,如果没有就启动GRPC远程调用.开发人员无感知.
* 基于seata分布式事务实现.支持有注解和无注解(开发人员无感知,理论上有不同步风险,个人感觉做好日志,风险不大)混合使用.
* 基于Istio实现微服务的发现,监控,熔断,限流.开发人员无感知.

## 限制
* 接口和实现的命名强制规范.
* 一个RPC接口只能有一个实现.
* 分布式事务,一定要避免A服务update表t,RPC调用B服务,B服务也update表t.这样A等待B结果,B等待A释放锁,造成死锁.
* 分布式无注解比较方便,理论上有不同步风险,个人感觉做好日志,风险不大
* Service层不可以使用Servlet API,例如 HttpRequest

## 体验单体到分层切换
* 修改springrain-system-web依赖springrain-system-service,不再依赖springrain-system-serviceimpl.
* springrain-system-serviceimpl添加springrain-grpc-server依赖,启用org.springrain.SystemServiceImplApplication的@SpringBootApplication注解
* seata-server的conf目录下file.conf,修改vgroup_mapping.my_test_tx_group = "default" 为 vgroup_mapping.seata_tx_group = "default",启动seata-server服务.
* 启动springrain-system-serviceimpl
* 启动springrain-system-web
* 访问http://127.0.0.1:8080/api/checkHealth

## 博客
项目名为springrain[春雨]我的个人博客是 http://www.jiagou.com </br>
## 文档
https://gitee.com/chunanyong/springrain/tree/master/springrain-system/springrain-system-web/doc  </br>
## 代码生成器
https://gitee.com/chunanyong/springrain/tree/master/springrain-gencode  </br>
## sql脚本
https://gitee.com/chunanyong/springrain/tree/master/springrain-system/springrain-system-web/sql  </br>


springrain是spring/springboot的封装,springboot开发范例.

springrain是一个Maven项目,包含spring core,spring jdbc,spring mvc.

springrain自带代码生成器,能够生成对表的增删改查的逻辑代码,以及前台页面样式和js文件

项目只依赖spring,没有hibernate,struts,ibatis.

使用jwt认证.

数据库调优可以使用druid输出慢sql,比分析xml中的语句更直观,springrain所有的sql语句都使用Finder封装管理,只要查看Finder的引用即可.

## 案例

```java
//就极简而言,一个数据库只需要一个Service,就可以管理这个数据库的任意一张表 
//@Test  查询基本类型
public void testObject() throws Exception{
       // Finder finder=new Finder("select id from t_user where 1=1 ");
        Finder finder=Finder.getSelectFinder(User.class,"id").append(" WHERE 1=1 "); 
         finder.append("and id=:userId ").setParam("userId", "admin");
        String id = baseDemoService.queryForObject(finder, String.class);
        System.out.println(id);

}

//@Test 查询一个对象
public void testObjectUser() throws Exception{
        //Finder finder=new Finder("select * from t_user where id=:userId order by id"); 
Finder finder=Finder.getSelectFinder(User.class).append(" WHERE id=:userId order by id desc "); 
        finder.setParam("userId", "admin");
        User u = baseDemoService.queryForObject(finder, User.class);
        System.out.println(u.getName());

}
//@Test 查询分页
public void testMsSql() throws Exception{
        //Finder finder=new Finder("select * from t_user order by id");
        Finder finder=Finder.getSelectFinder(User.class).append(" order by id desc ");
        Listlist = baseDemoService.queryForList(finder, User.class, new Page(2));
        System.out.println(list.size());
        for(User s:list){
         System.out.println(s.getName());
         }
}



//@Test 调用数据库存储过程
public void testProc() throws Exception{
        Finder finder=new Finder();
        finder.setParam("unitId", 0);
        finder.setProcName("proc_up");
        Map queryObjectByProc = (Map) baseDemoService.queryObjectByProc(finder);
        System.out.println(queryObjectByProc.get("#update-count-10"));
        

}

//@Test 调用数据库函数
public void testFunction() throws Exception{
        Finder finder=new Finder();
        finder.setFunName("fun_userId");
        finder.setParam("userId", "admin");
        String userName= baseDemoService.queryForObjectByByFunction(finder,String.class);
        System.out.println(userName);
}

```','<h2 id="微服务">微服务</h2>
<p>6.0.0 项目入口是springrain-system-web,基于Istio实现微服务,正在整理文档.</p>
<h2 id="前后分离">前后分离</h2>
<p>6.0.0 基于VUE前后端分离,使用JWT认证.前端项目是<a href="https://gitee.com/chunanyong/springrain-vue">springrain-vue</a></p>
<h2 id="实现了什么?">实现了什么?</h2>
<ul>
<li>不增加学习成本,像单体一样开发分布式微服务.</li>
<li>不修改业务代码,可以实现单体,分层,微服务多种部署模式切换.</li>
<li>集成seata分布式事务实现.</li>
</ul>
<h2 id="实现思路">实现思路</h2>
<ul>
<li>启动加载springbean时,先检查本地是否有实现,如果没有就启动GRPC远程调用.开发人员无感知.</li>
<li>基于seata分布式事务实现.支持有注解和无注解(开发人员无感知,理论上有不同步风险,个人感觉做好日志,风险不大)混合使用.</li>
<li>基于Istio实现微服务的发现,监控,熔断,限流.开发人员无感知.</li>
</ul>
<h2 id="限制">限制</h2>
<ul>
<li>接口和实现的命名强制规范.</li>
<li>一个RPC接口只能有一个实现.</li>
<li>分布式事务,一定要避免A服务update表t,RPC调用B服务,B服务也update表t.这样A等待B结果,B等待A释放锁,造成死锁.</li>
<li>分布式无注解比较方便,理论上有不同步风险,个人感觉做好日志,风险不大</li>
<li>Service层不可以使用Servlet API,例如 HttpRequest</li>
</ul>
<h2 id="体验单体到分层切换">体验单体到分层切换</h2>
<ul>
<li>修改springrain-system-web依赖springrain-system-service,不再依赖springrain-system-serviceimpl.</li>
<li>springrain-system-serviceimpl添加springrain-grpc-server依赖,启用org.springrain.SystemServiceImplApplication的@SpringBootApplication注解</li>
<li>seata-server的conf目录下file.conf,修改vgroup_mapping.my_test_tx_group = &quot;default&quot; 为 vgroup_mapping.seata_tx_group = &quot;default&quot;,启动seata-server服务.</li>
<li>启动springrain-system-serviceimpl</li>
<li>启动springrain-system-web</li>
<li>访问http://127.0.0.1:8080/api/checkHealth</li>
</ul>
<h2 id="博客">博客</h2>
<p>项目名为springrain[春雨]我的个人博客是 <a href="http://www.jiagou.com">http://www.jiagou.com</a> <!-- raw HTML omitted --></p>
<h2 id="文档">文档</h2>
<p><a href="https://gitee.com/chunanyong/springrain/tree/master/springrain-system/springrain-system-web/doc">https://gitee.com/chunanyong/springrain/tree/master/springrain-system/springrain-system-web/doc</a>  <!-- raw HTML omitted --></p>
<h2 id="代码生成器">代码生成器</h2>
<p><a href="https://gitee.com/chunanyong/springrain/tree/master/springrain-gencode">https://gitee.com/chunanyong/springrain/tree/master/springrain-gencode</a>  <!-- raw HTML omitted --></p>
<h2 id="sql脚本">sql脚本</h2>
<p><a href="https://gitee.com/chunanyong/springrain/tree/master/springrain-system/springrain-system-web/sql">https://gitee.com/chunanyong/springrain/tree/master/springrain-system/springrain-system-web/sql</a>  <!-- raw HTML omitted --></p>
<p>springrain是spring/springboot的封装,springboot开发范例.</p>
<p>springrain是一个Maven项目,包含spring core,spring jdbc,spring mvc.</p>
<p>springrain自带代码生成器,能够生成对表的增删改查的逻辑代码,以及前台页面样式和js文件</p>
<p>项目只依赖spring,没有hibernate,struts,ibatis.</p>
<p>使用jwt认证.</p>
<p>数据库调优可以使用druid输出慢sql,比分析xml中的语句更直观,springrain所有的sql语句都使用Finder封装管理,只要查看Finder的引用即可.</p>
<h2 id="案例">案例</h2>
<div class="highlight"><div class="chroma">
<table class="lntable"><tr><td class="lntd">
<pre tabindex="0" class="chroma"><code class="language-fallback" data-lang="fallback"><span class="lnt"> 1
</span><span class="lnt"> 2
</span><span class="lnt"> 3
</span><span class="lnt"> 4
</span><span class="lnt"> 5
</span><span class="lnt"> 6
</span><span class="lnt"> 7
</span><span class="lnt"> 8
</span><span class="lnt"> 9
</span><span class="lnt">10
</span><span class="lnt">11
</span><span class="lnt">12
</span><span class="lnt">13
</span><span class="lnt">14
</span><span class="lnt">15
</span><span class="lnt">16
</span><span class="lnt">17
</span><span class="lnt">18
</span><span class="lnt">19
</span><span class="lnt">20
</span><span class="lnt">21
</span><span class="lnt">22
</span><span class="lnt">23
</span><span class="lnt">24
</span><span class="lnt">25
</span><span class="lnt">26
</span><span class="lnt">27
</span><span class="lnt">28
</span><span class="lnt">29
</span><span class="lnt">30
</span><span class="lnt">31
</span><span class="lnt">32
</span><span class="lnt">33
</span><span class="lnt">34
</span><span class="lnt">35
</span><span class="lnt">36
</span><span class="lnt">37
</span><span class="lnt">38
</span><span class="lnt">39
</span><span class="lnt">40
</span><span class="lnt">41
</span><span class="lnt">42
</span><span class="lnt">43
</span><span class="lnt">44
</span><span class="lnt">45
</span><span class="lnt">46
</span><span class="lnt">47
</span><span class="lnt">48
</span><span class="lnt">49
</span><span class="lnt">50
</span><span class="lnt">51
</span><span class="lnt">52
</span><span class="lnt">53
</span></code></pre></td>
<td class="lntd">
<pre tabindex="0" class="chroma"><code class="language-java" data-lang="java"><span class="line"><span class="cl"><span class="c1">//就极简而言,一个数据库只需要一个Service,就可以管理这个数据库的任意一张表 
</span></span></span><span class="line"><span class="cl"><span class="c1">//@Test  查询基本类型
</span></span></span><span class="line"><span class="cl"><span class="c1"></span><span class="kd">public</span> <span class="kt">void</span> <span class="nf">testObject</span><span class="o">()</span> <span class="kd">throws</span> <span class="n">Exception</span><span class="o">{</span>
</span></span><span class="line"><span class="cl">       <span class="c1">// Finder finder=new Finder(&#34;select id from t_user where 1=1 &#34;);
</span></span></span><span class="line"><span class="cl"><span class="c1"></span>        <span class="n">Finder</span> <span class="n">finder</span><span class="o">=</span><span class="n">Finder</span><span class="o">.</span><span class="na">getSelectFinder</span><span class="o">(</span><span class="n">User</span><span class="o">.</span><span class="na">class</span><span class="o">,</span><span class="s">&#34;id&#34;</span><span class="o">).</span><span class="na">append</span><span class="o">(</span><span class="s">&#34; WHERE 1=1 &#34;</span><span class="o">);</span> 
</span></span><span class="line"><span class="cl">         <span class="n">finder</span><span class="o">.</span><span class="na">append</span><span class="o">(</span><span class="s">&#34;and id=:userId &#34;</span><span class="o">).</span><span class="na">setParam</span><span class="o">(</span><span class="s">&#34;userId&#34;</span><span class="o">,</span> <span class="s">&#34;admin&#34;</span><span class="o">);</span>
</span></span><span class="line"><span class="cl">        <span class="n">String</span> <span class="n">id</span> <span class="o">=</span> <span class="n">baseDemoService</span><span class="o">.</span><span class="na">queryForObject</span><span class="o">(</span><span class="n">finder</span><span class="o">,</span> <span class="n">String</span><span class="o">.</span><span class="na">class</span><span class="o">);</span>
</span></span><span class="line"><span class="cl">        <span class="n">System</span><span class="o">.</span><span class="na">out</span><span class="o">.</span><span class="na">println</span><span class="o">(</span><span class="n">id</span><span class="o">);</span>
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl"><span class="o">}</span>
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl"><span class="c1">//@Test 查询一个对象
</span></span></span><span class="line"><span class="cl"><span class="c1"></span><span class="kd">public</span> <span class="kt">void</span> <span class="nf">testObjectUser</span><span class="o">()</span> <span class="kd">throws</span> <span class="n">Exception</span><span class="o">{</span>
</span></span><span class="line"><span class="cl">        <span class="c1">//Finder finder=new Finder(&#34;select * from t_user where id=:userId order by id&#34;); 
</span></span></span><span class="line"><span class="cl"><span class="c1"></span><span class="n">Finder</span> <span class="n">finder</span><span class="o">=</span><span class="n">Finder</span><span class="o">.</span><span class="na">getSelectFinder</span><span class="o">(</span><span class="n">User</span><span class="o">.</span><span class="na">class</span><span class="o">).</span><span class="na">append</span><span class="o">(</span><span class="s">&#34; WHERE id=:userId order by id desc &#34;</span><span class="o">);</span> 
</span></span><span class="line"><span class="cl">        <span class="n">finder</span><span class="o">.</span><span class="na">setParam</span><span class="o">(</span><span class="s">&#34;userId&#34;</span><span class="o">,</span> <span class="s">&#34;admin&#34;</span><span class="o">);</span>
</span></span><span class="line"><span class="cl">        <span class="n">User</span> <span class="n">u</span> <span class="o">=</span> <span class="n">baseDemoService</span><span class="o">.</span><span class="na">queryForObject</span><span class="o">(</span><span class="n">finder</span><span class="o">,</span> <span class="n">User</span><span class="o">.</span><span class="na">class</span><span class="o">);</span>
</span></span><span class="line"><span class="cl">        <span class="n">System</span><span class="o">.</span><span class="na">out</span><span class="o">.</span><span class="na">println</span><span class="o">(</span><span class="n">u</span><span class="o">.</span><span class="na">getName</span><span class="o">());</span>
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl"><span class="o">}</span>
</span></span><span class="line"><span class="cl"><span class="c1">//@Test 查询分页
</span></span></span><span class="line"><span class="cl"><span class="c1"></span><span class="kd">public</span> <span class="kt">void</span> <span class="nf">testMsSql</span><span class="o">()</span> <span class="kd">throws</span> <span class="n">Exception</span><span class="o">{</span>
</span></span><span class="line"><span class="cl">        <span class="c1">//Finder finder=new Finder(&#34;select * from t_user order by id&#34;);
</span></span></span><span class="line"><span class="cl"><span class="c1"></span>        <span class="n">Finder</span> <span class="n">finder</span><span class="o">=</span><span class="n">Finder</span><span class="o">.</span><span class="na">getSelectFinder</span><span class="o">(</span><span class="n">User</span><span class="o">.</span><span class="na">class</span><span class="o">).</span><span class="na">append</span><span class="o">(</span><span class="s">&#34; order by id desc &#34;</span><span class="o">);</span>
</span></span><span class="line"><span class="cl">        <span class="n">Listlist</span> <span class="o">=</span> <span class="n">baseDemoService</span><span class="o">.</span><span class="na">queryForList</span><span class="o">(</span><span class="n">finder</span><span class="o">,</span> <span class="n">User</span><span class="o">.</span><span class="na">class</span><span class="o">,</span> <span class="k">new</span> <span class="n">Page</span><span class="o">(</span><span class="mi">2</span><span class="o">));</span>
</span></span><span class="line"><span class="cl">        <span class="n">System</span><span class="o">.</span><span class="na">out</span><span class="o">.</span><span class="na">println</span><span class="o">(</span><span class="n">list</span><span class="o">.</span><span class="na">size</span><span class="o">());</span>
</span></span><span class="line"><span class="cl">        <span class="k">for</span><span class="o">(</span><span class="n">User</span> <span class="n">s</span><span class="o">:</span><span class="n">list</span><span class="o">){</span>
</span></span><span class="line"><span class="cl">         <span class="n">System</span><span class="o">.</span><span class="na">out</span><span class="o">.</span><span class="na">println</span><span class="o">(</span><span class="n">s</span><span class="o">.</span><span class="na">getName</span><span class="o">());</span>
</span></span><span class="line"><span class="cl">         <span class="o">}</span>
</span></span><span class="line"><span class="cl"><span class="o">}</span>
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl"><span class="c1">//@Test 调用数据库存储过程
</span></span></span><span class="line"><span class="cl"><span class="c1"></span><span class="kd">public</span> <span class="kt">void</span> <span class="nf">testProc</span><span class="o">()</span> <span class="kd">throws</span> <span class="n">Exception</span><span class="o">{</span>
</span></span><span class="line"><span class="cl">        <span class="n">Finder</span> <span class="n">finder</span><span class="o">=</span><span class="k">new</span> <span class="n">Finder</span><span class="o">();</span>
</span></span><span class="line"><span class="cl">        <span class="n">finder</span><span class="o">.</span><span class="na">setParam</span><span class="o">(</span><span class="s">&#34;unitId&#34;</span><span class="o">,</span> <span class="mi">0</span><span class="o">);</span>
</span></span><span class="line"><span class="cl">        <span class="n">finder</span><span class="o">.</span><span class="na">setProcName</span><span class="o">(</span><span class="s">&#34;proc_up&#34;</span><span class="o">);</span>
</span></span><span class="line"><span class="cl">        <span class="n">Map</span> <span class="n">queryObjectByProc</span> <span class="o">=</span> <span class="o">(</span><span class="n">Map</span><span class="o">)</span> <span class="n">baseDemoService</span><span class="o">.</span><span class="na">queryObjectByProc</span><span class="o">(</span><span class="n">finder</span><span class="o">);</span>
</span></span><span class="line"><span class="cl">        <span class="n">System</span><span class="o">.</span><span class="na">out</span><span class="o">.</span><span class="na">println</span><span class="o">(</span><span class="n">queryObjectByProc</span><span class="o">.</span><span class="na">get</span><span class="o">(</span><span class="s">&#34;#update-count-10&#34;</span><span class="o">));</span>
</span></span><span class="line"><span class="cl">        
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl"><span class="o">}</span>
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl"><span class="c1">//@Test 调用数据库函数
</span></span></span><span class="line"><span class="cl"><span class="c1"></span><span class="kd">public</span> <span class="kt">void</span> <span class="nf">testFunction</span><span class="o">()</span> <span class="kd">throws</span> <span class="n">Exception</span><span class="o">{</span>
</span></span><span class="line"><span class="cl">        <span class="n">Finder</span> <span class="n">finder</span><span class="o">=</span><span class="k">new</span> <span class="n">Finder</span><span class="o">();</span>
</span></span><span class="line"><span class="cl">        <span class="n">finder</span><span class="o">.</span><span class="na">setFunName</span><span class="o">(</span><span class="s">&#34;fun_userId&#34;</span><span class="o">);</span>
</span></span><span class="line"><span class="cl">        <span class="n">finder</span><span class="o">.</span><span class="na">setParam</span><span class="o">(</span><span class="s">&#34;userId&#34;</span><span class="o">,</span> <span class="s">&#34;admin&#34;</span><span class="o">);</span>
</span></span><span class="line"><span class="cl">        <span class="n">String</span> <span class="n">userName</span><span class="o">=</span> <span class="n">baseDemoService</span><span class="o">.</span><span class="na">queryForObjectByByFunction</span><span class="o">(</span><span class="n">finder</span><span class="o">,</span><span class="n">String</span><span class="o">.</span><span class="na">class</span><span class="o">);</span>
</span></span><span class="line"><span class="cl">        <span class="n">System</span><span class="o">.</span><span class="na">out</span><span class="o">.</span><span class="na">println</span><span class="o">(</span><span class="n">userName</span><span class="o">);</span>
</span></span><span class="line"><span class="cl"><span class="o">}</span>
</span></span><span class="line"><span class="cl">
</span></span></code></pre></td></tr></table>
</div>
</div>','## 微服务
 6.0.0 项目入口是springrain-system-web,基于Istio实现微服务,正在整理文档.
## 前后分离
 6.0.0 基于VUE前后端分离,使用JWT认证.前端项目','<ul>
<li>
<ul>
<li>
<a href="#%E5%BE%AE%E6%9C%8D%E5%8A%A1">微服务</a></li>
<li>
<a href="#%E5%89%8D%E5%90%8E%E5%88%86%E7%A6%BB">前后分离</a></li>
<li>
<a href="#%E5%AE%9E%E7%8E%B0%E4%BA%86%E4%BB%80%E4%B9%88?">实现了什么?</a></li>
<li>
<a href="#%E5%AE%9E%E7%8E%B0%E6%80%9D%E8%B7%AF">实现思路</a></li>
<li>
<a href="#%E9%99%90%E5%88%B6">限制</a></li>
<li>
<a href="#%E4%BD%93%E9%AA%8C%E5%8D%95%E4%BD%93%E5%88%B0%E5%88%86%E5%B1%82%E5%88%87%E6%8D%A2">体验单体到分层切换</a></li>
<li>
<a href="#%E5%8D%9A%E5%AE%A2">博客</a></li>
<li>
<a href="#%E6%96%87%E6%A1%A3">文档</a></li>
<li>
<a href="#%E4%BB%A3%E7%A0%81%E7%94%9F%E6%88%90%E5%99%A8">代码生成器</a></li>
<li>
<a href="#sql%E8%84%9A%E6%9C%AC">sql脚本</a></li>
<li>
<a href="#%E6%A1%88%E4%BE%8B">案例</a></li>
</ul>
</li>
</ul>
','springrain','springrain',NULL,'web','web',NULL,NULL,NULL,NULL,'springrain项目说明',NULL,'springrain');

