package skel

// Add 定义新增操作
func (skel *Skel) Add() (err error) {
	c := skel.GetC()
	defer c.Database.Session.Close()

	err = c.Insert(skel)
	if err != nil {
		return
	}

	return
}
