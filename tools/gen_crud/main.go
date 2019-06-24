package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"
)

var (
	rootPkg        string
	name           string
	timeout        time.Duration
	connectTimeout time.Duration
	idType         string
	needMultiTable bool
	dbNum          int
	tblNum         int
	keys           string // ex: pid,int64,rid,int64
)

var tplOuterAPI = `package {{.proj}}_api

// {{.camel}} just for skel
type {{.camel}} struct {
	ID {{.id_type}}` + " `" + `json:"id" bson:"_id"` + "`" + `
	{{- range $elem := (split_keys .keys)}}
	{{camel (index $elem 0)}} {{index $elem 1}}` + " `" + `json:"{{snake (index $elem 0)}}" bson:"{{snake (index $elem 0)}}"` + "`" + `
	{{- end}}
}
`

var tplOuterAPIGet = `package {{.proj}}_api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/simplejia/utils"
	"{{.pkg}}/lib"
)

// {{.camel}}GetReq 定义输入
type {{.camel}}GetReq struct {
	{{- if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	{{- range $elem := .}}
	{{camel (index $elem 0)}} {{index $elem 1}}` + " `" + `json:"{{snake (index $elem 0)}}"` + "`" + `
	{{- end}}
	{{- end}}
	{{- else}}
	ID {{.id_type}}` + " `" + `json:"id"` + "`" + `
	{{- end}}
}

// Regular 用于参数校验
func ({{.lower}}GetReq *{{.camel}}GetReq) Regular() (ok bool) {
	if {{.lower}}GetReq == nil {
		return
	}

	ok = true
	return
}

// {{.camel}}GetResp 定义输出
type {{.camel}}GetResp {{.camel}}

func {{.camel}}Get(name string, req *{{.camel}}GetReq, trace *lib.Trace) (resp *{{.camel}}GetResp, result *lib.Resp, err error) {
	addr, err := lib.NameWrap(name)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("http://%s/%s", addr, "{{.snake}}/get")
	gpp := &utils.GPP{
		Uri:            uri,
		{{- with .connect_timeout}}
		ConnectTimeout: {{.}},
		{{- end}}
		{{- with .timeout}}
		Timeout: {{.}},
		{{- end}}
		Params:         req,
		Headers: map[string]string{
			"h_trace": trace.Encode(),
		},
	}

	body, err := utils.Post(gpp)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.camel}}GetResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}

	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}
`

var tplOuterAPIGets = `package {{.proj}}_api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/simplejia/utils"
	"{{.pkg}}/lib"
)

// {{.camel}}GetsReq 定义输入
type {{.camel}}GetsReq struct {
	{{- if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	{{camel (index . 0 0)}} {{index . 0 1}}` + " `" + `json:"{{snake (index . 0 0)}}"` + "`" + `
	{{camel (index . 1 0)}}s []{{index . 1 1}}` + " `" + `json:"{{snake (index . 1 0)}}s"` + "`" + `
	{{- end}}
	{{- else}}
	IDS []{{.id_type}}` + " `" + `json:"ids"` + "`" + `
	{{- end}}
}

// Regular 用于参数校验
func ({{.lower}}GetsReq *{{.camel}}GetsReq) Regular() (ok bool) {
	if {{.lower}}GetsReq == nil {
		return
	}

	ok = true
	return
}

{{- define "map_key"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- index . 1 1}}
{{- end}}
{{- else -}}
{{.id_type}}
{{- end}}
{{- end}}

// {{.camel}}GetsResp 定义输出
type {{.camel}}GetsResp map[{{template "map_key" dict "keys" .keys "id_type" .id_type}}]{{if gt (len (split_keys .keys)) 1}}[]{{end}}*{{.camel}}

func {{.camel}}Gets(name string, req *{{.camel}}GetsReq, trace *lib.Trace) (resp *{{.camel}}GetsResp, result *lib.Resp, err error) {
	addr, err := lib.NameWrap(name)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("http://%s/%s", addr, "{{.snake}}/gets")
	gpp := &utils.GPP{
		Uri:            uri,
		{{- with .connect_timeout}}
		ConnectTimeout: {{.}},
		{{- end}}
		{{- with .timeout}}
		Timeout: {{.}},
		{{- end}}
		Params:         req,
		Headers: map[string]string{
			"h_trace": trace.Encode(),
		},
	}

	body, err := utils.Post(gpp)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.camel}}GetsResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}

	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}
`

var tplOuterAPIAdd = `package {{.proj}}_api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/simplejia/utils"
	"{{.pkg}}/lib"
)

// {{.camel}}AddReq 定义输入
type {{.camel}}AddReq {{.camel}}

// Regular 用于参数校验
func ({{.lower}}AddReq *{{.camel}}AddReq) Regular() (ok bool) {
	if {{.lower}}AddReq == nil {
		return
	}

	ok = true
	return
}

// {{.camel}}AddResp 定义输出
type {{.camel}}AddResp {{.camel}}

func {{.camel}}Add(name string, req *{{.camel}}AddReq, trace *lib.Trace) (resp *{{.camel}}AddResp, result *lib.Resp, err error) {
	addr, err := lib.NameWrap(name)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("http://%s/%s", addr, "{{.snake}}/add")
	gpp := &utils.GPP{
		Uri:            uri,
		{{- with .connect_timeout}}
		ConnectTimeout: {{.}},
		{{- end}}
		{{- with .timeout}}
		Timeout: {{.}},
		{{- end}}
		Params:         req,
		Headers: map[string]string{
			"h_trace": trace.Encode(),
		},
	}

	body, err := utils.Post(gpp)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.camel}}AddResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}

	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}
`

var tplOuterAPIDel = `package {{.proj}}_api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/simplejia/utils"
	"{{.pkg}}/lib"
)

// {{.camel}}DelReq 定义输入
type {{.camel}}DelReq struct {
	{{- if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	{{- range $elem := .}}
	{{camel (index $elem 0)}} {{index $elem 1}}` + " `" + `json:"{{snake (index $elem 0)}}"` + "`" + `
	{{- end}}
	{{- end}}
	{{- else}}
	ID {{.id_type}}` + " `" + `json:"id"` + "`" + `
	{{- end}}
}

// Regular 用于参数校验
func ({{.lower}}DelReq *{{.camel}}DelReq) Regular() (ok bool) {
	if {{.lower}}DelReq == nil {
		return
	}

	ok = true
	return
}

// {{.camel}}DelResp 定义输出
type {{.camel}}DelResp struct {
}

func {{.camel}}Del(name string, req *{{.camel}}DelReq, trace *lib.Trace) (resp *{{.camel}}DelResp, result *lib.Resp, err error) {
	addr, err := lib.NameWrap(name)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("http://%s/%s", addr, "{{.snake}}/del")
	gpp := &utils.GPP{
		Uri:            uri,
		{{- with .connect_timeout}}
		ConnectTimeout: {{.}},
		{{- end}}
		{{- with .timeout}}
		Timeout: {{.}},
		{{- end}}
		Params:         req,
		Headers: map[string]string{
			"h_trace": trace.Encode(),
		},
	}

	body, err := utils.Post(gpp)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.camel}}DelResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}

	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}
`

var tplOuterAPIUpdate = `package {{.proj}}_api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/simplejia/utils"
	"{{.pkg}}/lib"
)

// {{.camel}}UpdateReq 定义输入
type {{.camel}}UpdateReq {{.camel}}

// Regular 用于参数校验
func ({{.lower}}UpdateReq *{{.camel}}UpdateReq) Regular() (ok bool) {
	if {{.lower}}UpdateReq == nil {
		return
	}

	ok = true
	return
}

// {{.camel}}UpdateResp 定义输出
type {{.camel}}UpdateResp {{.camel}}

func {{.camel}}Update(name string, req *{{.camel}}UpdateReq, trace *lib.Trace) (resp *{{.camel}}UpdateResp, result *lib.Resp, err error) {
	addr, err := lib.NameWrap(name)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("http://%s/%s", addr, "{{.snake}}/update")
	gpp := &utils.GPP{
		Uri:            uri,
		{{- with .connect_timeout}}
		ConnectTimeout: {{.}},
		{{- end}}
		{{- with .timeout}}
		Timeout: {{.}},
		{{- end}}
		Params:         req,
		Headers: map[string]string{
			"h_trace": trace.Encode(),
		},
	}

	body, err := utils.Post(gpp)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.camel}}UpdateResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}

	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}
`

