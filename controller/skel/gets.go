package skel

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/lib"
	"github.com/simplejia/skel/service"
	"github.com/simplejia/skel_api"

	clog "github.com/simplejia/clog/api"
)

// Gets just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func (skel *Skel) Gets(w http.ResponseWriter, r *http.Request) {
	fun := "controller.skel.Skel.Gets"

	var req *skel_api.SkelGetsReq
	if err := json.Unmarshal(skel.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, lib.CodePara)
		return
	}

	trace := lib.GetTrace(skel)

	skelsAPI, err := service.NewSkel().WithTrace(trace).Gets(req.IDS)
	if err != nil {
		clog.Error("%s skel.Gets err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := skel_api.SkelGetsResp(skelsAPI)
	skel.ReplyOk(w, resp)

	return
}
