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

env=$1
if [ ! $env ];then
    env="prod"
fi


proj=$(basename $PWD)
controllers=$(find controller -type d -exec basename {} \;|grep -v controller)

if [ $env == "dev" ];then
    go generate xxx.com/$proj
    go build xxx.com/$proj

    for v in ${controllers[@]}; do
        go test -c xxx.com/$proj/controller/$v
    done
else
    export GOOS=linux
    go generate xxx.com/$proj
    go build xxx.com/$proj

    for v in ${controllers[@]}; do
        go test -c xxx.com/$proj/controller/$v
    done
fi
