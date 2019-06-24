package filter

import (
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/simplejia/lib"
	"github.com/simplejia/skel/conf"
)

// Trace 前置过滤器，用于填充trace
func Trace(w http.ResponseWriter, r *http.Request, m map[string]interface{}) (ok bool) {
	c := m["__C__"].(lib.IBase)
	path := m["__P__"].(string)

	trace := &lib.Trace{
		SrvDst:  conf.C.App.Name,
		NameDst: path,
	}

	hTrace := r.Header.Get("h_trace")
	if hTrace == "" {
		trace.ID = bson.NewObjectId().Hex()
	} else {
		trace.Decode(hTrace)
	}

	c.SetParam(lib.KeyTrace, trace)
	return true

}
