package skel

import (
	"github.com/globalsign/mgo/bson"
	"github.com/simplejia/skel_api"
	"github.com/simplejia/utils"
)

// PageList 定义page_list操作
func (skel *Skel) PageList(offset, limit int) (skelsAPI []*skel_api.Skel, err error) {
	fun := "model.skel.Skel.PageList"
	defer utils.TraceMe(skel.Trace, fun)()

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
func (skel *Skel) FlowList(lastID string, limit int) (skelsAPI []*skel_api.Skel, err error) {
	fun := "model.skel.Skel.FlowList"
	defer utils.TraceMe(skel.Trace, fun)()

	c := skel.GetC()
	defer c.Database.Session.Close()

	q := bson.M{}
	err = c.Find(q).Sort("-_id").Limit(limit).All(&skelsAPI)
	if err != nil {
		return
	}

	return
}
