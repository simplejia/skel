# 删除model,service,controller,api各层代码

curdir=`pwd`
basedir=`dirname $curdir`
pkgpath=${basedir##*src}

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

del_crud -pkg=$pkgpath -name=$name
