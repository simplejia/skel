// 此文件封装仅用于单元测试公共函数

package skel

import (
	"encoding/json"
	"fmt"

	"xxx.com/lib"
	"xxx.com/skel/api"
	"xxx.com/skel/test"
)

// Set 封装controller.Set操作
func Set(id int64) (err error) {
	c := &Skel{}
	p := &SetReq{
		ID: id,
	}
	body, err := lib.TestPost(c.Set, p)
	if err != nil {
		return
	}

	resp := &struct {
		lib.Resp
		Data SetRsp `json:"data"`
	}{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return
	}
	if resp.Ret != lib.CodeOk {
		err = fmt.Errorf("resp ret invalid: %v, body: %s", resp, body)
		return
	}

	test.Sleep()
	return
}

// Del 封装controller.Del操作
func Del(id int64) (err error) {
	c := &Skel{}
	p := &SetReq{
		ID: id,
	}
	body, err := lib.TestPost(c.Del, p)
	if err != nil {
		return
	}

	resp := &struct {
		lib.Resp
		Data SetRsp `json:"data"`
	}{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return
	}
	if resp.Ret != lib.CodeOk {
		err = fmt.Errorf("resp ret invalid: %v, body: %s", resp, body)
		return
	}

	test.Sleep()
	return
}

// Get 封装controller.Get操作
func Get(id int64) (skel *api.Skel, err error) {
	c := &Skel{}
	p := &GetReq{
		ID: id,
	}
	body, err := lib.TestPost(c.Get, p)
	if err != nil {
		return
	}

	resp := &struct {
		lib.Resp
		Data GetRsp `json:"data"`
	}{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return
	}
	if resp.Ret != lib.CodeOk {
		err = fmt.Errorf("resp ret invalid: %v, body: %s", resp, body)
		return
	}

	skel = resp.Data.Skel
	return
}

// GetSkel 返回一个全新的skel对象
func GetSkel() (skel *api.Skel) {
	id := test.GetID()
	for idTemp := id; ; idTemp++ {
		if skelTemp, err := Get(idTemp); err != nil {
			panic(err)
		} else if skelTemp != nil {
			if err := Del(idTemp); err != nil {
				panic(err)
			}
		} else {
			break
		}
	}
	if err := Set(id); err != nil {
		panic(err)
	}
	skel, err := Get(id)
	if err != nil {
		panic(err)
	} else if skel == nil {
		panic("GetSkel fail")
	}

	return
}
