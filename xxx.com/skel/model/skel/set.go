package skel

import "xxx.com/skel/mongo"

func (skel *Skel) Set() (err error) {
	db, table := skel.Db(), skel.Table()
	session := mongo.DBS[db]
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	c := sessionCopy.DB(db).C(table)

	err = c.Insert(skel)
	if err != nil {
		return
	}

	return
}
