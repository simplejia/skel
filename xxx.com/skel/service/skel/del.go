package skel

import "xxx.com/skel/model"

func (skel *Skel) Del(id int64) (err error) {
	fun := "skel.Skel.Del"
	_ = fun

	skelModel := model.NewSkel()
	skelModel.ID = id

	err = skelModel.Del()
	if err != nil {
		return
	}

	return
}
