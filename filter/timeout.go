package filter

import (
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/simplejia/lib"
)

// Timeout 前置过滤器，用于接口超时控制
func Timeout(w http.ResponseWriter, r *http.Request, m map[string]interface{}) (ok bool) {
	c := m["__C__"].(lib.IBase)

	durStr, _ := m["dur"].(string)
	if durStr == "" {
		return true
	}

	dur, err := time.ParseDuration(durStr)
	if err != nil {
		panic(err)
	}

	if dur == 0 {
		return true
	}

	c.SetParam(lib.KeyTimeout, nil)

	i := new(int32)
	c.SetParam(lib.KeyTimeoutMutex, i)

	j := new(int32)
	c.SetParam(lib.KeyTimeoutOccur, j)

	done := make(chan struct{})
	c.SetParam(lib.KeyTimeoutDone, done)

	go func() {
		timer := time.NewTimer(dur)
		defer timer.Stop()

		select {
		case <-timer.C:
			ok := atomic.CompareAndSwapInt32(i, 0, 1)
			if !ok {
				break
			}

			atomic.StoreInt32(j, 1)
			w.WriteHeader(http.StatusServiceUnavailable)
			io.WriteString(w, "Timeout error")

			if ctxDone, ok := r.Context().Value(lib.CtxDone).(chan struct{}); ok {
				close(ctxDone)
			}
		case <-done:
			// nothing
		}
	}()

	return true
}
