# 说明

一个方便切换分支和执行脚本的web小工具。部署在测试服务器上，方便测试同学拉取切换分支测试，以及执行一些脚本。

### 配置文件

config.ini

```
[projects]
project1 = /data/htdocs/project_1
project2 = /data/htdocs/project_2

[commands]
reload crontab = crontab /etc/support/crontab.conf
restart supervisor = supervisorctl restart all
```

- projects 是项目配置，配置多个项目的名称和源码地址。
- commands 是项目依赖的一些服务的执行命令配置。

#### 编译

linux:

```
$ export GOARCH=amd64
$ export GOOS=linux
$ go build
```

![screenshot](https://github.com/lisijie/repo_tool/raw/master/screenshot.png)