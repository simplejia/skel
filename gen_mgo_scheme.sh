# 生成mongo scheme

env=$1
if [ ! $env ];then
    env="prod"
fi

gen_mgo_scheme -env $env
