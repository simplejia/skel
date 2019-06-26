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
	rootPkg string
	name    string
)

func exit(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	println()
	println("Failed!")
	os.Exit(-1)
}

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
		if d, ok := n.(*ast.FuncDecl); ok && d.Recv == nil && d.Name.String() == "New"+camel(name) {
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
	flag.StringVar(&rootPkg, "pkg", "", "Specify package name")
	flag.StringVar(&name, "name", "", "Specify module name")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "A tiny tool, used to delete module's crud code\n")
		fmt.Fprintf(os.Stderr, "version: 1.11, Created by simplejia [11/2018]\n\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	rootPkg = strings.TrimPrefix(rootPkg, "/")
	if rootPkg != "" {
		if !strings.HasSuffix(rootPkg, "/") {
			rootPkg += "/"
		}
	}

	if name == "" {
		flag.Usage()
		return
	}

	println("Begin delete crud")

	wd, err := os.Getwd()
	if err != nil {
		return
	}
	base := filepath.Base(wd)

	for _, level := range []string{"controller", "model", "service"} {
		dir := filepath.Join(level, snake(name))
		if err := os.RemoveAll(dir); err != nil {
			exit("remove dir fail, dir: %s, err: %v", dir, err)
		}
	}

	for _, n := range []string{
		snake(name) + ".go",
		snake(name) + "_get.go",
		snake(name) + "_gets.go",
		snake(name) + "_add.go",
		snake(name) + "_del.go",
		snake(name) + "_update.go",
		snake(name) + "_upsert.go",
		snake(name) + "_page_list.go",
		snake(name) + "_flow_list.go",
	} {
		file := filepath.Join("..", base+"_api", n)
		if err := os.Remove(file); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			exit("remove file fail, file: %s, err: %v", file, err)
		}
	}

	for _, level := range []string{"model", "service"} {
		file := filepath.Join(level, level+".go")
		pkg := fmt.Sprintf("%s%s/%s/%s", rootPkg, base, level, snake(name))
		if err := delImport(file, pkg); err != nil {
			exit("del import fail, file: %s, err: %v, package: %s", file, err, pkg)
		}

		if err := delFunc(name, file); err != nil {
			exit("del func fail, file: %s, err: %v, func: %s", file, err, camel(name))
		}
	}

	println("Success!")
}