var tplOuterAPIPageList = `package {{.proj}}_api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/simplejia/utils"
	"{{.pkg}}/lib"
)

// {{.camel}}PageListReq 定义输入
type {{.camel}}PageListReq struct {
	{{- if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	{{camel (index . 0 0)}} {{index . 0 1}}` + " `" + `json:"{{snake (index . 0 0)}}"` + "`" + `
	{{- end}}
	{{- end}}
	Offset int` + " `" + `json:"offset,omitempty"` + "`" + `
	Limit int` + " `" + `json:"limit,omitempty"` + "`" + `
}

// Regular 用于参数校验
func ({{.lower}}PageListReq *{{.camel}}PageListReq) Regular() (ok bool) {
	if {{.lower}}PageListReq == nil {
		return
	}

	if {{.lower}}PageListReq.Limit <= 0 {
		{{.lower}}PageListReq.Limit = 20
	}

	ok = true
	return
}

// {{.camel}}PageListResp 定义输出
type {{.camel}}PageListResp struct {
	List []*{{.camel}}` + " `" + `json:"list,omitempty"` + "`" + `
	Offset int` + " `" + `json:"offset,omitempty"` + "`" + `
	Limit int` + " `" + `json:"limit,omitempty"` + "`" + `
	More bool` + " `" + `json:"more,omitempty"` + "`" + `
	Total int` + " `" + `json:"total,omitempty"` + "`" + `
}

func {{.camel}}PageList(name string, req *{{.camel}}PageListReq, trace *lib.Trace) (resp *{{.camel}}PageListResp, result *lib.Resp, err error) {
	addr, err := lib.NameWrap(name)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("http://%s/%s", addr, "{{.snake}}/page_list")
	gpp := &utils.GPP{
		Uri:            uri,
		{{- with .connect_timeout}}
		ConnectTimeout: {{.}},
		{{- end}}
		{{- with .timeout}}
		Timeout: {{.}},
		{{- end}}
		Params:         req,
		Headers: map[string]string{
			"h_trace": trace.Encode(),
		},
	}

	body, err := utils.Post(gpp)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.camel}}PageListResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}

	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}
`

var tplOuterAPIFlowList = `package {{.proj}}_api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/simplejia/utils"
	"{{.pkg}}/lib"
)

// {{.camel}}FlowListReq 定义输入
type {{.camel}}FlowListReq struct {
	{{- if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	{{camel (index . 0 0)}} {{index . 0 1}}` + " `" + `json:"{{snake (index . 0 0)}}"` + "`" + `
	{{- end}}
	{{- end}}
	LastID string` + " `" + `json:"last_id,omitempty"` + "`" + `
	Limit int` + " `" + `json:"limit,omitempty"` + "`" + `
}

// Regular 用于参数校验
func ({{.lower}}FlowListReq *{{.camel}}FlowListReq) Regular() (ok bool) {
	if {{.lower}}FlowListReq == nil {
		return
	}

	if {{.lower}}FlowListReq.Limit <= 0 {
		{{.lower}}FlowListReq.Limit = 20
	}

	ok = true
	return
}

// {{.camel}}FlowListResp 定义输出
type {{.camel}}FlowListResp struct {
	List []*{{.camel}}` + " `" + `json:"list,omitempty"` + "`" + `
	LastID string` + " `" + `json:"last_id,omitempty"` + "`" + `
	Limit int` + " `" + `json:"limit,omitempty"` + "`" + `
	More bool` + " `" + `json:"more,omitempty"` + "`" + `
	Total int` + " `" + `json:"total,omitempty"` + "`" + `
}

func {{.camel}}FlowList(name string, req *{{.camel}}FlowListReq, trace *lib.Trace) (resp *{{.camel}}FlowListResp, result *lib.Resp, err error) {
	addr, err := lib.NameWrap(name)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("http://%s/%s", addr, "{{.snake}}/flow_list")
	gpp := &utils.GPP{
		Uri:            uri,
		{{- with .connect_timeout}}
		ConnectTimeout: {{.}},
		{{- end}}
		{{- with .timeout}}
		Timeout: {{.}},
		{{- end}}
		Params:         req,
		Headers: map[string]string{
			"h_trace": trace.Encode(),
		},
	}

	body, err := utils.Post(gpp)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.camel}}FlowListResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}

	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}
`

var tplOuterAPIUpsert = `package {{.proj}}_api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/simplejia/utils"
	"{{.pkg}}/lib"
)

// {{.camel}}UpsertReq 定义输入
type {{.camel}}UpsertReq {{.camel}}

// Regular 用于参数校验
func ({{.lower}}UpsertReq *{{.camel}}UpsertReq) Regular() (ok bool) {
	if {{.lower}}UpsertReq == nil {
		return
	}

	ok = true
	return
}

// {{.camel}}UpsertResp 定义输出
type {{.camel}}UpsertResp {{.camel}}

func {{.camel}}Upsert(name string, req *{{.camel}}UpsertReq, trace *lib.Trace) (resp *{{.camel}}UpsertResp, result *lib.Resp, err error) {
	addr, err := lib.NameWrap(name)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("http://%s/%s", addr, "{{.snake}}/upsert")
	gpp := &utils.GPP{
		Uri:            uri,
		{{- with .connect_timeout}}
		ConnectTimeout: {{.}},
		{{- end}}
		{{- with .timeout}}
		Timeout: {{.}},
		{{- end}}
		Params:         req,
		Headers: map[string]string{
			"h_trace": trace.Encode(),
		},
	}

	body, err := utils.Post(gpp)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.camel}}UpsertResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}

	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}
`

