package skel

import (
	"encoding/json"
	"net/http"

	"xxx.com/lib"
	"xxx.com/skel/api"
	"xxx.com/skel/service"

	"github.com/simplejia/clog"
)

type SetReq struct {
	ID int64 `json:"id"`
}

// Regular 用于参数校验
func (setReq *SetReq) Regular() (ok bool) {
	if setReq == nil {
		return
	}

	if setReq.ID <= 0 {
		return
	}

	ok = true
	return
}

type SetRsp struct {
}

// Set just for demo
// @postfilter("Boss")
func (skel *Skel) Set(w http.ResponseWriter, r *http.Request) {
	fun := "skel.Skel.Set"

	var setReq *SetReq
	err := json.Unmarshal(skel.ReadBody(r), &setReq)
	if err != nil || !setReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, setReq)
		skel.ReplyFail(w, lib.CodePara)
		return
	}

	skelService := service.NewSkel()
	err = skelService.Set(setReq.ID)
	if err != nil {
		clog.Error("%s set err: %v, req: %v", fun, err, setReq)
		skel.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &SetRsp{}
	skel.ReplyOk(w, resp)

	skelAPI := api.NewSkel()
	skelAPI.ID = setReq.ID
	// 进行一些异步处理的工作
	go lib.Updates(skelAPI, lib.ADD, nil)

	return
}
