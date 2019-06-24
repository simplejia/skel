# 单元测试启动脚本

this="$0"
while [ -h "$this" ]; do
    ls=`ls -ld "$this"`
    link=`expr "$ls" : '.*-> \(.*\)$'`
    if expr "$link" : '.*/.*' > /dev/null; then
        this="$link"
    else
        this=`dirname "$this"`/"$link"
    fi
done

env=$1
if [ ! $env ];then
    env="prod"
fi

basedir=`dirname $this`
cd $basedir
curdir=`pwd`
pkgpath=${curdir##*src/}

controllers=$(find controller -type d -exec basename {} \;|grep -v controller)

if [ $env == "dev" ];then
    for v in ${controllers[@]}; do
        go test $pkgpath/controller/$v -v -test=true -env=dev 
    done
else
    for v in ${controllers[@]}; do
        ./$v.test -test.v -test=true
    done
fi
