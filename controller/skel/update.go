package skel

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/skel/service"
	"github.com/simplejia/skel_api"
	"github.com/simplejia/utils"

	clog "github.com/simplejia/clog/api"
)

// Update just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func (skel *Skel) Update(w http.ResponseWriter, r *http.Request) {
	fun := "controller.skel.Skel.Update"

	var req *skel_api.SkelUpdateReq
	if err := json.Unmarshal(skel.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, utils.CodePara)
		return
	}

	trace := utils.GetTrace(skel)

	skelAPI := (*skel_api.Skel)(req)
	if err := service.NewSkel().WithTrace(trace).Update(skelAPI); err != nil {
		clog.Error("%s skel.Update err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, utils.CodeSrv)
		return
	}

	resp := &skel_api.SkelUpdateResp{}
	skel.ReplyOk(w, resp)

	// 进行一些异步处理的工作
	go utils.Updates(skelAPI, utils.UPDATE, nil)

	return
}