var tplModel = `package {{.snake}}

import (
	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/conf"
	"{{.pkg}}/{{.proj}}/mongo"
)

{{- if .need_multi_table}}
var (
	// DbNum 分库分表，代表库总数
	DbNum = {{.db_num}}
	// TableNum 分库分表，代表表总数
	TableNum = {{.table_num}}
)

func init() {
	// 方便本地测试
	if conf.Env == lib.DEV {
		DbNum = 1
		TableNum = 1
	}
}
{{- end}}

// {{.camel}} 定义{{.camel}}类型
type {{.camel}} struct {
	*lib.Trace
}

func ({{.lower}} *{{.camel}}) WithTrace(trace *lib.Trace) *{{.camel}} {
	if {{.lower}} == nil {
		return nil
	}

	{{.lower}}.Trace = trace
	return {{.lower}}
}

{{- $snake := .snake}}
{{- $proj := .proj}}
{{- $id_type := .id_type}}

{{- define "id"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- lower (index . 0 0)}}
{{- end}}
{{- else -}}
id
{{- end}}
{{- end}}

{{- define "id_type"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- index . 0 1}}
{{- end}}
{{- else}}
{{- .id_type}}
{{- end}}
{{- end}}

// Db 返回db name
func ({{.lower}} *{{.camel}}) Db({{if .need_multi_table}}{{template "id" .keys}} {{template "id_type" dict "keys" .keys "id_type" .id_type}}{{end}}) (db string) {
	{{- if .need_multi_table}}
	{{- if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	return fmt.Sprintf("{{dbprefix $proj $snake}}_%d", {{if eq (index . 0 1) "string"}}utils.Hash33({{lower (index . 0 0)}}){{else}}int({{lower (index . 0 0)}}){{end}}%TableNum%DbNum)
	{{- end}}
	{{- else}}
	return fmt.Sprintf("{{dbprefix $proj $snake}}_%d", {{if eq $id_type "string"}}utils.Hash33(id){{else}}int(id){{end}}%TableNum%DbNum)
	{{- end}}
	{{- else}}
	return "{{dbprefix $proj $snake}}"
	{{- end}}
}

// Table 返回table name
func ({{.lower}} *{{.camel}}) Table({{if .need_multi_table}}{{template "id" .keys}} {{template "id_type" dict "keys" .keys "id_type" .id_type}}{{end}}) (table string) {
	{{- if .need_multi_table}}
	{{- if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	return fmt.Sprintf("{{$snake}}_%d", {{if eq (index . 0 1) "string"}}utils.Hash33({{lower (index . 0 0)}}) {{else}}int({{lower (index . 0 0)}}){{end}}%TableNum)
	{{- end}}
	{{- else}}
	return fmt.Sprintf("{{$snake}}_%d", {{if eq $id_type "string"}}utils.Hash33(id){{else}}int(id){{end}}%TableNum)
	{{- end}}
	{{- else}}
	return "{{$snake}}"
	{{- end}}
}

// GetC 返回db col
func ({{.lower}} *{{.camel}}) GetC({{if .need_multi_table}}{{template "id" .keys}} {{template "id_type" dict "keys" .keys "id_type" .id_type}}{{end}}) (c *mgo.Collection) {
	db, table := {{.lower}}.Db({{if .need_multi_table}}{{template "id" .keys}}{{end}}), {{.lower}}.Table({{if .need_multi_table}}{{template "id" .keys}}{{end}})
	session := mongo.DBS[db]
	sessionCopy := session.Copy()
	c = sessionCopy.DB(db).C(table)
	return
}

{{if .need_multi_table}}
// Groups 对多个对象按库表分组
func ({{.lower}} *{{.camel}}) Groups({{template "id" .keys}}s []{{template "id_type" dict "keys" .keys "id_type" .id_type}}) ({{template "id" .keys}}sGroup [][]{{template "id_type" dict "keys" .keys "id_type" .id_type}}) {
	if len({{template "id" .keys}}s) == 0 {
		return
	}

	groupsM := map[[2]string][]{{template "id_type" dict "keys" .keys "id_type" .id_type}}{} // key: [2]string{db, table}
	for _, {{template "id" .keys}} := range {{template "id" .keys}}s {
		db, table := {{.lower}}.Db({{template "id" .keys}}), {{.lower}}.Table({{template "id" .keys}})
		key := [2]string{db, table}
		groupsM[key] = append(groupsM[key], {{template "id" .keys}})
	}

	for _, {{template "id" .keys}}s := range groupsM {
		{{template "id" .keys}}sGroup = append({{template "id" .keys}}sGroup, {{template "id" .keys}}s)
	}
	return
}
{{end}}
`

var tplService = `package {{.snake}}

import "{{.pkg}}/lib"

// {{.camel}} 定义{{.camel}}类型
type {{.camel}} struct {
	*lib.Trace
}

func ({{.lower}} *{{.camel}}) WithTrace(trace *lib.Trace) *{{.camel}} {
	if {{.lower}} == nil {
		return nil
	}

	{{.lower}}.Trace = trace
	return {{.lower}}
}
`

var tplServiceGet = `package {{.snake}}

import (
	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/model"
	"{{.pkg}}/{{.proj}}_api"
)

{{- define "id_and_id_type"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- range $pos, $elem := .}} 
{{- if ne $pos 0 -}}
,
{{- end}}
{{- lower (index $elem 0)}} {{index $elem 1}}
{{- end}}
{{- end}}
{{- else -}}
id {{.id_type}}
{{- end}}
{{- end}}

{{- define "id"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- range $pos, $elem := .}} 
{{- if ne $pos 0 -}}
,
{{- end}}
{{- lower (index $elem 0)}}
{{- end}}
{{- end}}
{{- else -}}
id
{{- end}}
{{- end}}

{{- define "id_and_id_type_map"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- lower (index . 0 0)}} {{index . 0 1}}, {{lower (index . 1 0)}}s []{{index . 1 1}}
{{- end}}
{{- else -}}
ids []{{.id_type}}
{{- end}}
{{- end}}

{{- define "id_map"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- lower (index . 0 0)}}, {{lower (index . 1 0)}}s
{{- end}}
{{- else -}}
ids
{{- end}}
{{- end}}

{{- define "map_key"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- index . 1 1}}
{{- end}}
{{- else}}
{{- .id_type}}
{{- end}}
{{- end}}

{{- define "map_id"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- camel (index . 1 0)}}
{{- end}}
{{- else -}}
ID
{{- end}}
{{- end}}

// Get 定义获取操作
func ({{.lower}} *{{.camel}}) Get({{template "id_and_id_type" dict "keys" .keys "id_type" .id_type}}) ({{.lower}}API *{{.proj}}_api.{{.camel}}, err error) {
	fun := "service.{{.snake}}.{{.camel}}.Get"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	if {{.lower}}API, err = model.New{{.camel}}().WithTrace({{.lower}}.Trace).Get({{template "id" dict "keys" .keys "id_type" .id_type}}); err != nil {
		return
	}

	return
}

// Gets 定义批量获取操作
func ({{.lower}} *{{.camel}}) Gets({{template "id_and_id_type_map" dict "keys" .keys "id_type" .id_type}}) ({{.lower}}sAPI map[{{template "map_key" dict "keys" .keys "id_type" .id_type}}]{{if gt (len (split_keys .keys)) 1}}[]{{end}}*{{.proj}}_api.{{.camel}}, err error) {
	fun := "service.{{.snake}}.{{.camel}}.Gets"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	{{.lower}}sSliceAPI, err := model.New{{.camel}}().WithTrace({{.lower}}.Trace).Gets({{template "id_map" dict "keys" .keys "id_type" .id_type}})
	if err != nil {
		return
	}

	if len({{.lower}}sSliceAPI) == 0 {
		return
	}

	{{.lower}}sAPI = map[{{template "map_key" dict "keys" .keys "id_type" .id_type}}]{{if gt (len (split_keys .keys)) 1}}[]{{end}}*{{.proj}}_api.{{.camel}}{}
	for _, {{.lower}}API := range {{.lower}}sSliceAPI {
		{{- if gt (len (split_keys .keys)) 1}}
			{{.lower}}sAPI[{{.lower}}API.{{template "map_id" .keys}}] = append({{.lower}}sAPI[{{.lower}}API.{{template "map_id" .keys}}], {{.lower}}API)
		{{- else}}
			{{.lower}}sAPI[{{.lower}}API.{{template "map_id" .keys}}] = {{.lower}}API
		{{- end}}
	}

	return
}
`

var tplServiceAdd = `package {{.snake}}

import (
	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/model"
	"{{.pkg}}/{{.proj}}_api"
)

// Add 定义新增操作
func ({{.lower}} *{{.camel}}) Add({{.lower}}API *{{.proj}}_api.{{.camel}}) (err error) {
	fun := "service.{{.snake}}.{{.camel}}.Add"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	if err = model.New{{.camel}}().WithTrace({{.lower}}.Trace).Add({{.lower}}API); err != nil {
		return
	}

	return
}
`

var tplServiceDel = `package {{.snake}}

import (
	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/model"
)

{{- define "id_and_id_type"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- range $pos, $elem := .}} 
{{- if ne $pos 0 -}}
,
{{- end}}
{{- lower (index $elem 0)}} {{index $elem 1}}
{{- end}}
{{- end}}
{{- else -}}
id {{.id_type}}
{{- end}}
{{- end}}

{{- define "id"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- range $pos, $elem := .}} 
{{- if ne $pos 0 -}}
,
{{- end}}
{{- lower (index $elem 0)}}
{{- end}}
{{- end}}
{{- else -}}
id
{{- end}}
{{- end}}

// Del 定义删除操作
func ({{.lower}} *{{.camel}}) Del({{template "id_and_id_type" dict "keys" .keys "id_type" .id_type}}) (err error) {
	fun := "service.{{.snake}}.{{.camel}}.Del"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	if err = model.New{{.camel}}().WithTrace({{.lower}}.Trace).Del({{template "id" .keys}}); err != nil {
		return
	}

	return
}
`

