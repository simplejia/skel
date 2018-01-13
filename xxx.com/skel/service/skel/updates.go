package skel

import (
	"reflect"

	"xxx.com/lib"
	"xxx.com/skel/api"
)

func (skel *Skel) UpdateIndex(d interface{}, t lib.UpdateType, props map[string]interface{}) {
	fun := "skel.Skel.UpdateIndex"
	_ = fun

	switch reflect.TypeOf(d) {
	case reflect.TypeOf(&api.Skel{}):
		skel := d.(*api.Skel)
		_ = skel
		switch t {
		case lib.ADD:
		case lib.DELETE:
		}
	}
}

func init() {
	lib.RegisterIndex(&Skel{})
}
