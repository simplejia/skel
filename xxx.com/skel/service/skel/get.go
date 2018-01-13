package skel

import (
	"xxx.com/skel/api"
	"xxx.com/skel/model"
)

func (skel *Skel) Get(id int64) (skelAPI *api.Skel, err error) {
	fun := "skel.Skel.Get"
	_ = fun

	skelModel := model.NewSkel()
	skelModel.ID = id
	skelModelRet, err := skelModel.Get()
	if err != nil {
		return
	}

	if skelModelRet == nil {
		return
	}

	skelAPI = &api.Skel{
		Skel: skelModelRet,
	}
	return
}
