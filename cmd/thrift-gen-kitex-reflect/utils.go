package main

import (
	"os"
	"strings"

	"github.com/cloudwego/thriftgo/parser"

	idl "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl"
)

var (
	enumTypeName  = "i32"
	idlTypeStruct = idl.TypePtr(idl.Type_STRUCT)
	trueVal       = true
)

func setNonZeroValue[T comparable](dst **T, val T) {
	var zero T
	if val != zero {
		*dst = &val
	}
}

func int32p(x int) *int32 {
	ret := int32(x)
	return &ret
}

func debugln(fn func()) {
	if os.Getenv("DEBUG_KITEX_REFLECT_PLUGIN") != "" {
		fn()
	}
}

func splitType(t string) (pkg, name string) {
	idx := strings.LastIndex(t, ".")
	if idx == -1 {
		return "", t
	}
	return t[:idx], t[idx+1:]
}

func findStructLike(tree *parser.Thrift, name string) (st *parser.StructLike, found bool) {
	st, found = tree.GetStruct(name)
	if !found {
		st, found = tree.GetUnion(name)
	}
	if !found {
		st, found = tree.GetException(name)
	}
	return
}
