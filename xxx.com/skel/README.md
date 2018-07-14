# skel just for demo

## 依赖
    wsp: github.com/simplejia/wsp
    lc: github.com/simplejia/lc
    lm: github.com/simplejia/lm
    clog: github.com/simplejia/clog
    utils: github.com/simplejia/utils
    namecli: github.com/simplejia/namecli
    mongo: gopkg.in/mgo.v2
    redis: github.com/garyburd/redigo/redis
    gomvpkg: golang.org/x/tools/cmd/gomvpkg

## 注意
    如果在controller里修改了路由，编译前需执行go generate
    如果想在本机启动服务，然后本机也没有安装mongo服务端，请删除./mongo/skel.json
    如果要新生成一个项目，请执行./new.sh $proj, $proj替换成你的新项目名，执行完后，当前目录会被替换成新项目
    如果要修改xxx.com，换成新的域名，比如：yyy.zzz，请通过全局替换，处理xxx.com目录下的所有文件，只要是出现了xxx.com，全换成yyy.zzz，最后再把xxx.com这个目录替换成yyy.zzz
