package skel

import (
	"github.com/simplejia/lib"
	"github.com/simplejia/skel_api"
)

// Update 定义更新操作
func (skel *Skel) Update(skelAPI *skel_api.Skel) (err error) {
	fun := "model.skel.Skel.Update"
	defer lib.TraceMe(skel.Trace, fun)()

	c := skel.GetC()
	defer c.Database.Session.Close()

	err = c.UpdateId(skelAPI.ID, skelAPI)
	if err != nil {
		return
	}

	return
}

// Upsert 定义upsert操作
func (skel *Skel) Upsert(skelAPI *skel_api.Skel) (err error) {
	fun := "model.skel.Skel.Upsert"
	defer lib.TraceMe(skel.Trace, fun)()

	c := skel.GetC()
	defer c.Database.Session.Close()

	_, err = c.UpsertId(skelAPI.ID, skelAPI)
	if err != nil {
		return
	}

	return
}
