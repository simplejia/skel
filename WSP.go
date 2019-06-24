// Code generated by wsp, DO NOT EDIT.

package main

import "net/http"
import "time"
import "github.com/simplejia/skel/controller/skel"
import "github.com/simplejia/skel/filter"

func init() {
	http.HandleFunc("/skel/add", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(skel.Skel)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/add"}); !ok {
				return
			}
		}()
		if ok := filter.Timeout(w, r, map[string]interface{}{"dur": "1s", "__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/add"}); !ok {
			return
		}
		if ok := filter.Trace(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/add"}); !ok {
			return
		}
		c.Add(w, r)
	})

	http.HandleFunc("/skel/del", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(skel.Skel)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/del"}); !ok {
				return
			}
		}()
		if ok := filter.Timeout(w, r, map[string]interface{}{"dur": "1s", "__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/del"}); !ok {
			return
		}
		if ok := filter.Trace(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/del"}); !ok {
			return
		}
		c.Del(w, r)
	})

	http.HandleFunc("/skel/flow_list", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(skel.Skel)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/flow_list"}); !ok {
				return
			}
		}()
		if ok := filter.Timeout(w, r, map[string]interface{}{"dur": "1s", "__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/flow_list"}); !ok {
			return
		}
		if ok := filter.Trace(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/flow_list"}); !ok {
			return
		}
		c.FlowList(w, r)
	})

	http.HandleFunc("/skel/get", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(skel.Skel)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/get"}); !ok {
				return
			}
		}()
		if ok := filter.Timeout(w, r, map[string]interface{}{"dur": "1s", "__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/get"}); !ok {
			return
		}
		if ok := filter.Trace(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/get"}); !ok {
			return
		}
		c.Get(w, r)
	})

	http.HandleFunc("/skel/gets", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(skel.Skel)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/gets"}); !ok {
				return
			}
		}()
		if ok := filter.Timeout(w, r, map[string]interface{}{"dur": "1s", "__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/gets"}); !ok {
			return
		}
		if ok := filter.Trace(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/gets"}); !ok {
			return
		}
		c.Gets(w, r)
	})

	http.HandleFunc("/skel/page_list", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(skel.Skel)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/page_list"}); !ok {
				return
			}
		}()
		if ok := filter.Timeout(w, r, map[string]interface{}{"dur": "1s", "__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/page_list"}); !ok {
			return
		}
		if ok := filter.Trace(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/page_list"}); !ok {
			return
		}
		c.PageList(w, r)
	})

	http.HandleFunc("/skel/update", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(skel.Skel)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/update"}); !ok {
				return
			}
		}()
		if ok := filter.Timeout(w, r, map[string]interface{}{"dur": "1s", "__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/update"}); !ok {
			return
		}
		if ok := filter.Trace(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/update"}); !ok {
			return
		}
		c.Update(w, r)
	})

	http.HandleFunc("/skel/upsert", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_ = t
		var e interface{}
		c := new(skel.Skel)
		defer func() {
			e = recover()
			if ok := filter.Boss(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/upsert"}); !ok {
				return
			}
		}()
		if ok := filter.Timeout(w, r, map[string]interface{}{"dur": "1s", "__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/upsert"}); !ok {
			return
		}
		if ok := filter.Trace(w, r, map[string]interface{}{"__T__": t, "__C__": c, "__E__": e, "__P__": "/skel/upsert"}); !ok {
			return
		}
		c.Upsert(w, r)
	})

}