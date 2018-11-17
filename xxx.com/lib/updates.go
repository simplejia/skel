// 此文件定义一种订阅，发布模式

package lib

import (
	"runtime/debug"

	"github.com/simplejia/clog/api"
)

// UpdateType 表示类型
type UpdateType int

const (
	// ADD 代表insert和upsert
	ADD UpdateType = iota
	// DELETE 代表delete
	DELETE
	// UPDATE 代表update
	UPDATE
	// GET 代表get
	GET
)

var (
	indices = make([]Index, 0)
)

type Index interface {
	UpdateIndex(d interface{}, t UpdateType, props map[string]interface{})
}

func Updates(d interface{}, t UpdateType, props map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			clog.Error("lib.Updates recover err: %v, stack: %s", err, debug.Stack())
		}
	}()
	for _, index := range indices {
		index.UpdateIndex(d, t, props)
	}
}

func RegisterIndex(index Index) {
	indices = append(indices, index)
}
