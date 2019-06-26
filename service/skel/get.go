package skel

import (
	"github.com/simplejia/skel/model"
	"github.com/simplejia/skel_api"
	"github.com/simplejia/utils"
)

// Get 定义获取操作
func (skel *Skel) Get(id int64) (skelAPI *skel_api.Skel, err error) {
	fun := "service.skel.Skel.Get"
	defer utils.TraceMe(skel.Trace, fun)()

	if skelAPI, err = model.NewSkel().WithTrace(skel.Trace).Get(id); err != nil {
		return
	}

	return
}

// Gets 定义批量获取操作
func (skel *Skel) Gets(ids []int64) (skelsAPI map[int64]*skel_api.Skel, err error) {
	fun := "service.skel.Skel.Gets"
	defer utils.TraceMe(skel.Trace, fun)()

	skelsSliceAPI, err := model.NewSkel().WithTrace(skel.Trace).Gets(ids)
	if err != nil {
		return
	}

	if len(skelsSliceAPI) == 0 {
		return
	}

	skelsAPI = map[int64]*skel_api.Skel{}
	for _, skelAPI := range skelsSliceAPI {
		skelsAPI[skelAPI.ID] = skelAPI
	}

	return
}
