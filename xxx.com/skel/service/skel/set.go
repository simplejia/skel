package skel

import "xxx.com/skel/model"

func (skel *Skel) Set(id int64) (err error) {
	fun := "skel.Skel.Set"
	_ = fun

	skelModel := model.NewSkel()
	skelModel.ID = id

	err = skelModel.Set()
	if err != nil {
		return
	}

	return
}
