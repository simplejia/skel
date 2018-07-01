package skel

import (
	"xxx.com/skel/api"
	"xxx.com/skel/model"
)

// Get 定义获取操作
func (skel *Skel) Get(id int64) (skelApi *api.Skel, err error) {
	skelModel := model.NewSkel()
	skelModel.ID = id
	if skelModel, err = skelModel.Get(); err != nil {
		return
	}

	if skelModel == nil {
		return
	}

	skelApi = (*api.Skel)(skelModel)

	return
}
