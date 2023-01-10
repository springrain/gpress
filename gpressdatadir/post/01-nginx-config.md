---
title: "01.Nginx安装配置"
date: 2017-01-10T17:11:00+08:00
draft: false
tags: ["web"]
categories: ["web"]
author: "springrain"
---

## 1. 准备环境
```shell
yum -y install openssl-devel zlib-devel libtool automake autoconf make  
yum -y install gcc gcc-c++  

#下载Nginx   
wget http://nginx.org/download/nginx-1.22.1.tar.gz
#下载PCRE 
wget https://jaist.dl.sourceforge.net/project/pcre/pcre/8.45/pcre-8.45.tar.gz
```
## 2. 安装nginx
```shell

tar -zxf nginx-1.22.1.tar.gz
tar -zxf pcre-8.45.tar.gz

#进入nginx解压后的目录
cd nginx-1.22.1

#执行检查
./configure --with-http_stub_status_module --with-http_ssl_module --with-http_sub_module --with-http_v2_module --with-http_realip_module --with-pcre=/root/pcre-8.45

#然后执行:
make
make install
```
运行一下,查看是否正常,Nginx默认监听的是80端口  
启动:```/usr/local/nginx/sbin/nginx```  
重启:```/usr/local/nginx/sbin/nginx -s reload```  

![nginx](/01/01-nginx-config-01.jpg)  

**[Windows版本的Nginx](/01/nginx-windows.zip)**  

## 3. 配置Nginx
一般是把nginx主文件和server文件分开处理,这样便于管理.  
主文件是:```/usr/local/nginx/conf/nginx.conf```    
给server文件创建一个单独的www文件夹 
```shell 
mkdir /usr/local/nginx/www  
```
server配置文件名以 .conf为后缀.  
  

**[下载范例配置文件](/01/conf.zip)**  


## 4. 开机启动
编辑 ```vi /etc/rc.d/rc.local``` 在最后加上一行    
```/usr/local/nginx/sbin/nginx```   
建议把nginx目录迁移至数据盘.例如数据库挂载为/data目录    
关闭nginx服务:```/usr/local/nginx/sbin/nginx -s stop```   
迁移目录:```mv /usr/local/nginx /data/nginx```  
建立软链接:```ln -s /usr/local/nginx /usr/local/nginx```   
启动nginx服务:```/usr/local/nginx/sbin/nginx```    

## 5. 切割日志
Nginx默认并没有实现日志切割,这样所有的日志都在一个文件里,文件很大时会影响访问的性能.通过每天零点调用日志切割的脚本,实现Nginx日志切割.

cut_nginx_logs.sh 脚本如下:
```shell
#!/bin/bash
#function:cut nginx log files 
#author: http://www.jiagou.com

###设置日志文件的路径####
log_files_path="/usr/local/nginx/logs/"
###日志文件备份的路径####
log_files_dir=${log_files_path}history/$(date -d "yesterday" +"%Y")/$(date -d "yesterday" +"%m")
###设置日志文件名称###
log_files_name=(domain1.access domain2.access error)
###nginx主程序###
nginx_sbin="/usr/local/nginx/sbin/nginx"
###保留日志天数###
save_days=100

############################################
#Please do not modify the following script #
############################################
mkdir -p $log_files_dir

log_files_num=${#log_files_name[@]}

#cut nginx log files
for((i=0;i<$log_files_num;i++));do
mv ${log_files_path}${log_files_name[i]}.log ${log_files_dir}/${log_files_name[i]}_$(date -d "yesterday" +"%Y%m%d").log
done

#delete save_days ago nginx log files
find $log_files_path -mtime +$save_days -exec rm -rf {} \; 

##reload会重启nginx,重新加载配置文件,reopen只会重建日志文件##
#$nginx_sbin -s reload
$nginx_sbin -s reopen
```
把cut_nginx_logs.sh脚本放到```/usr/local/nginx/sbin/``` 目录下,修改可执行权限:
```shell
chmod 755 /usr/local/nginx/sbin/cut_nginx_logs.sh
```
添加定时调用:crontab -e 编辑输入   
```shell
00 00 * * * /bin/bash /usr/local/nginx/sbin/cut_nginx_logs.sh
```


