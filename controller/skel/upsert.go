package skel

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/lib"
	"github.com/simplejia/skel/service"
	"github.com/simplejia/skel_api"

	clog "github.com/simplejia/clog/api"
)

// Upsert just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func (skel *Skel) Upsert(w http.ResponseWriter, r *http.Request) {
	fun := "controller.skel.Skel.Upsert"

	var req *skel_api.SkelUpsertReq
	if err := json.Unmarshal(skel.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, lib.CodePara)
		return
	}

	trace := lib.GetTrace(skel)

	skelAPI := (*skel_api.Skel)(req)

	if err := service.NewSkel().WithTrace(trace).Upsert(skelAPI); err != nil {
		clog.Error("%s skel.Upsert err: %v, req: %v", fun, err, req)
		skel.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &skel_api.SkelUpsertResp{}
	skel.ReplyOk(w, resp)

	// 进行一些异步处理的工作
	go lib.Updates(skelAPI, lib.ADD, nil)

	return
}
