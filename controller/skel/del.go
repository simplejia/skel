package skel

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/skel_api"

	"github.com/simplejia/skel/service"
	"github.com/simplejia/utils"

	clog "github.com/simplejia/clog/api"
)

// Del just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func (skel *Skel) Del(w http.ResponseWriter, r *http.Request) {
	fun := "controller.skel.Skel.Del"

	var req *skel_api.SkelDelReq
	if err := json.Unmarshal(skel.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, utils.CodePara)
		return
	}

	trace := utils.GetTrace(skel)

	if err := service.NewSkel().WithTrace(trace).Del(req.ID); err != nil {
		clog.Error("%s skel.Del err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, utils.CodeSrv)
		return
	}

	resp := &skel_api.SkelDelResp{}
	skel.ReplyOk(w, resp)

	skelAPI := &skel_api.Skel{ID: req.ID}

	// 进行一些异步处理的工作
	go utils.Updates(skelAPI, utils.DELETE, nil)

	return
}
