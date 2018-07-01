package skel

import "xxx.com/skel/model"

// Del 定义删除操作
func (skel *Skel) Del(id int64) (err error) {
	skelModel := model.NewSkel()
	skelModel.ID = id

	if err = skelModel.Del(); err != nil {
		return
	}

	return
}