var tplServiceUpdate = `package {{.snake}}

import (
	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/model"
	"{{.pkg}}/{{.proj}}_api"
)

// Update 定义更新操作
func ({{.lower}} *{{.camel}}) Update({{.lower}}API *{{.proj}}_api.{{.camel}}) (err error) {
	fun := "service.{{.snake}}.{{.camel}}.Update"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	if err = model.New{{.camel}}().WithTrace({{.lower}}.Trace).Update({{.lower}}API); err != nil {
		return
	}

	return
}

// Upsert 定义upsert操作
func ({{.lower}} *{{.camel}}) Upsert({{.lower}}API *{{.proj}}_api.{{.camel}}) (err error) {
	fun := "service.{{.snake}}.{{.camel}}.Upsert"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	if err = model.New{{.camel}}().WithTrace({{.lower}}.Trace).Upsert({{.lower}}API); err != nil {
		return
	}

	return
}
`

var tplServiceList = `package {{.snake}}

import (
	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/model"
	"{{.pkg}}/{{.proj}}_api"
)

{{- define "id_and_id_type_page"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}{{lower (index . 0 0)}} {{index . 0 1}}, {{end}}
{{- end}}
{{- end}}

{{- define "id_page"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}{{lower (index . 0 0)}}, {{end}}
{{- end}}
{{- end}}

{{- define "id_flow"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- lower (index . 0 0)}}, {{lower (index . 1 0)}}
{{- end}}
{{- else -}}
id
{{- end}}
{{- end}}

{{define "id_and_id_type_flow"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- lower (index . 0 0)}} {{index . 0 1}}, {{lower (index . 1 0)}} {{index . 1 1}}
{{- end}}
{{- else -}}
id {{.id_type}}
{{- end}}
{{- end}}

// PageList 定义page_list操作
func ({{.lower}} *{{.camel}}) PageList({{template "id_and_id_type_page" .keys}}offset, limit int) ({{.lower}}sAPI []*{{.proj}}_api.{{.camel}}, err error) {
	fun := "service.{{.snake}}.{{.camel}}.PageList"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	if {{.lower}}sAPI, err = model.New{{.camel}}().WithTrace({{.lower}}.Trace).PageList({{template "id_page" .keys}}offset, limit); err != nil {
		return
	}

	return
}

// FlowList 定义list操作
func ({{.lower}} *{{.camel}}) FlowList({{template "id_and_id_type_flow" dict "keys" .keys "id_type" .id_type}}, limit int) ({{.lower}}sAPI []*{{.proj}}_api.{{.camel}}, err error) {
	fun := "service.{{.snake}}.{{.camel}}.FlowList"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	if {{.lower}}sAPI, err = model.New{{.camel}}().WithTrace({{.lower}}.Trace).FlowList({{template "id_flow" dict "keys" .keys "id_type" .id_type}}, limit); err != nil {
		return
	}

	return
}
`

var tplModelGet = `package {{.snake}}

import (
	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}_api"
)

{{- $lower := .lower}}

{{- define "id_and_id_type_map"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- lower (index . 0 0)}} {{index . 0 1}}, {{lower (index . 1 0)}}s []{{index . 1 1}}
{{- end}}
{{- else -}}
ids []{{.id_type}}
{{- end}}
{{- end}}

{{- define "id_and_id_type"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- range $pos, $elem := .}} 
{{- if ne $pos 0 -}}
,
{{- end}}
{{- lower (index $elem 0)}} {{index $elem 1}}
{{- end}}
{{- end}}
{{- else -}}
id {{.id_type}}
{{- end}}
{{- end}}

{{- define "id"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- lower (index . 0 0)}}
{{- end}}
{{- else -}}
id
{{- end}}
{{- end}}

// Get 定义获取操作
func ({{.lower}} *{{.camel}}) Get({{template "id_and_id_type" dict "keys" .keys "id_type" .id_type}}) ({{.lower}}API *{{.proj}}_api.{{.camel}}, err error) {
	fun := "model.{{.snake}}.{{.camel}}.Get"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	c := {{.lower}}.GetC({{if .need_multi_table}}{{template "id" .keys}}{{end}})
	defer c.Database.Session.Close()

	{{if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	q := bson.M{
		{{- range $elem := .}}
		"{{snake (index $elem 0)}}": {{lower (index $elem 0)}},
		{{- end}}
	}
	err = c.Find(q).One(&{{$lower}}API)
	{{- end}}
	{{- else}}
	err = c.FindId(id).One(&{{$lower}}API)
	{{- end}}
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}
		err = nil
		return
	}

	return
}

// Gets 定义批量获取操作
func ({{.lower}} *{{.camel}}) Gets({{template "id_and_id_type_map" dict "keys" .keys "id_type" .id_type}}) ({{.lower}}sAPI []*{{.proj}}_api.{{.camel}}, err error) {
	fun := "model.{{.snake}}.{{.camel}}.Gets"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	{{if .need_multi_table}}
	{{- if gt (len (split_keys .keys)) 1}}
	c := {{.lower}}.GetC({{template "id" .keys}})
	defer c.Database.Session.Close()

	{{with (split_keys .keys) -}}
	q := bson.M{
		"{{snake (index . 0 0)}}": {{lower (index . 0 0)}},
		"{{snake (index . 1 0)}}": bson.M{
			"$in": {{lower (index . 1 0)}}s,
		},
	}
	{{- end}}

	err = c.Find(q).All(&{{.lower}}sAPI)
	if err != nil {
		return
	}
	{{- else}}
	{{template "id" .keys}}sGroup := {{.lower}}.Groups(ids)

	for _, {{template "id" .keys}}s := range {{template "id" .keys}}sGroup {
		{{template "id" .keys}} := {{template "id" .keys}}s[0]
		c := {{.lower}}.GetC({{template "id" .keys}})
		defer c.Database.Session.Close()

		q := bson.M{
			"_id": bson.M{
				"$in": {{template "id" .keys}}s,
			},
		}

		var {{.lower}}sTmpAPI []*{{.proj}}_api.{{.camel}}
		err = c.Find(q).All(&{{.lower}}sTmpAPI)
		if err != nil {
			return
		}

		{{.lower}}sAPI = append({{.lower}}sAPI, {{.lower}}sTmpAPI...)
	}
	{{- end}}
	{{- else}}
	c := {{.lower}}.GetC()
	defer c.Database.Session.Close()

	{{if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	q := bson.M{
		"{{snake (index . 0 0)}}": {{lower (index . 0 0)}},
		"{{snake (index . 1 0)}}": bson.M{
			"$in": {{lower (index . 1 0)}}s,
		},
	}
	{{- end}}
	{{- else}}
	q := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	{{- end}}

	err = c.Find(q).All(&{{.lower}}sAPI)
	if err != nil {
		return
	}
	{{- end}}

	return
}
`

var tplModelAdd = `package {{.snake}}

import (
	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}_api"
)

{{- define "id"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- camel (index . 0 0)}}
{{- end}}
{{- else -}}
ID
{{- end}}
{{- end}}

// Add 定义新增操作
func ({{.lower}} *{{.camel}}) Add({{.lower}}API *{{.proj}}_api.{{.camel}}) (err error) {
	fun := "model.{{.snake}}.{{.camel}}.Add"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	c := {{.lower}}.GetC({{if .need_multi_table}}{{.lower}}API.{{template "id" .keys}}{{end}})
	defer c.Database.Session.Close()

	err = c.Insert({{.lower}}API)
	if err != nil {
		return
	}

	return
}
`

