package api

import "xxx.com/skel/model/skel"

type Skel struct {
	*skel.Skel
}

func NewSkel() *Skel {
	return &Skel{
		Skel: &skel.Skel{},
	}
}
