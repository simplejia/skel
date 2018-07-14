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

# configuration variables
basedir=`dirname $this`
cd $basedir
basedir=`pwd`
dirname=`basename $(dirname $basedir)`
basename=`basename $basedir`

proj=$1
if [ ! $proj ];then
    echo "no proj name"
    exit
fi

gomvpkg -from $dirname/$basename -to $dirname/$proj && \
cd ../$proj && \
sed "s/package $proj/package main/g" <main.go >main.go.new && mv main.go.new main.go && \
sed "s/package $proj/package main/g" <WSP.go >WSP.go.new && mv WSP.go.new WSP.go
