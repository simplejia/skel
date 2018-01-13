package skel

import (
	"encoding/json"
	"net/http"

	"xxx.com/lib"
	"xxx.com/skel/api"
	"xxx.com/skel/service"

	"github.com/simplejia/clog"
)

type DelReq struct {
	ID int64 `json:"id"`
}

// Regular 用于参数校验
func (delReq *DelReq) Regular() (ok bool) {
	if delReq == nil {
		return
	}

	if delReq.ID <= 0 {
		return
	}

	ok = true
	return
}

type DelRsp struct {
}

// Del just for demo
// @postfilter("Boss")
func (skel *Skel) Del(w http.ResponseWriter, r *http.Request) {
	fun := "skel.Skel.Del"

	var delReq *DelReq
	err := json.Unmarshal(skel.ReadBody(r), &delReq)
	if err != nil || !delReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, delReq)
		skel.ReplyFail(w, lib.CodePara)
		return
	}

	skelService := service.NewSkel()
	err = skelService.Del(delReq.ID)
	if err != nil {
		clog.Error("%s del err: %v, req: %v", fun, err, delReq)
		skel.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &DelRsp{}
	skel.ReplyOk(w, resp)

	skelAPI := api.NewSkel()
	skelAPI.ID = delReq.ID
	// 进行一些异步处理的工作
	go lib.Updates(skelAPI, lib.DELETE, nil)

	return
}
