package skel

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/lib"
	"github.com/simplejia/skel/service"
	"github.com/simplejia/skel_api"

	clog "github.com/simplejia/clog/api"
)

// Get just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func (skel *Skel) Get(w http.ResponseWriter, r *http.Request) {
	fun := "controller.skel.Skel.Get"

	var req *skel_api.SkelGetReq
	if err := json.Unmarshal(skel.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, lib.CodePara)
		return
	}

	trace := lib.GetTrace(skel)

	skelAPI, err := service.NewSkel().WithTrace(trace).Get(req.ID)
	if err != nil {
		clog.Error("%s skel.Get err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := (*skel_api.SkelGetResp)(skelAPI)
	skel.ReplyOk(w, resp)

	return
}
