# 生成model,service,controller,api各层代码

curdir=`pwd`
basedir=`dirname $curdir`
pkgpath=${basedir##*src}

args=''

while test $# -gt 0; do
    case "$1" in
        -name)
            shift
            args+=" -name="$1
            shift
            ;;
        -id_type)
            shift
            args+=" -id_type="$1
            shift
            ;;
        -keys)
            shift
            args+=" -keys="$1
            shift
            ;;
        -need_multi_table)
            args+=" -need_multi_table"
            shift
            ;;
        -db_num)
            shift
            args+=" -db_num="$1
            shift
            ;;
        -table_num)
            shift
            args+=" -table_num="$1
            shift
            ;;
        -connect_timeout)
            shift
            args+=" -connect_timeout="$1
            shift
            ;;
        -timeout)
            shift
            args+=" -timeout="$1
            shift
            ;;
        *)
            break
            ;;
    esac
done

gen_crud -pkg=$pkgpath $args
