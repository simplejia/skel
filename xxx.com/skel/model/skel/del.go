package skel

import (
	"github.com/globalsign/mgo"
)

// Del 定义删除操作
func (skel *Skel) Del() (err error) {
	c := skel.GetC()
	defer c.Database.Session.Close()

	err = c.Remove(skel)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}
		err = nil
		return
	}

	return
}
