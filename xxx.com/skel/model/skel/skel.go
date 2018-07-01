/*
Package skel just for demo
*/
package skel

import (
	mgo "gopkg.in/mgo.v2"
	"xxx.com/skel/api"
	"xxx.com/skel/mongo"
)

// Skel 定义db对应的类型
type Skel api.Skel

// Db 返回db name
func (skel *Skel) Db() (db string) {
	return "skel"
}

// Table 返回table name
func (skel *Skel) Table() (table string) {
	return "skel"
}

// GetC 返回db col
func (skel *Skel) GetC() (c *mgo.Collection) {
	db, table := skel.Db(), skel.Table()
	session := mongo.DBS[db]
	sessionCopy := session.Copy()
	c = sessionCopy.DB(db).C(table)
	return
}
