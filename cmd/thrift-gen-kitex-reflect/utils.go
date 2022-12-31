package main

import (
	"os"
)

func debugln(fn func()) {
	if os.Getenv("DEBUG_KITEX_REFLECT_PLUGIN") != "" {
		fn()
	}
}

func int32p(x int) *int32 {
	ret := int32(x)
	return &ret
}

func setNonZeroValue[T comparable](dst **T, val T) {
	var zero T
	if val != zero {
		*dst = &val
	}
}
