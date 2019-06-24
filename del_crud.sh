# 删除model,service,controller,api各层代码

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

basedir=`dirname $this`
cd $basedir
curdir=`pwd`
basedir=`dirname $curdir`
pkgpath=${basedir##*src/}

while test $# -gt 0; do
    case "$1" in
        -name)
            shift
            ;;
        *)
            break
            ;;
    esac
done

name=$1
if [ ! $name ];then
    echo "no name provided"
    exit
fi

del_crud -pkg $pkgpath -name $name
