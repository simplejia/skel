package skel

import "github.com/simplejia/utils"

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
