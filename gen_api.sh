# 生成controller,api代码

curdir=`pwd`
basedir=`dirname $curdir`
pkgpath=${basedir##*src}

path=$1
if [ ! $path ];then
    echo "no path provided"
    exit
fi

gen_api -pkg=$pkgpath -path=$path
