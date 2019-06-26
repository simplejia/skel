# skel: 简单实用的golang webserver框架，skel取自单词：skeleton

## Installation
> go get -u github.com/simplejia/skel

以下内容无特殊说明均是在skel项目目录下执行。

```
依赖安装：
运行：./dep.sh
运行：./install_tools.sh (注意，默认安装工具程序到$GOBIN目录下，请设置$PATH环境变量)

以skel项目为模板生成新项目：demo
运行：./new.sh demo
运行完以上命令后，会在$GOPATH/src目录下新建demo和demo_api两个文件夹


以skel项目为模板生成新项目：xxx.com/demo
运行：./new.sh xxx.com/demo
运行完以上命令后，会在$GOPATH/src目录下新建xxx.com/demo和xxx.com/demo_api两级目录结构


以下内容无特殊说明均是在demo项目目录下执行。

生成user相关的controller, service, model层代码：
运行：./gen_crud.sh -name user (其他功能请直接执行：./gen_crud.sh，查看输出命令列表)

生成/user/get接口：
运行：./gen_api.sh /user/get

del_crud.sh和del_api.sh用于删除操作
```