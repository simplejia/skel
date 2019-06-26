# 编译脚本

env=$1

if [ $env -a $env == "dev" ];then
    go generate $pkgpath && go build $pkgpath
else
    export GOOS=linux && go generate $pkgpath && go build $pkgpath
fi
