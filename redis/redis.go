/*
Package redis 代表redis连接配置，所有redis连接配置均在此目录配置，一个redis对应一个配置文件。
*/
package redis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/simplejia/skel/conf"

	"github.com/garyburd/redigo/redis"
	"github.com/simplejia/clog/api"
	"github.com/simplejia/utils"
)

// Conf 用于redis连接配置
type Conf struct {
	Addr string // addr
	Auth string // auth
	Db   int    // db select
}

var (
	// RDS 表示redis连接，key是业务名，value是redis连接
	RDS map[string]*redis.Pool = map[string]*redis.Pool{}
)

func init() {
	dir := "redis"
	for i := 0; i < 3; i++ {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			break
		}
		dir = filepath.Join("..", dir)
	}
	err := filepath.Walk(
		dir,
		func(path string, info os.FileInfo, err error) (reterr error) {
			if err != nil {
				reterr = err
				return
			}
			if info.IsDir() {
				return
			}
			if strings.HasPrefix(filepath.Base(path), ".") {
				return
			}
			if filepath.Ext(path) != ".json" {
				return
			}

			fcontent, err := ioutil.ReadFile(path)
			if err != nil {
				reterr = err
				return
			}
			fcontent = utils.RemoveAnnotation(fcontent)
			var envs map[string]*Conf
			if err := json.Unmarshal(fcontent, &envs); err != nil {
				reterr = err
				return
			}

			c := envs[conf.Env]
			if c == nil {
				reterr = fmt.Errorf("env not right: %s", conf.Env)
				return
			}

			rd := &redis.Pool{
				MaxIdle:     30,
				IdleTimeout: time.Minute,
				Dial: func() (conn redis.Conn, err error) {
					conn, err = redis.Dial("tcp", c.Addr,
						redis.DialReadTimeout(time.Second),
						redis.DialConnectTimeout(time.Second),
					)
					if err != nil {
						clog.Error("redis.Dial err: %v, req: %v", err, c.Addr)
						return
					}

					if auth := c.Auth; auth != "" {
						if _, err = conn.Do("AUTH", auth); err != nil {
							clog.Error("redis AUTH err: %v, req: %v,%v", err, c.Addr, auth)
							conn.Close()
							return
						}
					}
					if db := c.Db; db > 0 {
						if _, err = conn.Do("SELECT", db); err != nil {
							clog.Error("redis SELECT err: %v, req: %v,%v", err, c.Addr, db)
							conn.Close()
							return
						}
					}
					return
				},
			}

			key := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			RDS[key] = rd
			return
		},
	)
	if err != nil {
		log.Printf("conf(redis) not right: %v\n", err)
		os.Exit(-1)
	}

	log.Printf("conf(redis): %v\n", RDS)
}
