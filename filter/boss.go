package filter

import (
	"net/http"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/simplejia/lib"

	"github.com/simplejia/clog/api"
)

// Boss 后置过滤器，用于数据上报，比如调用延时，出错等
func Boss(w http.ResponseWriter, r *http.Request, m map[string]interface{}) bool {
	err := m["__E__"]
	path := m["__P__"]
	c := m["__C__"].(lib.IBase)
	bt := m["__T__"].(time.Time)

	trace, _ := c.GetParam(lib.KeyTrace)

	timeout, _ := c.GetParam(lib.KeyTimeoutOccur)
	if timeout != nil {
		timeout = atomic.LoadInt32(timeout.(*int32))
	}

	if err != nil {
		clog.Error("Boss() path: %v, body: %s, err: %v, stack: %s, timeout: %v, trace: %v, elapse: %s", path, c.ReadBody(r), err, debug.Stack(), timeout, trace, time.Since(bt))

		if _, ok := c.GetParam(lib.KeyTimeout); ok {
			muI, _ := c.GetParam(lib.KeyTimeoutMutex)
			mu := muI.(*int32)
			ok := atomic.CompareAndSwapInt32(mu, 0, 1)
			if !ok {
				return true
			}

			doneI, _ := c.GetParam(lib.KeyTimeoutDone)
			done := doneI.(chan struct{})
			close(done)
		}

		w.WriteHeader(http.StatusInternalServerError)
		return true
	}

	resp, ok := c.GetParam(lib.KeyResp)
	if !ok {
		resp = []byte(nil)
	}
	clog.Info("Boss() path: %v, body: %s, resp: %s, timeout: %v, trace: %v, elapse: %s", path, c.ReadBody(r), resp.([]byte), timeout, trace, time.Since(bt))

	return true
}
