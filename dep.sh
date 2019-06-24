# 下载依赖组件脚本

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


go get -u -v github.com/simplejia/wsp
go get -u -v golang.org/x/tools/cmd/gomvpkg