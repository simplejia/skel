package skel

import (
	"github.com/simplejia/lib"
	"github.com/simplejia/skel_api"
)

// Add 定义新增操作
func (skel *Skel) Add(skelAPI *skel_api.Skel) (err error) {
	fun := "model.skel.Skel.Add"
	defer lib.TraceMe(skel.Trace, fun)()

	c := skel.GetC()
	defer c.Database.Session.Close()

	err = c.Insert(skelAPI)
	if err != nil {
		return
	}

	return
}
