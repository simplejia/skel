# 编译脚本

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

env=""
test=""

while test $# -gt 0; do
    case "$1" in
        -env)
            shift
            env=$1
            shift
            ;;
        -test)
            shift
            test="true"
            ;;
        *)
            break
            ;;
    esac
done

if [ ! $env ];then
    env="prod"
fi

basedir=`dirname $this`
cd $basedir
curdir=`pwd`
pkgpath=${curdir##*src/}

controllers=$(find controller -type d -exec basename {} \;|grep -v controller)

if [ $env == "dev" ];then
    go generate $pkgpath
    go build $pkgpath

    if [ $test ];then
        for v in ${controllers[@]}; do
            go test -c $pkgpath/controller/$v
        done
    fi
else
    export GOOS=linux
    go generate $pkgpath
    go build $pkgpath

    if [ $test ];then
        for v in ${controllers[@]}; do
            go test -c $pkgpath/controller/$v
        done
    fi
fi
