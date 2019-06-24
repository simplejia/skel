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

path=$1
if [ ! $path ];then
    echo "no path provided"
    exit
fi

del_api -path $path
