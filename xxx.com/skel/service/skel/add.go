package skel

import (
	"xxx.com/skel/api"
	"xxx.com/skel/model"
)

// Add 定义新增操作
func (skel *Skel) Add(skelApi *api.Skel) (err error) {
	skelModel := model.NewSkel()
	skelModel.ID = skelApi.ID
	if err = skelModel.Add(); err != nil {
		return
	}

	return
}
