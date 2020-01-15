package skel

import (
	"github.com/simplejia/skel/model"
	"github.com/simplejia/skel_api"
	"github.com/simplejia/utils"
)

// PageList 定义page_list操作
func (skel *Skel) PageList(offset, limit int) (skelsAPI []*skel_api.Skel, err error) {
	fun := "service.skel.Skel.PageList"
	defer utils.TraceMe(skel.Trace, fun)()

	if skelsAPI, err = model.NewSkel().WithTrace(skel.Trace).PageList(offset, limit); err != nil {
		return
	}

	return
}

// FlowList 定义list操作
func (skel *Skel) FlowList(lastID string, limit int) (skelsAPI []*skel_api.Skel, err error) {
	fun := "service.skel.Skel.FlowList"
	defer utils.TraceMe(skel.Trace, fun)()

	if skelsAPI, err = model.NewSkel().WithTrace(skel.Trace).FlowList(lastID, limit); err != nil {
		return
	}

	return
}
