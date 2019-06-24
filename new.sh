# 生成新项目脚本

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

proj=$1
if [ ! $proj ];then
    echo "no proj name"
    exit
fi

basedir=`dirname $this`
cd $basedir
curdir=`pwd`
basedir=`dirname $curdir`
pkgpath=${basedir##*src/}
pkg=`basename $proj`

if [ $pkg == $proj ];then
    proj=$pkgpath/$proj
else
    mkdir -p ${basedir%%/src*}/src/$(dirname $proj)
fi

echo "begin generate new project..."
(gomvpkg -from $pkgpath/skel -to $proj && \
sed "s/package $pkg/package main/g" <main.go >main.go.new && mv main.go.new main.go && \
sed "s/package $pkg/package main/g" <WSP.go >WSP.go.new && mv WSP.go.new WSP.go && \
rm -rf tools skel new.sh install_tools.sh .git .gitignore ../${pkg}_api/.git && \
gomvpkg -from $pkgpath/skel_api -to ${proj}_api && \
rm -rf ../${pkg}_api/.git)

echo "begin download dependence..."
go get -v github.com/simplejia/skel
