package skel

import (
	"encoding/json"

	"github.com/simplejia/skel_api"
	"github.com/simplejia/utils"
)

// Add 封装controller.Add操作
func Add(req *skel_api.SkelAddReq) (resp *skel_api.SkelAddResp, result *utils.Resp, err error) {
	c := &Skel{}
	body, err := utils.TestPost(c.Add, req)
	if err != nil {
		return
	}

	s := &struct {
		utils.Resp
		Data *skel_api.SkelAddResp `json:"data"`
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != utils.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data

	return
}

// Update 封装controller.Update操作
func Update(req *skel_api.SkelUpdateReq) (resp *skel_api.SkelUpdateResp, result *utils.Resp, err error) {
	c := &Skel{}
	body, err := utils.TestPost(c.Update, req)
	if err != nil {
		return
	}

	s := &struct {
		utils.Resp
		Data *skel_api.SkelUpdateResp `json:"data"`
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != utils.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data

	return
}

// Del 封装controller.Del操作
func Del(req *skel_api.SkelDelReq) (resp *skel_api.SkelDelResp, result *utils.Resp, err error) {
	c := &Skel{}
	body, err := utils.TestPost(c.Del, req)
	if err != nil {
		return
	}

	s := &struct {
		utils.Resp
		Data *skel_api.SkelDelResp `json:"data"`
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != utils.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data

	return
}

// Get 封装controller.Get操作
func Get(req *skel_api.SkelGetReq) (resp *skel_api.SkelGetResp, result *utils.Resp, err error) {
	c := &Skel{}
	body, err := utils.TestPost(c.Get, req)
	if err != nil {
		return
	}

	s := &struct {
		utils.Resp
		Data *skel_api.SkelGetResp `json:"data"`
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != utils.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}

// Gets 封装controller.Gets操作
func Gets(req *skel_api.SkelGetsReq) (resp *skel_api.SkelGetsResp, result *utils.Resp, err error) {
	c := &Skel{}
	body, err := utils.TestPost(c.Gets, req)
	if err != nil {
		return
	}

	s := &struct {
		utils.Resp
		Data *skel_api.SkelGetsResp `json:"data"`
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != utils.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}

// PageList 封装controller.PageList操作
func PageList(req *skel_api.SkelPageListReq) (resp *skel_api.SkelPageListResp, result *utils.Resp, err error) {
	c := &Skel{}
	body, err := utils.TestPost(c.PageList, req)
	if err != nil {
		return
	}

	s := &struct {
		utils.Resp
		Data *skel_api.SkelPageListResp `json:"data"`
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != utils.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}

// FlowList 封装controller.FlowList操作
func FlowList(req *skel_api.SkelFlowListReq) (resp *skel_api.SkelFlowListResp, result *utils.Resp, err error) {
	c := &Skel{}
	body, err := utils.TestPost(c.FlowList, req)
	if err != nil {
		return
	}

	s := &struct {
		utils.Resp
		Data *skel_api.SkelFlowListResp `json:"data"`
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != utils.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}

// Upsert 封装controller.Upsert操作
func Upsert(req *skel_api.SkelUpsertReq) (resp *skel_api.SkelUpsertResp, result *utils.Resp, err error) {
	c := &Skel{}
	body, err := utils.TestPost(c.Upsert, req)
	if err != nil {
		return
	}

	s := &struct {
		utils.Resp
		Data *skel_api.SkelUpsertResp `json:"data"`
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != utils.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}
