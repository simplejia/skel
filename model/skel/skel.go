package skel

import (
	"github.com/globalsign/mgo"
	"github.com/simplejia/skel/mongo"
	"github.com/simplejia/utils"
)

// Skel 定义Skel类型
type Skel struct {
	*utils.Trace
}

func (skel *Skel) WithTrace(trace *utils.Trace) *Skel {
	if skel == nil {
		return nil
	}

	skel.Trace = trace
	return skel
}

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