var tplModelDel = `package {{.snake}}

import (
	"{{.pkg}}/lib"
)

{{- define "id_and_id_type"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- range $pos, $elem := .}} 
{{- if ne $pos 0 -}}
,
{{- end}}
{{- lower (index $elem 0)}} {{index $elem 1}}
{{- end}}
{{- end}}
{{- else -}}
id {{.id_type}}
{{- end}}
{{- end}}

{{- define "id"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- lower (index . 0 0)}}
{{- end}}
{{- else -}}
id
{{- end}}
{{- end}}


// Del 定义删除操作
func ({{.lower}} *{{.camel}}) Del({{template "id_and_id_type" dict "keys" .keys "id_type" .id_type}}) (err error) {
	fun := "model.{{.snake}}.{{.camel}}.Del"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	c := {{.lower}}.GetC({{if .need_multi_table}}{{template "id" .keys}}{{end}})
	defer c.Database.Session.Close()

	{{if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	q := bson.M{
		{{- range $elem := .}}
		"{{snake (index $elem 0)}}": {{lower (index $elem 0)}},
		{{- end}}
	}
	err = c.Remove(q)
	{{- end}}
	{{- else}}
	err = c.RemoveId(id)
	{{- end}}
	if err != nil {
		return
	}

	return
}
`

var tplModelUpdate = `package {{.snake}}

import (
	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}_api"
)

{{- $lower := .lower}}

{{- define "id"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- camel (index . 0 0)}}
{{- end}}
{{- else -}}
ID
{{- end}}
{{- end}}

{{- define "id_and_id_type"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- range $pos, $elem := .}} 
{{- if ne $pos 0 -}}
,
{{- end}}
{{- lower (index $elem 0)}} {{index $elem 1}}
{{- end}}
{{- end}}
{{- else -}}
id {{.id_type}}
{{- end}}
{{- end}}

// Update 定义更新操作
func ({{.lower}} *{{.camel}}) Update({{.lower}}API *{{.proj}}_api.{{.camel}}) (err error) {
	fun := "model.{{.snake}}.{{.camel}}.Update"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	c := {{.lower}}.GetC({{if .need_multi_table}}{{.lower}}API.{{template "id" .keys}}{{end}})
	defer c.Database.Session.Close()

	{{if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	q := bson.M{
		{{- range $elem := .}}
		"{{snake (index $elem 0)}}": {{$lower}}API.{{camel (index $elem 0)}},
		{{- end}}
	}
	err = c.Update(q, {{$lower}}API)
	{{- end}}
	{{- else}}
	err = c.UpdateId({{$lower}}API.ID, {{$lower}}API)
	{{- end}}
	if err != nil {
		return
	}

	return
}

// Upsert 定义upsert操作
func ({{.lower}} *{{.camel}}) Upsert({{.lower}}API *{{.proj}}_api.{{.camel}}) (err error) {
	fun := "model.{{.snake}}.{{.camel}}.Upsert"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	c := {{.lower}}.GetC({{if .need_multi_table}}{{.lower}}API.{{template "id" .keys}}{{end}})
	defer c.Database.Session.Close()

	{{if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	q := bson.M{
		{{- range $elem := .}}
		"{{snake (index $elem 0)}}": {{$lower}}API.{{camel (index $elem 0)}},
		{{- end}}
	}
	_, err = c.Upsert(q, {{$lower}}API)
	{{- end}}
	{{- else}}
	_, err = c.UpsertId({{$lower}}API.ID, {{$lower}}API)
	{{- end}}
	if err != nil {
		return
	}

	return
}
`

var tplModelList = `package {{.snake}}

import (
	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}_api"
)

{{- $id_type := .id_type}}

{{- define "id_and_id_type_page"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}{{lower (index . 0 0)}} {{index . 0 1}}, {{end}}
{{- end}}
{{- end}}

{{- define "id_page"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- lower (index . 0 0)}}
{{- end}}
{{- else -}}
*new({{.id_type}})
{{- end}}
{{- end}}


{{- define "id_flow"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- lower (index . 0 0)}}
{{- end}}
{{- else -}}
id
{{- end}}
{{- end}}

{{- define "id_and_id_type_flow"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys)}}
{{- lower (index . 0 0)}} {{index . 0 1}}, {{lower (index . 1 0)}} {{index . 1 1}}
{{- end}}
{{- else -}}
id {{.id_type}}
{{- end}}
{{- end}}

{{- define "id_sort"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- snake (index . 1 0)}}
{{- end}}
{{- else -}}
_id
{{- end}}
{{- end}}

// PageList 定义page_list操作
func ({{.lower}} *{{.camel}}) PageList({{template "id_and_id_type_page" .keys}}offset, limit int) ({{.lower}}sAPI []*{{.proj}}_api.{{.camel}}, err error) {
	fun := "model.{{.snake}}.{{.camel}}.PageList"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	c := {{.lower}}.GetC({{if .need_multi_table}}{{template "id_page" dict "keys" .keys "id_type" .id_type}}{{end}})
	defer c.Database.Session.Close()

	q := bson.M{}

	{{- if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	q["{{snake (index . 0 0)}}"] = {{lower (index . 0 0)}}
	{{- end}}
	{{- end}}
	err = c.Find(q).Sort("-{{template "id_sort" .keys}}").Skip(offset).Limit(limit).All(&{{.lower}}sAPI)
	if err != nil {
		return
	}

	return
}

// FlowList 定义flow_list操作
func ({{.lower}} *{{.camel}}) FlowList({{template "id_and_id_type_flow" dict "keys" .keys "id_type" .id_type}}, limit int) ({{.lower}}sAPI []*{{.proj}}_api.{{.camel}}, err error) {
	fun := "model.{{.snake}}.{{.camel}}.FlowList"
	defer lib.TraceMe({{.lower}}.Trace, fun)()

	c := {{.lower}}.GetC({{if .need_multi_table}}{{template "id_flow" .keys}}{{end}})
	defer c.Database.Session.Close()

	q := bson.M{}

	{{- if gt (len (split_keys .keys)) 1}}
	{{- with (split_keys .keys)}}
	q["{{snake (index . 0 0)}}"] = {{lower (index . 0 0)}}
	if {{lower (index . 1 0)}} != *new({{lower (index . 1 1)}}) {
		q["{{snake (index . 1 0)}}"] = bson.M{
			"$lt": {{lower (index . 1 0)}},
		}
	}
	{{- end}}
	{{- else}}
	if id != *new({{.id_type}}) {
		q["_id"] = bson.M{
			"$lt": id,
		}
	}
	{{- end}}
	err = c.Find(q).Sort("-{{template "id_sort" .keys}}").Limit(limit).All(&{{.lower}}sAPI)
	if err != nil {
		return
	}

	return
}
`

var tplFunc = `func New{{.camel}}() *{{.snake}}.{{.camel}} {
	return &{{.snake}}.{{.camel}}{}
}
`

var tplController = `package {{.snake}}

import "{{.pkg}}/lib"

// {{.camel}} 定义主对象
type {{.camel}} struct {
	lib.Base
}
`

var tplControllerGet = `package {{.snake}}

import (
	"encoding/json"
	"net/http"

	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/service"
	"{{.pkg}}/{{.proj}}_api"

	"github.com/simplejia/clog/api"
)

{{- define "id"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- range $pos, $elem := .}} 
{{- if ne $pos 0 -}}
,
{{- end -}}
req.{{camel (index $elem 0)}}
{{- end}}
{{- end}}
{{- else -}}
req.ID
{{- end}}
{{- end}}


// Get just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func ({{.lower}} *{{.camel}}) Get(w http.ResponseWriter, r *http.Request) {
	fun := "controller.{{.snake}}.{{.camel}}.Get"

	var req *{{.proj}}_api.{{.camel}}GetReq
	if err := json.Unmarshal({{.lower}}.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodePara)
		return
	}

	trace := lib.GetTrace({{.lower}})

	{{.lower}}API, err := service.New{{.camel}}().WithTrace(trace).Get({{template "id" .keys}})
	if err != nil {
		clog.Error("%s {{.lower}}.Get err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := (*{{.proj}}_api.{{.camel}}GetResp)({{.lower}}API)
	{{.lower}}.ReplyOk(w, resp)

	return
}
`

