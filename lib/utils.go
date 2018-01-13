package lib

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"

	"github.com/simplejia/namecli/api"
	"github.com/simplejia/utils"
)

func TestPost(h http.HandlerFunc, params interface{}) (body []byte, err error) {
	v, err := json.Marshal(params)
	if err != nil {
		return
	}
	r, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(v))
	if err != nil {
		return
	}
	w := httptest.NewRecorder()
	h(w, r)
	body = w.Body.Bytes()
	if g, e := w.Code, http.StatusOK; g != e {
		err = fmt.Errorf("http resp status not ok: %s", http.StatusText(g))
		return
	}
	return
}

func SearchInt64s(a []int64, x int64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

func Int64s(a []int64) {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
}

func ZipInt64s(a []int64) (result []byte, err error) {
	if len(a) == 0 {
		return
	}

	var b bytes.Buffer
	zw, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return
	}

	err = json.NewEncoder(zw).Encode(a)
	if err != nil {
		return
	}
	zw.Close()

	result = b.Bytes()
	return
}

func UnzipInt64s(a []byte) (result []int64, err error) {
	if len(a) == 0 {
		return
	}

	zr := flate.NewReader(bytes.NewReader(a))
	err = json.NewDecoder(zr).Decode(&result)
	if err != nil {
		return
	}
	return
}

func ZipBytes(a []byte) (result []byte, err error) {
	if len(a) == 0 {
		return
	}

	var b bytes.Buffer
	zw, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return
	}

	zw.Write(a)
	zw.Close()

	result = b.Bytes()
	return
}

func UnzipBytes(a []byte) (result []byte, err error) {
	if len(a) == 0 {
		return
	}

	zr := flate.NewReader(bytes.NewReader(a))
	bs, err := ioutil.ReadAll(zr)
	if err != nil {
		return
	}

	result = bs
	return
}

func NameWrap(name string) (addr string, err error) {
	ip := strings.Replace(strings.Replace(name, ".", "", 4), ":", "", 1)
	if _, err := strconv.Atoi(ip); err == nil {
		return name, nil
	}
	addr, err = api.Name(name)
	return
}

func PostProxy(name, path string, req []byte) (rsp []byte, err error) {
	addr, err := NameWrap(name)
	if err != nil {
		return
	}
	url := fmt.Sprintf("http://%s/%s", addr, strings.TrimPrefix(path, "/"))

	gpp := &utils.GPP{
		Uri:    url,
		Params: req,
	}
	rsp, err = utils.Post(gpp)
	if err != nil {
		return
	}

	return
}
