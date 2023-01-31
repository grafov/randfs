// This file defines how to make directory tree.
// This is a part of randfs library.
//
// Copyright 2023 The Project Developers. Under MIT license. See
// LICENSE file at the top-level directory of this distribution and at
// https://github.com/grafov/randfs
//
// ॐ तारे तुत्तारे तुरे स्व
package randfs

type Tree struct {
	depth              int
	minLimit, maxLimit int
}

func Make(root string, opts ...option) error {
	var t Tree
	for _, o := range opts {
		o(&t)
	}
	return nil
}

type option func(*Tree)

func Depth(n int) option {
	return func(t *Tree) {
		t.depth = n
	}
}

func Limit(min, max int) option {
	return func(t *Tree) {
		t.minLimit = min
		t.maxLimit = max
	}
}
