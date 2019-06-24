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
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"
)

var (
	rootPkg        string
	path           string
	timeout        time.Duration
	connectTimeout time.Duration
)

var tplOuterAPI = `package {{.proj}}_api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/simplejia/utils"
	"{{.pkg}}/lib"
)

// {{.camel_type}}{{.camel_method}}Req 定义输入
type {{.camel_type}}{{.camel_method}}Req struct {
}

// Regular 用于参数校验
func ({{.lower_type}}{{.camel_method}}Req *{{.camel_type}}{{.camel_method}}Req) Regular() (ok bool) {
	if {{.lower_type}}{{.camel_method}}Req == nil {
		return
	}

	ok = true
	return
}

// {{.camel_type}}{{.camel_method}}Resp 定义输出
type {{.camel_type}}{{.camel_method}}Resp struct {
}

func {{.camel_type}}{{.camel_method}}(name string, req *{{.camel_type}}{{.camel_method}}Req, trace *lib.Trace) (resp *{{.camel_type}}{{.camel_method}}Resp, result *lib.Resp, err error) {
	addr, err := lib.NameWrap(name)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("http://%s/%s", addr, "{{.snake_type}}/{{.snake_method}}")
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
		Data *{{.camel_type}}{{.camel_method}}Resp` + " `" + `json:"data"` + "`" + `
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

var tplController = `package {{.snake_type}}

import "{{.pkg}}/lib"

// {{.camel_type}} 定义主对象
type {{.camel_type}} struct {
	lib.Base
}
`

var tplControllerAction = `package {{.snake_type}}

import (
	"encoding/json"
	"net/http"

	clog "github.com/simplejia/clog/api"
	"{{.pkg}}/lib"
	"{{.pkg}}/{{.proj}}_api"
)

// {{.camel_method}} just for demo
// @prefilter({"Timeout":{"dur":"1s"}}, "Trace")
// @postfilter("Boss")
func ({{.lower_type}} *{{.camel_type}}) {{.camel_method}}(w http.ResponseWriter, r *http.Request) {
	fun := "controller.{{.lower_type}}.{{.camel_type}}.{{.camel_method}}"

	var req *{{.proj}}_api.{{.camel_type}}{{.camel_method}}Req
	if err := json.Unmarshal({{.lower_type}}.ReadBody(r), &req); err != nil || !req.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, req)
		{{.lower_type}}.ReplyFail(w, lib.CodePara)
		return
	}

	resp := &{{.proj}}_api.{{.camel_type}}{{.camel_method}}Resp{}
	{{.lower_type}}.ReplyOk(w, resp)

	return
}
`

var tplTestFunc = `// {{.camel_method}} 封装controller.{{.camel_method}}操作
func {{.camel_method}}(req *{{.proj}}_api.{{.camel_type}}{{.camel_method}}Req) (resp *{{.proj}}_api.{{.camel_type}}{{.camel_method}}Resp, result *lib.Resp, err error) {
	c := &{{.camel_type}}{}
	body, err := lib.TestPost(c.{{.camel_method}}, req)
	if err != nil {
		return
	}

	s := &struct {
		lib.Resp
		Data *{{.proj}}_api.{{.camel_type}}{{.camel_method}}Resp` + " `" + `json:"data"` + "`" + `
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

func gen(t, proj, typ, method, file string) (err error) {
	m := map[string]interface{}{
		"pkg":          rootPkg,
		"proj":         proj,
		"snake_type":   snake(typ),
		"lower_type":   lower(typ),
		"camel_type":   camel(typ),
		"snake_method": snake(method),
		"lower_method": lower(method),
		"camel_method": camel(method),
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

	tpl := template.Must(template.New("tpl").Parse(t))

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

func addFunc(t, proj, typ, method, file string) (err error) {
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

		if mdecl.Name.String() == camel(method) {
			bingo = true
			break
		}
	}

	if bingo {
		return
	}

	content := bytes.NewBuffer(nil)

	m := map[string]interface{}{
		"proj":         proj,
		"snake_type":   snake(typ),
		"lower_type":   lower(typ),
		"camel_type":   camel(typ),
		"snake_method": snake(method),
		"lower_method": lower(method),
		"camel_method": camel(method),
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

func main() {
	flag.StringVar(&rootPkg, "pkg", "", "Specify package name")
	flag.StringVar(&path, "path", "", "Specify url path")
	flag.DurationVar(&timeout, "timeout", time.Minute, "Specify timeout")
	flag.DurationVar(&connectTimeout, "connect_timeout", time.Millisecond*40, "Specify connect timeout")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "A tiny tool, used to generate project's api\n")
		fmt.Fprintf(os.Stderr, "version: 1.11, Created by simplejia [11/2018]\n\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if path == "" {
		flag.Usage()
		return
	}

	path = strings.TrimPrefix(path, "/")
	a := strings.Split(path, "/")
	if len(a) != 2 {
		fmt.Println("path must be alike xxx/xxx")
		return
	}

	typ, method := a[0], a[1]

	println("Begin generate api")

	wd, err := os.Getwd()
	if err != nil {
		return
	}
	base := filepath.Base(wd)

	m := map[string]string{
		filepath.Join("..", base+"_api", snake(typ)+"_"+snake(method)+".go"): tplOuterAPI,
		filepath.Join("controller", snake(typ), snake(typ)+".go"):            tplController,
		filepath.Join("controller", snake(typ), snake(method)+".go"):         tplControllerAction,
	}

	for file, tpl := range m {
		if err := gen(tpl, base, typ, method, file); err != nil {
			exit("gen file: %s, err: %v", file, err)
		}
	}

	file := filepath.Join("controller", snake(typ), "test.go")

	if _, err := os.Stat(file); err != nil {
		if !os.IsNotExist(err) {
			exit("file stat err: %v, file: %s", err, file)
		}
		dir := filepath.Dir(file)
		if _, err := os.Stat(dir); err != nil {
			if !os.IsNotExist(err) {
				exit("dir stat err: %v, dir: %s", err, dir)
			}
			if err := os.MkdirAll(dir, 0755); err != nil {
				exit("mkdir err: %v, dir: %s", err, dir)
			}
		}

		if err := ioutil.WriteFile(file, []byte("package "+snake(typ)), 0666); err != nil {
			exit("write file err: %v, file: %s", err, file)
		}
	}

	for _, pkg := range []string{
		"encoding/json",
		fmt.Sprintf("%s/lib", rootPkg),
		fmt.Sprintf("%s/%s/test", rootPkg, base),
		fmt.Sprintf("%s/%s_api", rootPkg, base),
	} {
		if err := addImport(file, pkg); err != nil {
			exit("add file: %s, err: %v, package: %s", file, err, pkg)
		}
	}

	if err := addFunc(tplTestFunc, base, typ, method, file); err != nil {
		exit("add file: %s, err: %v", file, err)
	}

	println("Success!")
	return
}
