package skel

import (
	"github.com/simplejia/lib"
	"github.com/simplejia/skel/model"
)

// Del 定义删除操作
func (skel *Skel) Del(id int64) (err error) {
	fun := "service.skel.Skel.Del"
	defer lib.TraceMe(skel.Trace, fun)()

	if err = model.NewSkel().WithTrace(skel.Trace).Del(id); err != nil {
		return
	}

	return
}
