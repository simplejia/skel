package skel

import (
	"encoding/json"
	"net/http"

	"xxx.com/lib"
	"xxx.com/skel/service"

	"github.com/simplejia/clog"
)

// DelReq 定义输入
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

// DelResp 定义输出
type DelResp struct {
}

// Del just for demo
// @postfilter("Boss")
func (skel *Skel) Del(w http.ResponseWriter, r *http.Request) {
	fun := "skel.Skel.Del"

	var delReq *DelReq
	if err := json.Unmarshal(skel.ReadBody(r), &delReq); err != nil || !delReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, delReq)
		skel.ReplyFail(w, lib.CodePara)
		return
	}

	skelApi, err := service.NewSkel().Get(delReq.ID)
	if err != nil {
		clog.Error("%s skel.Get err: %v, req: %v", fun, err, delReq)
		skel.ReplyFail(w, lib.CodeSrv)
		return
	}

	if skelApi == nil {
		detail := "skel not found"
		skel.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if err := service.NewSkel().Del(delReq.ID); err != nil {
		clog.Error("%s skel.Del err: %v, req: %v", fun, err, delReq)
		skel.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &DelResp{}
	skel.ReplyOk(w, resp)

	// 进行一些异步处理的工作
	go lib.Updates(skelApi, lib.DELETE, nil)

	return
}
