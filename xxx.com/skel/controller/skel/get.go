package skel

import (
	"encoding/json"
	"net/http"

	"xxx.com/lib"
	"xxx.com/skel/api"
	"xxx.com/skel/service"

	"github.com/simplejia/clog/api"
)

// GetReq 定义输入
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

// GetResp 定义输出
type GetResp struct {
	Skel *api.Skel `json:"skel,omitempty"`
}

// Get just for demo
// @postfilter("Boss")
func (skel *Skel) Get(w http.ResponseWriter, r *http.Request) {
	fun := "skel.Skel.Get"

	var getReq *GetReq
	if err := json.Unmarshal(skel.ReadBody(r), &getReq); err != nil || !getReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, getReq)
		skel.ReplyFail(w, lib.CodePara)
		return
	}

	skelApi, err := service.NewSkel().Get(getReq.ID)
	if err != nil {
		clog.Error("%s skel.Get err: %v, req: %v", fun, err, getReq)
		skel.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &GetResp{
		Skel: skelApi,
	}
	skel.ReplyOk(w, resp)

	// 进行一些异步处理的工作
	go lib.Updates(skelApi, lib.GET, nil)

	return
}