var tplControllerGets = `package {{.snake}}

import (
	"encoding/json"
	"net/http"

	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/service"
	"{{.pkg}}/{{.proj}}_api"

	"github.com/simplejia/clog/api"
)

{{- define "id_map"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys) -}}
req.{{camel (index . 0 0)}}, req.{{camel (index . 1 0)}}s
{{- end}}
{{- else -}}
req.IDS
{{- end}}
{{- end}}

// Gets just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func ({{.lower}} *{{.camel}}) Gets(w http.ResponseWriter, r *http.Request) {
	fun := "controller.{{.snake}}.{{.camel}}.Gets"

	var req *{{.proj}}_api.{{.camel}}GetsReq
	if err := json.Unmarshal({{.lower}}.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodePara)
		return
	}

	trace := lib.GetTrace({{.lower}})

	{{.lower}}sAPI, err := service.New{{.camel}}().WithTrace(trace).Gets({{template "id_map" dict "keys" .keys "id_type" .id_type}})
	if err != nil {
		clog.Error("%s {{.lower}}.Gets err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := {{.proj}}_api.{{.camel}}GetsResp({{.lower}}sAPI)
	{{.lower}}.ReplyOk(w, resp)

	return
}
`

var tplControllerAdd = `package {{.snake}}

import (
	"encoding/json"
	"net/http"

	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/api"
	"{{.pkg}}/{{.proj}}/service"
	"{{.pkg}}/{{.proj}}_api"

	"github.com/simplejia/clog/api"
)

{{- define "id"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- range $pos, $elem := .}} 
{{- if ne $pos 0 -}}
,
{{- end -}}
req.{{camel (index $elem 0)}}
{{- end}}
{{- end}}
{{- else -}}
req.ID
{{- end}}
{{- end}}

// Add just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func ({{.lower}} *{{.camel}}) Add(w http.ResponseWriter, r *http.Request) {
	fun := "controller.{{.snake}}.{{.camel}}.Add"

	var req *{{.proj}}_api.{{.camel}}AddReq
	if err := json.Unmarshal({{.lower}}.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodePara)
		return
	}

	trace := lib.GetTrace({{.lower}})

	{{.lower}}API := (*{{.proj}}_api.{{.camel}})(req)
	if err := service.New{{.camel}}().WithTrace(trace).Add({{.lower}}API); err != nil {
		clog.Error("%s {{.lower}}.Add err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &{{.proj}}_api.{{.camel}}AddResp{}
	{{.lower}}.ReplyOk(w, resp)

	// 进行一些异步处理的工作
	go lib.Updates({{.lower}}API, lib.ADD, nil)

	return
}
`

var tplControllerDel = `package {{.snake}}

import (
	"encoding/json"
	"net/http"

	"{{.pkg}}/{{.proj}}/api"
	"{{.pkg}}/{{.proj}}_api"

	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/service"

	"github.com/simplejia/clog/api"
)

{{- define "id"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- range $pos, $elem := .}} 
{{- if ne $pos 0 -}}
,
{{- end -}}
req.{{camel (index $elem 0)}}
{{- end}}
{{- end}}
{{- else -}}
req.ID
{{- end}}
{{- end}}

// Del just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func ({{.lower}} *{{.camel}}) Del(w http.ResponseWriter, r *http.Request) {
	fun := "controller.{{.snake}}.{{.camel}}.Del"

	var req *{{.proj}}_api.{{.camel}}DelReq
	if err := json.Unmarshal({{.lower}}.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodePara)
		return
	}

	trace := lib.GetTrace({{.lower}})

	if err := service.New{{.camel}}().WithTrace(trace).Del({{template "id" .keys}}); err != nil {
		clog.Error("%s {{.lower}}.Del err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &{{.proj}}_api.{{.camel}}DelResp{}
	{{.lower}}.ReplyOk(w, resp)

	{{.lower}}API := &{{.proj}}_api.{{.camel}}{
		{{- if gt (len (split_keys .keys)) 1}}
		{{- with (split_keys .keys)}}
		{{- range $elem := .}}
		{{camel (index $elem 0)}}: req.{{camel (index $elem 0)}},
		{{- end}}
		{{- end}}
		{{- else -}}
		ID: req.ID,
		{{- end}}
	}

	// 进行一些异步处理的工作
	go lib.Updates({{.lower}}API, lib.DELETE, nil)

	return
}
`

var tplControllerUpdate = `package {{.snake}}

import (
	"encoding/json"
	"net/http"

	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/api"
	"{{.pkg}}/{{.proj}}/service"
	"{{.pkg}}/{{.proj}}_api"

	"github.com/simplejia/clog/api"
)

{{- define "id"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}
{{- range $pos, $elem := .}} 
{{- if ne $pos 0 -}}
,
{{- end -}}
req.{{camel (index $elem 0)}}
{{- end}}
{{- end}}
{{- else -}}
req.ID
{{- end}}
{{- end}}

// Update just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func ({{.lower}} *{{.camel}}) Update(w http.ResponseWriter, r *http.Request) {
	fun := "controller.{{.snake}}.{{.camel}}.Update"

	var req *{{.proj}}_api.{{.camel}}UpdateReq
	if err := json.Unmarshal({{.lower}}.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodePara)
		return
	}

	trace := lib.GetTrace({{.lower}})

	{{.lower}}API := (*{{.proj}}_api.{{.camel}})(req)
	if err := service.New{{.camel}}().WithTrace(trace).Update({{.lower}}API); err != nil {
		clog.Error("%s {{.lower}}.Update err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &{{.proj}}_api.{{.camel}}UpdateResp{}
	{{.lower}}.ReplyOk(w, resp)

	// 进行一些异步处理的工作
	go lib.Updates({{.lower}}API, lib.UPDATE, nil)

	return
}
`

var tplControllerPageList = `package {{.snake}}

import (
	"encoding/json"
	"net/http"

	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/service"
	"{{.pkg}}/{{.proj}}_api"

	"github.com/simplejia/clog/api"
)

{{- define "id"}}
{{- if gt (len (split_keys .)) 1}}
{{- with (split_keys .)}}req.{{camel (index . 0 0)}}, {{end}}
{{- end}}
{{- end}}

// PageList just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func ({{.lower}} *{{.camel}}) PageList(w http.ResponseWriter, r *http.Request) {
	fun := "controller.{{.snake}}.{{.camel}}.PageList"

	var req *{{.proj}}_api.{{.camel}}PageListReq
	if err := json.Unmarshal({{.lower}}.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodePara)
		return
	}

	trace := lib.GetTrace({{.lower}})

	limitMore := req.Limit + 1

	{{.lower}}sAPI, err := service.New{{.camel}}().WithTrace(trace).PageList({{template "id" .keys}}req.Offset, limitMore)
	if err != nil {
		clog.Error("%s {{.lower}}.PageList err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodeSrv)
		return
	}

	n := len({{.lower}}sAPI)
	if n == 0 {
		{{.lower}}.ReplyOk(w, nil)
		return
	}

	more := false
	if n == limitMore {
		more = true
		{{.lower}}sAPI = {{.lower}}sAPI[:req.Limit]
	}

	resp := &{{.proj}}_api.{{.camel}}PageListResp{
		List:   {{.lower}}sAPI,
		Offset: req.Offset + len({{.lower}}sAPI),
		More: more,
	}
	{{.lower}}.ReplyOk(w, resp)

	return
}
`

