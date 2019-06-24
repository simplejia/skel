# 安装tools

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

go install $pkgpath/skel/tools/gen_api
go install $pkgpath/skel/tools/gen_crud
go install $pkgpath/skel/tools/del_api
go install $pkgpath/skel/tools/del_crud
go install $pkgpath/skel/tools/gen_mgo_scheme
