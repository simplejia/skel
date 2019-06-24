package skel

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/simplejia/lib"
	"github.com/simplejia/skel_api"
)

// Get 定义获取操作
func (skel *Skel) Get(id int64) (skelAPI *skel_api.Skel, err error) {
	fun := "model.skel.Skel.Get"
	defer lib.TraceMe(skel.Trace, fun)()

	c := skel.GetC()
	defer c.Database.Session.Close()

	err = c.FindId(id).One(&skelAPI)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}
		err = nil
		return
	}

	return
}

// Gets 定义批量获取操作
func (skel *Skel) Gets(ids []int64) (skelsAPI []*skel_api.Skel, err error) {
	fun := "model.skel.Skel.Gets"
	defer lib.TraceMe(skel.Trace, fun)()

	c := skel.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	err = c.Find(q).All(&skelsAPI)
	if err != nil {
		return
	}

	return
}