var tplControllerFlowList = `package {{.snake}}

import (
	"encoding/json"
	"net/http"

	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/service"
	"{{.pkg}}/{{.proj}}_api"

	"github.com/simplejia/clog/api"
)

{{- define "id"}}
{{- if gt (len (split_keys .keys)) 1}}
{{- with (split_keys .keys) -}}
req.{{camel (index . 0 0)}}, *new({{lower (index . 1 1)}})
{{- end}}
{{- else -}}
*new({{.id_type}})
{{- end}}
{{- end}}

// FlowList just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func ({{.lower}} *{{.camel}}) FlowList(w http.ResponseWriter, r *http.Request) {
	fun := "controller.{{.snake}}.{{.camel}}.FlowList"

	var req *{{.proj}}_api.{{.camel}}FlowListReq
	if err := json.Unmarshal({{.lower}}.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodePara)
		return
	}

	trace := lib.GetTrace({{.lower}})

	limitMore := req.Limit + 1

	{{.lower}}sAPI, err := service.New{{.camel}}().WithTrace(trace).FlowList({{template "id" dict "keys" .keys "id_type" .id_type}}, limitMore)
	if err != nil {
		clog.Error("%s {{.lower}}.FlowList err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodeSrv)
		return
	}

	n := len({{.lower}}sAPI)
	if n == 0 {
		{{.lower}}.ReplyOk(w, nil)
		return
	}

	more := false
	if n == limitMore {
		more = true
		{{.lower}}sAPI = {{.lower}}sAPI[:req.Limit]
	}

	resp := &{{.proj}}_api.{{.camel}}FlowListResp{
		List:   {{.lower}}sAPI,
		More: more,
	}
	{{.lower}}.ReplyOk(w, resp)

	return
}
`

var tplControllerUpsert = `package {{.snake}}

import (
	"encoding/json"
	"net/http"

	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}/api"
	"{{.pkg}}/{{.proj}}/service"
	"{{.pkg}}/{{.proj}}_api"

	"github.com/simplejia/clog/api"
)

// Upsert just for skel
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func ({{.lower}} *{{.camel}}) Upsert(w http.ResponseWriter, r *http.Request) {
	fun := "controller.{{.snake}}.{{.camel}}.Upsert"

	var req *{{.proj}}_api.{{.camel}}UpsertReq
	if err := json.Unmarshal({{.lower}}.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodePara)
		return
	}

	trace := lib.GetTrace({{.lower}})

	{{.lower}}API := (*{{.proj}}_api.{{.camel}})(req)

	if err := service.New{{.camel}}().WithTrace(trace).Upsert({{.lower}}API); err != nil {
		clog.Error("%s {{.lower}}.Upsert err: %v, req: %v", fun, err, req)
		{{.lower}}.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &{{.proj}}_api.{{.camel}}UpsertResp{}
	{{.lower}}.ReplyOk(w, resp)

	// 进行一些异步处理的工作
	go lib.Updates({{.lower}}API, lib.ADD, nil)

	return
}
`

var tplTest = `package {{.snake}}

import (
	"encoding/json"
	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}_api"
)

// Add 封装controller.Add操作
func Add(req *{{.proj}}_api.{{.camel}}AddReq) (resp *{{.proj}}_api.{{.camel}}AddResp, result *lib.Resp, err error) {
	c := &{{.camel}}{}
	body, err := lib.TestPost(c.Add, req)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.proj}}_api.{{.camel}}AddResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data

	return
}

// Update 封装controller.Update操作
func Update(req *{{.proj}}_api.{{.camel}}UpdateReq) (resp *{{.proj}}_api.{{.camel}}UpdateResp, result *lib.Resp, err error) {
	c := &{{.camel}}{}
	body, err := lib.TestPost(c.Update, req)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.proj}}_api.{{.camel}}UpdateResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data

	return
}

// Del 封装controller.Del操作
func Del(req *{{.proj}}_api.{{.camel}}DelReq) (resp *{{.proj}}_api.{{.camel}}DelResp, result *lib.Resp, err error) {
	c := &{{.camel}}{}
	body, err := lib.TestPost(c.Del, req)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.proj}}_api.{{.camel}}DelResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data

	return
}

// Get 封装controller.Get操作
func Get(req *{{.proj}}_api.{{.camel}}GetReq) (resp *{{.proj}}_api.{{.camel}}GetResp, result *lib.Resp, err error) {
	c := &{{.camel}}{}
	body, err := lib.TestPost(c.Get, req)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.proj}}_api.{{.camel}}GetResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}

// Gets 封装controller.Gets操作
func Gets(req *{{.proj}}_api.{{.camel}}GetsReq) (resp *{{.proj}}_api.{{.camel}}GetsResp, result *lib.Resp, err error) {
	c := &{{.camel}}{}
	body, err := lib.TestPost(c.Gets, req)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.proj}}_api.{{.camel}}GetsResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}

// PageList 封装controller.PageList操作
func PageList(req *{{.proj}}_api.{{.camel}}PageListReq) (resp *{{.proj}}_api.{{.camel}}PageListResp, result *lib.Resp, err error) {
	c := &{{.camel}}{}
	body, err := lib.TestPost(c.PageList, req)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.proj}}_api.{{.camel}}PageListResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}

// FlowList 封装controller.FlowList操作
func FlowList(req *{{.proj}}_api.{{.camel}}FlowListReq) (resp *{{.proj}}_api.{{.camel}}FlowListResp, result *lib.Resp, err error) {
	c := &{{.camel}}{}
	body, err := lib.TestPost(c.FlowList, req)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.proj}}_api.{{.camel}}FlowListResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}

// Upsert 封装controller.Upsert操作
func Upsert(req *{{.proj}}_api.{{.camel}}UpsertReq) (resp *{{.proj}}_api.{{.camel}}UpsertResp, result *lib.Resp, err error) {
	c := &{{.camel}}{}
	body, err := lib.TestPost(c.Upsert, req)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.proj}}_api.{{.camel}}UpsertResp` + " `" + `json:"data"` + "`" + `
	}{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return
	}
	if s.Ret != lib.CodeOk {
		result = &s.Resp
		return
	}

	resp = s.Data
	return
}
`

func snake(src string) string {
	thisUpper := false
	prevUpper := false
	thisNumber := false
	prevNumber := false

	buf := bytes.NewBufferString("")
	for i, v := range src {
		if v >= '0' && v <= '9' {
			thisNumber = true
			thisUpper = false
		} else if v >= 'A' && v <= 'Z' {
			thisNumber = false
			thisUpper = true
		} else {
			thisNumber = false
			thisUpper = false
		}
		nextLower := false
		if i+1 < len(src) {
			vNext := src[i+1]
			if vNext >= 'a' && vNext <= 'z' {
				nextLower = true
			}
		}
		if i > 0 && ((thisNumber && !prevNumber) || (!thisNumber && prevNumber) || (thisUpper && (!prevUpper || nextLower))) {
			buf.WriteRune('_')
		}
		prevUpper = thisUpper
		prevNumber = thisNumber
		buf.WriteRune(v)
	}
	return strings.ToLower(buf.String())
}

func camel(src string) string {
	prevUnderline := true

	buf := bytes.NewBufferString("")
	for _, v := range src {
		if v == '_' {
			prevUnderline = true
			continue
		}

		if prevUnderline {
			buf.WriteRune(unicode.ToUpper(v))
			prevUnderline = false
		} else {
			buf.WriteRune(v)
		}
	}

	return buf.String()
}

func lower(str string) string {
	if len(str) == 0 {
		return ""
	}

	str = camel(str)
	return strings.ToLower(str[0:1]) + str[1:]
}

func exit(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	println()
	println("Failed!")
	os.Exit(-1)
}

