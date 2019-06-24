package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/ast/astutil"
)

var (
	path string
)

func exit(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	println()
	println("Failed!")
	os.Exit(-1)
}

func snake(src string) string {
	if strings.IndexByte(src, '_') != -1 {
		return strings.ToLower(src)
	}

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

func delImport(file, pkg string) (err error) {
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

	if !bingo {
		return
	}

	if !astutil.DeleteImport(fset, f, pkg) {
		return errors.New("del import package fail")
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

func delFunc(name, file string) (err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return
	}

	n := astutil.Apply(f, func(c *astutil.Cursor) bool {
		n := c.Node()
		if d, ok := n.(*ast.FuncDecl); ok && d.Recv == nil && d.Name.String() == camel(name) {
			c.Delete()
		}
		return true
	}, nil)

	buffer := bytes.NewBuffer(nil)
	if err = format.Node(buffer, fset, n); err != nil {
		return
	}

	if err = ioutil.WriteFile(file, buffer.Bytes(), 0644); err != nil {
		return
	}
	return
}

func main() {
	flag.StringVar(&path, "path", "", "Specify url path")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "A tiny tool, used to delete project's api\n")
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

	println("Begin delete api")

	wd, err := os.Getwd()
	if err != nil {
		return
	}
	base := filepath.Base(wd)

	for _, n := range []string{
		filepath.Join("..", base+"_api", snake(typ)+"_"+snake(method)+".go"),
		filepath.Join("controller", snake(typ), snake(method)+".go"),
	} {
		file := n
		if err := os.Remove(file); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			exit("remove file fail, file: %s, err: %v", file, err)
		}
	}

	file := filepath.Join("controller", snake(typ), "test.go")
	if _, err := os.Stat(file); err != nil {
		if !os.IsNotExist(err) {
			exit("file stat err: %v, file: %s", err, file)
		}
	} else {
		if err := delFunc(method, file); err != nil {
			exit("del func fail, file: %s, err: %v, func: %s", file, err, camel(method))
		}
	}

	dir := filepath.Join("controller", snake(typ))
	if _, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			exit("dir stat err: %v, dir: %s", err, dir)
		}
	} else {
		bingo := false
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) (reterr error) {
			if err != nil {
				reterr = err
				return
			}

			file := info.Name()
			if strings.HasPrefix(file, ".") || file == snake(typ) || file == snake(typ)+".go" || file == snake(method)+".go" || file == "test.go" {
				return
			}

			bingo = true
			return
		}); err != nil {
			exit("walk dir err, dir: %s, err: %v", dir, err)
		}

		if !bingo {
			if err := os.RemoveAll(dir); err != nil {
				exit("remove dir fail, dir: %s, err: %v", dir, err)
			}
		}
	}

	println("Success!")
}
