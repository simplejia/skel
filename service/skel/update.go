package skel

import (
	"github.com/simplejia/skel/model"
	"github.com/simplejia/skel_api"
	"github.com/simplejia/utils"
)

// Update 定义更新操作
func (skel *Skel) Update(skelAPI *skel_api.Skel) (err error) {
	fun := "service.skel.Skel.Update"
	defer utils.TraceMe(skel.Trace, fun)()

	if err = model.NewSkel().WithTrace(skel.Trace).Update(skelAPI); err != nil {
		return
	}

	return
}

// Upsert 定义upsert操作
func (skel *Skel) Upsert(skelAPI *skel_api.Skel) (err error) {
	fun := "service.skel.Skel.Upsert"
	defer utils.TraceMe(skel.Trace, fun)()

	if err = model.NewSkel().WithTrace(skel.Trace).Upsert(skelAPI); err != nil {
		return
	}

	return
}
