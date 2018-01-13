package skel

import (
	"gopkg.in/mgo.v2"
	"xxx.com/skel/mongo"
)

func (skel *Skel) Del() (err error) {
	db, table := skel.Db(), skel.Table()
	session := mongo.DBS[db]
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	c := sessionCopy.DB(db).C(table)

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
