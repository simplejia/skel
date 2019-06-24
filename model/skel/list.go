package skel

import (
	"github.com/globalsign/mgo/bson"
	"github.com/simplejia/lib"
	"github.com/simplejia/skel_api"
)

// PageList 定义page_list操作
func (skel *Skel) PageList(offset, limit int) (skelsAPI []*skel_api.Skel, err error) {
	fun := "model.skel.Skel.PageList"
	defer lib.TraceMe(skel.Trace, fun)()

	c := skel.GetC()
	defer c.Database.Session.Close()

	q := bson.M{}
	err = c.Find(q).Sort("-_id").Skip(offset).Limit(limit).All(&skelsAPI)
	if err != nil {
		return
	}

	return
}

// FlowList 定义flow_list操作
func (skel *Skel) FlowList(id int64, limit int) (skelsAPI []*skel_api.Skel, err error) {
	fun := "model.skel.Skel.FlowList"
	defer lib.TraceMe(skel.Trace, fun)()

	c := skel.GetC()
	defer c.Database.Session.Close()

	q := bson.M{}
	if id != *new(int64) {
		q["_id"] = bson.M{
			"$lt": id,
		}
	}
	err = c.Find(q).Sort("-_id").Limit(limit).All(&skelsAPI)
	if err != nil {
		return
	}

	return
}
