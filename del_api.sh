# 生成controller,api代码

path=$1
if [ ! $path ];then
    echo "no path provided"
    exit
fi

del_api -path=$path
