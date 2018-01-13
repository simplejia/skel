/*
Package skel just for demo
*/
package skel

type Skel struct {
	ID int64 `json:"id" bson:"_id"`
}

func (skel *Skel) Db() (db string) {
	return "skel"
}

func (skel *Skel) Table() (table string) {
	return "skel"
}

func (skel *Skel) Regular() (ok bool) {
	if skel == nil {
		return
	}

	ok = true
	return
}
