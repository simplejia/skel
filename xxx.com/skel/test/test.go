// Package test 用于封装测试用例用到的公共库
package test

import (
	"time"

	"github.com/simplejia/lc"
	"xxx.com/lib"

	"xxx.com/skel/conf"
)

const (
	// IDStart 定义起始id
	IDStart = 1
)

var (
	// ID 记录测试用id
	ID int64 = IDStart
)

// GetID 返回id
func GetID() (id int64) {
	id = ID
	ID++
	return
}

// Setup 用于每一个测试用例的初始化
func Setup() {
	ID = IDStart
	// 运行测试用例时关闭lc的使用，避免由于缓存数据导致数据误判
	lc.Disabled = true
}

// Sleep 主要用于等待异步执行过程结束，目前至少要大于网络+db执行时间
func Sleep() {
	if conf.Env == lib.DEV {
		time.Sleep(time.Millisecond * 10)
	} else {
		time.Sleep(time.Millisecond * 200)
	}
}