func gen(t, proj, name, file string) (err error) {
	m := map[string]interface{}{
		"pkg":              rootPkg,
		"proj":             proj,
		"snake":            snake(name),
		"lower":            lower(name),
		"camel":            camel(name),
		"id_type":          idType,
		"need_multi_table": needMultiTable,
		"db_num":           dbNum,
		"table_num":        tblNum,
		"keys":             keys,
	}

	if timeout != 0 {
		if timeout < time.Second {
			m["timeout"] = fmt.Sprintf("time.Millisecond * %d", timeout/time.Millisecond)
		} else {
			m["timeout"] = fmt.Sprintf("time.Second * %d", timeout/time.Second)
		}
	}
	if connectTimeout != 0 {
		if connectTimeout < time.Second {
			m["connect_timeout"] = fmt.Sprintf("time.Millisecond * %d", connectTimeout/time.Millisecond)
		} else {
			m["connect_timeout"] = fmt.Sprintf("time.Second * %d", connectTimeout/time.Second)
		}
	}

	funcMap := template.FuncMap{
		"split_keys": splitKeys,
		"dict":       dict,
		"camel":      camel,
		"lower":      lower,
		"snake":      snake,
		"dbprefix":   dbprefix,
	}
	tpl := template.Must(template.New("tpl").Funcs(funcMap).Parse(t))

	buf := bytes.NewBuffer(nil)
	if err = tpl.Execute(buf, m); err != nil {
		return
	}

	content, err := imports.Process("", buf.Bytes(), nil)
	if err != nil {
		return
	}

	if _, err = os.Stat(file); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		err = nil
	} else {
		return
	}

	dir := filepath.Dir(file)
	if _, err = os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		err = nil
		if err = os.MkdirAll(dir, 0755); err != nil {
			return
		}
	}
	if err = ioutil.WriteFile(file, content, 0644); err != nil {
		return
	}

	return
}

func addImport(file, pkg string) (err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return
	}

	bingo := false

	for _, decl := range f.Decls {
		v, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if v.Tok != token.IMPORT {
			continue
		}
		for _, spec := range v.Specs {
			if strconv.Quote(pkg) == spec.(*ast.ImportSpec).Path.Value {
				bingo = true
				break
			}
		}
	}

	if bingo {
		return
	}

	if !astutil.AddImport(fset, f, pkg) {
		return errors.New("add import package fail")
	}

	ast.SortImports(fset, f)

	buffer := bytes.NewBuffer(nil)
	if err = printer.Fprint(buffer, fset, f); err != nil {
		return
	}

	if err = ioutil.WriteFile(file, buffer.Bytes(), 0644); err != nil {
		return
	}
	return
}

func addFunc(t, name, file string) (err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return
	}

	bingo := false

	for _, decl := range f.Decls {
		mdecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if mdecl.Recv != nil {
			continue
		}

		if mdecl.Name.String() == "New"+camel(name) {
			bingo = true
			break
		}
	}

	if bingo {
		return
	}

	content := bytes.NewBuffer(nil)

	m := map[string]interface{}{
		"snake": snake(name),
		"camel": camel(name),
	}
	tpl := template.Must(template.New("tpl").Parse(t))
	if err = tpl.Execute(content, m); err != nil {
		return
	}

	fd, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer fd.Close()

	if _, err = fd.WriteString("\n" + content.String()); err != nil {
		return
	}
	return
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func splitKeys(keys string) [][2]string {
	a := strings.Split(keys, ",")
	if len(a)%2 != 0 {
		return nil
	}

	r := [][2]string{}
	for i := 0; i < len(a); i += 2 {
		r = append(r, [2]string{a[i], a[i+1]})
	}

	return r
}

func dbprefix(proj, db string) string {
	if strings.HasPrefix(db, proj) {
		return db
	}
	return proj + "_" + db
}

func main() {
	flag.StringVar(&rootPkg, "pkg", "", "Specify package name")
	flag.StringVar(&name, "name", "", "Specify module name")
	flag.DurationVar(&timeout, "timeout", time.Minute, "Specify timeout")
	flag.DurationVar(&connectTimeout, "connect_timeout", time.Millisecond*40, "Specify connect timeout")
	flag.StringVar(&idType, "id_type", "int64", "Specify id type")
	flag.StringVar(&keys, "keys", "", "Specify keys, ex: pid,int64,rid,int64")
	flag.BoolVar(&needMultiTable, "need_multi_table", false, "Specify if need multi table")
	flag.IntVar(&dbNum, "db_num", 1, "Specify number of db if multi table")
	flag.IntVar(&tblNum, "table_num", 64, "Specify number of table if multi table")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "A tiny tool, used to generate module's crud code\n")
		fmt.Fprintf(os.Stderr, "version: 1.11, Created by simplejia [11/2018]\n\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if rootPkg == "" || name == "" {
		flag.Usage()
		return
	}

	if keys != "" {
		if len(strings.Split(keys, ","))%2 != 0 {
			flag.Usage()
			return
		}
	}

	println("Begin generate crud")

	wd, err := os.Getwd()
	if err != nil {
		return
	}
	base := filepath.Base(wd)

	m := map[string]string{
		filepath.Join("..", base+"_api", snake(name)+".go"):           tplOuterAPI,
		filepath.Join("..", base+"_api", snake(name)+"_get.go"):       tplOuterAPIGet,
		filepath.Join("..", base+"_api", snake(name)+"_gets.go"):      tplOuterAPIGets,
		filepath.Join("..", base+"_api", snake(name)+"_add.go"):       tplOuterAPIAdd,
		filepath.Join("..", base+"_api", snake(name)+"_del.go"):       tplOuterAPIDel,
		filepath.Join("..", base+"_api", snake(name)+"_update.go"):    tplOuterAPIUpdate,
		filepath.Join("..", base+"_api", snake(name)+"_page_list.go"): tplOuterAPIPageList,
		filepath.Join("..", base+"_api", snake(name)+"_flow_list.go"): tplOuterAPIFlowList,
		filepath.Join("..", base+"_api", snake(name)+"_upsert.go"):    tplOuterAPIUpsert,
		filepath.Join("model", snake(name), snake(name)+".go"):        tplModel,
		filepath.Join("model", snake(name), "get.go"):                 tplModelGet,
		filepath.Join("model", snake(name), "add.go"):                 tplModelAdd,
		filepath.Join("model", snake(name), "del.go"):                 tplModelDel,
		filepath.Join("model", snake(name), "update.go"):              tplModelUpdate,
		filepath.Join("model", snake(name), "list.go"):                tplModelList,
		filepath.Join("service", snake(name), snake(name)+".go"):      tplService,
		filepath.Join("service", snake(name), "get.go"):               tplServiceGet,
		filepath.Join("service", snake(name), "add.go"):               tplServiceAdd,
		filepath.Join("service", snake(name), "del.go"):               tplServiceDel,
		filepath.Join("service", snake(name), "update.go"):            tplServiceUpdate,
		filepath.Join("service", snake(name), "list.go"):              tplServiceList,
		filepath.Join("controller", snake(name), snake(name)+".go"):   tplController,
		filepath.Join("controller", snake(name), "get.go"):            tplControllerGet,
		filepath.Join("controller", snake(name), "gets.go"):           tplControllerGets,
		filepath.Join("controller", snake(name), "add.go"):            tplControllerAdd,
		filepath.Join("controller", snake(name), "del.go"):            tplControllerDel,
		filepath.Join("controller", snake(name), "update.go"):         tplControllerUpdate,
		filepath.Join("controller", snake(name), "page_list.go"):      tplControllerPageList,
		filepath.Join("controller", snake(name), "flow_list.go"):      tplControllerFlowList,
		filepath.Join("controller", snake(name), "upsert.go"):         tplControllerUpsert,
		filepath.Join("controller", snake(name), "test.go"):           tplTest,
	}

	for file, tpl := range m {
		if err := gen(tpl, base, name, file); err != nil {
			exit("gen file: %s, err: %v", file, err)
		}
	}

	for _, level := range []string{"model", "service"} {
		file := filepath.Join(level, level+".go")
		pkg := fmt.Sprintf("%s/%s/%s/%s", rootPkg, base, level, snake(name))
		if err := addImport(file, pkg); err != nil {
			exit("add import fail, file: %s, err: %v, package: %s", file, err, pkg)
		}

		if err := addFunc(tplFunc, name, file); err != nil {
			exit("add func fail, file: %s, err: %v, func: %s", file, err, tplFunc)
		}
	}

	println("Success!")
}
