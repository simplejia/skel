package api

// Code 用于定义返回码(detail)
type Code int

// 定义各种返回码(detail)
const (
	CodeNone Code = iota
	CodeSkelExist
	CodeSkelNotExist
)

// CodeMap 定义返回码对应的描述
var CodeMap = map[Code]string{
	CodeSkelExist:    "skel已存在",
	CodeSkelNotExist: "skel不存在",
}
