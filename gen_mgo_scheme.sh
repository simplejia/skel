# 生成mongo scheme

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

env=$1
if [ ! $env ];then
    echo "no env provided"
    exit
fi

gen_mgo_scheme -env $env
