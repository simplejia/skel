// Package test 用于封装测试用例用到的公共库
package test

import (
	"time"

	"xxx.com/lib"
	"xxx.com/skel/conf"
)

const (
	IDStart = 10
)

var (
	ID int64 = IDStart
)

func GetID() (id int64) {
	id = ID
	ID++
	return
}

func Setup() {
	ID = IDStart
}

// Sleep 主要用于等待异步执行过程结束，目前至少要大于网络+db执行时间
func Sleep() {
	if conf.Env == lib.DEV {
		time.Sleep(time.Millisecond * 10)
	} else {
		time.Sleep(time.Millisecond * 200)
	}
}
