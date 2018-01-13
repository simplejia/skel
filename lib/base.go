package lib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/simplejia/clog"
)

const (
	// KeyBody request body内容会存在此key对应params参数里
	KeyBody = "_body_"
)

// IBase 所有Controller必须实现此接口
type IBase interface {
	SetParam(string, interface{})
	GetParam(string) (interface{}, bool)
	ReadBody(*http.Request) []byte
}

type Base struct {
	params map[string]interface{}
}

func (base *Base) SetParam(key string, value interface{}) {
	if base.params == nil {
		base.params = make(map[string]interface{})
	}
	base.params[key] = value
}

func (base *Base) GetParam(key string) (value interface{}, ok bool) {
	value, ok = base.params[key]
	return
}

func (base *Base) ReadBody(r *http.Request) (body []byte) {
	key := KeyBody
	value, ok := base.GetParam(key)
	if ok {
		body = value.([]byte)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		clog.Error("lib.Base.ReadBody err: %v", err)
		return
	}
	base.SetParam(key, body)
	return
}

func (base *Base) ReplyOk(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(&Resp{
		Ret:  CodeOk,
		Data: data,
	})
}

func (base *Base) ReplyFail(w http.ResponseWriter, code Code) {
	json.NewEncoder(w).Encode(&Resp{
		Ret: code,
		Msg: CodeMap[Code(code)],
	})
}

func (base *Base) ReplyFailWithDetail(w http.ResponseWriter, code Code, detail string) {
	json.NewEncoder(w).Encode(&Resp{
		Ret:    code,
		Detail: detail,
		Msg:    CodeMap[Code(code)],
	})
}
