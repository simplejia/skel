package filter

import (
	"net/http"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/simplejia/clog/api"
	"github.com/simplejia/utils"
)

// Boss 后置过滤器，用于数据上报，比如调用延时，出错等
func Boss(w http.ResponseWriter, r *http.Request, m map[string]interface{}) bool {
	err := m["__E__"]
	path := m["__P__"]
	c := m["__C__"].(utils.IBase)
	bt := m["__T__"].(time.Time)

	trace, _ := c.GetParam(utils.KeyTrace)

	timeout, _ := c.GetParam(utils.KeyTimeoutOccur)
	if timeout != nil {
		timeout = atomic.LoadInt32(timeout.(*int32))
	}

	if err != nil {
		clog.Error("Boss() path: %v, body: %s, elapse: %s, timeout: %v, trace: %v, err: %v, stack: %s", path, c.ReadBody(r), time.Since(bt), timeout, trace, err, debug.Stack())

		if _, ok := c.GetParam(utils.KeyTimeout); ok {
			muI, _ := c.GetParam(utils.KeyTimeoutMutex)
			mu := muI.(*int32)
			ok := atomic.CompareAndSwapInt32(mu, 0, 1)
			if !ok {
				return true
			}

			doneI, _ := c.GetParam(utils.KeyTimeoutDone)
			done := doneI.(chan struct{})
			close(done)
		}

		w.WriteHeader(http.StatusInternalServerError)
		return true
	}

	resp, ok := c.GetParam(utils.KeyResp)
	if !ok {
		resp = []byte(nil)
	}
	clog.Info("Boss() path: %v, body: %s, elapse: %s, timeout: %v, trace: %v, resp: %s", path, c.ReadBody(r), time.Since(bt), timeout, trace, resp.([]byte))

	return true
}
