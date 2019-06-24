# 生成controller,api代码

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

path=$1
if [ ! $path ];then
    echo "no path provided"
    exit
fi

gen_api -pkg $pkgpath -path $path
