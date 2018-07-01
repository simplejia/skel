package skel

import (
	"xxx.com/lib"
	"xxx.com/skel/api"
)

// UpdateIndex 定义订阅者响应方法
func (skel *Skel) UpdateIndex(d interface{}, t lib.UpdateType, props map[string]interface{}) {
	switch o := d.(type) {
	case *api.Skel:
		switch t {
		case lib.ADD:
			_ = o
		case lib.DELETE:
			_ = o
		case lib.GET:
			_ = o
		}
	}
}

func init() {
	lib.RegisterIndex(&Skel{})
}
