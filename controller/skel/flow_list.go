package skel

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/skel/service"
	"github.com/simplejia/skel_api"
	"github.com/simplejia/utils"

	clog "github.com/simplejia/clog/api"
)

// FlowList just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func (skel *Skel) FlowList(w http.ResponseWriter, r *http.Request) {
	fun := "controller.skel.Skel.FlowList"

	var req *skel_api.SkelFlowListReq
	if err := json.Unmarshal(skel.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, utils.CodePara)
		return
	}

	trace := utils.GetTrace(skel)

	limitMore := req.Limit + 1

	skelsAPI, err := service.NewSkel().WithTrace(trace).FlowList(req.LastID, limitMore)
	if err != nil {
		clog.Error("%s skel.FlowList err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, utils.CodeSrv)
		return
	}

	n := len(skelsAPI)
	if n == 0 {
		skel.ReplyOk(w, nil)
		return
	}

	more := false
	if n == limitMore {
		more = true
		skelsAPI = skelsAPI[:req.Limit]
	}

	resp := &skel_api.SkelFlowListResp{
		List: skelsAPI,
		More: more,
	}
	skel.ReplyOk(w, resp)

	return
}
