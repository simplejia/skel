package skel

import (
	"github.com/simplejia/skel/model"
	"github.com/simplejia/skel_api"
	"github.com/simplejia/utils"
)

// Add 定义新增操作
func (skel *Skel) Add(skelAPI *skel_api.Skel) (err error) {
	fun := "service.skel.Skel.Add"
	defer utils.TraceMe(skel.Trace, fun)()

	if err = model.NewSkel().WithTrace(skel.Trace).Add(skelAPI); err != nil {
		return
	}

	return
}
