package skel

import (
	"encoding/json"
	"net/http"

	"xxx.com/lib"
	"xxx.com/skel/api"
	"xxx.com/skel/service"

	"github.com/simplejia/clog"
)

type GetReq struct {
	ID int64 `json:"id"`
}

// Regular 用于参数校验
func (getReq *GetReq) Regular() (ok bool) {
	if getReq == nil {
		return
	}

	if getReq.ID <= 0 {
		return
	}

	ok = true
	return
}

type GetRsp struct {
	Skel *api.Skel `json:"skel,omitempty"`
}

// Get just for demo
// @postfilter("Boss")
func (skel *Skel) Get(w http.ResponseWriter, r *http.Request) {
	fun := "skel.Skel.Get"

	var getReq *GetReq
	err := json.Unmarshal(skel.ReadBody(r), &getReq)
	if err != nil || !getReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, getReq)
		skel.ReplyFail(w, lib.CodePara)
		return
	}

	skelService := service.NewSkel()
	skelAPI, err := skelService.Get(getReq.ID)
	if err != nil {
		clog.Error("%s get err: %v, req: %v", fun, err, getReq)
		skel.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &GetRsp{
		Skel: skelAPI,
	}
	skel.ReplyOk(w, resp)

	return
}
