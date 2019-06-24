package skel

import (
	"github.com/simplejia/lib"
)

// Del 定义删除操作
func (skel *Skel) Del(id int64) (err error) {
	fun := "model.skel.Skel.Del"
	defer lib.TraceMe(skel.Trace, fun)()

	c := skel.GetC()
	defer c.Database.Session.Close()

	err = c.RemoveId(id)
	if err != nil {
		return
	}

	return
}
