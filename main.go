// skel，just for demo
// author: simplejia
// date: 2017/12/8

//go:generate wsp -s -d

package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/simplejia/clog/api"
	"github.com/simplejia/lc"
	"github.com/simplejia/namecli/api"
	"github.com/simplejia/skel/conf"
	"github.com/simplejia/utils"
)

func init() {
	lc.Init(1e5)

	clog.AddrFunc = func() (string, error) {
		return api.Name(conf.C.Addrs.Clog)
	}
	clog.Init(conf.C.App.Name, "", conf.C.Clog.Level, conf.C.Clog.Mode)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clog.Error("%s is not found", r.RequestURI)
		http.NotFound(w, r)
	})
}

func main() {
	fun := "main"
	clog.Info("%s rlimit nofile: %s", fun, utils.RlimitNofile())

	addr := fmt.Sprintf(":%d", conf.C.App.Port)

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		done := make(chan struct{})
		ctx, cancelCtx := context.WithCancel(context.WithValue(r.Context(), utils.CtxDone, done))
		r = r.WithContext(ctx)

		go func() {
			defer cancelCtx()
			http.DefaultServeMux.ServeHTTP(w, r)
		}()

		select {
		case <-done:
		case <-ctx.Done():
		}
	})
	if err := utils.ListenAndServe(addr, rootHandler); err != nil {
		clog.Error("%s err: %v, addr: %v", fun, err, addr)
	}

	utils.AsyncTaskShutdown(time.Second * 3)
}
