package skel

import "github.com/simplejia/lib"

// Skel 定义Skel类型
type Skel struct {
	*lib.Trace
}

func (skel *Skel) WithTrace(trace *lib.Trace) *Skel {
	if skel == nil {
		return nil
	}

	skel.Trace = trace
	return skel
}
