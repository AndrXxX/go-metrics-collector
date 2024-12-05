package utils

import "golang.org/x/exp/constraints"

func Pointer[T constraints.Ordered](v T) *T {
	return &v
}
