package skel

import (
	"gopkg.in/mgo.v2"
)

// Get 定义获取操作
func (skel *Skel) Get() (skelRet *Skel, err error) {
	c := skel.GetC()
	defer c.Database.Session.Close()

	err = c.Find(skel).One(&skelRet)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}
		err = nil
		return
	}

	return
}
