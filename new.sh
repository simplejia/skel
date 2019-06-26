# 生成新项目脚本

proj=$1
if [ ! $proj ];then
    echo "no proj name"
    exit
fi

curdir=`pwd`
basedir=`dirname $curdir`
pkg=`basename $proj`

mkdir -p ${basedir%%/src*}/src/$(dirname $proj)

echo "begin generate new project..."
(gomvpkg -from github.com/simplejia/skel -to $proj && \
sed "s/package $pkg/package main/g" <main.go >main.go.new && mv main.go.new main.go && \
sed "s/package $pkg/package main/g" <WSP.go >WSP.go.new && mv WSP.go.new WSP.go && \
rm -rf tools skel new.sh install_tools.sh .git .gitignore ../${pkg}_api/.git && \
gomvpkg -from github.com/simplejia/skel_api -to ${proj}_api && \
rm -rf ../${pkg}_api/.git)

echo "begin download dependence..."
go get -v github.com/simplejia/skel
