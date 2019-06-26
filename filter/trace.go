package filter

import (
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/simplejia/utils"
	"github.com/simplejia/skel/conf"
)

// Trace 前置过滤器，用于填充trace
func Trace(w http.ResponseWriter, r *http.Request, m map[string]interface{}) (ok bool) {
	c := m["__C__"].(utils.IBase)
	path := m["__P__"].(string)

	trace := &utils.Trace{
		SrvDst:  conf.C.App.Name,
		NameDst: path,
	}

	hTrace := r.Header.Get("h_trace")
	if hTrace == "" {
		trace.ID = bson.NewObjectId().Hex()
	} else {
		trace.Decode(hTrace)
	}

	c.SetParam(utils.KeyTrace, trace)
	return true

}
