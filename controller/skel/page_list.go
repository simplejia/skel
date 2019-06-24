package skel

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/lib"
	"github.com/simplejia/skel/service"
	"github.com/simplejia/skel_api"

	clog "github.com/simplejia/clog/api"
)

// PageList just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func (skel *Skel) PageList(w http.ResponseWriter, r *http.Request) {
	fun := "controller.skel.Skel.PageList"

	var req *skel_api.SkelPageListReq
	if err := json.Unmarshal(skel.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, lib.CodePara)
		return
	}

	trace := lib.GetTrace(skel)

	limitMore := req.Limit + 1

	skelsAPI, err := service.NewSkel().WithTrace(trace).PageList(req.Offset, limitMore)
	if err != nil {
		clog.Error("%s skel.PageList err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, lib.CodeSrv)
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

	resp := &skel_api.SkelPageListResp{
		List:   skelsAPI,
		Offset: req.Offset + len(skelsAPI),
		More:   more,
	}
	skel.ReplyOk(w, resp)

	return
}
