// This file defines how to make names for files and directories.
// This is a part of randfs library.
//
// Copyright 2023 The Project Developers. Under MIT license. See
// LICENSE file at the top-level directory of this distribution and at
// https://github.com/grafov/randfs
//
// ॐ तारे तुत्तारे तुरे स्व
package name

import (
	"math/rand"
	"strings"
	"time"
	"unicode"
)

type Name struct {
	rnd       rand.Source
	maxLength int
	minLength int
	alphabet  []rune
	strcase   int // -1 lower, 0 mixed, 1 upper
}

// New returns a name for file or directory that could be formatted to
// a text string. The generation of names of single Name instance is
// not safe for concurrent use by multiple goroutines.
func New(opts ...option) *Name {
	n := Name{
		rnd:       rand.NewSource(time.Now().UnixNano()),
		minLength: 1,
		maxLength: 16,
	}
	for _, o := range opts {
		o(&n)
	}
	if len(n.alphabet) == 0 {
		// all latin letters by default
		n.alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	return &n
}

// String returns a string representing a name for file or directory.
func (n *Name) String() string {
	nlen := int(n.rnd.Int63() % int64(n.maxLength))
	if nlen < n.minLength {
		nlen = n.minLength
	}
	var res strings.Builder
	res.Grow(nlen)
	for i := 0; i < nlen; i++ {
		r := n.alphabet[n.rnd.Int63()%int64(len(n.alphabet))]
		switch {
		case n.strcase < 0:
			r = unicode.ToLower(r)
		case n.strcase > 0:
			r = unicode.ToUpper(r)

		}

		res.WriteRune(r)
	}
	return res.String()
}

type option func(*Name)

func Length(min, max int) option {
	return func(n *Name) {
		n.minLength = min
		n.maxLength = max
	}
}

func UseAlpha() option {
	const c = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return func(n *Name) {
		if len(n.alphabet) > 0 {
			n.alphabet = append(n.alphabet, []rune(c)...)
			return
		}
		n.alphabet = []rune(c)
	}
}

func UseHex() option {
	const c = "0123456789abcdef"
	return func(n *Name) {
		if len(n.alphabet) > 0 {
			n.alphabet = append(n.alphabet, []rune(c)...)
			return
		}
		n.alphabet = []rune(c)
	}
}

func UseDigits() option {
	const c = "0123456789"
	return func(n *Name) {
		if len(n.alphabet) > 0 {
			n.alphabet = append(n.alphabet, []rune(c)...)
			return
		}
		n.alphabet = []rune(c)
	}
}

func UseCustom(c string) option {
	return func(n *Name) {
		if len(n.alphabet) > 0 {
			n.alphabet = append(n.alphabet, []rune(c)...)
			return
		}
		n.alphabet = []rune(c)
	}
}

func LowerCase() option {
	return func(n *Name) {
		n.strcase = -1
	}
}

func UpperCase() option {
	return func(n *Name) {
		n.strcase = 1
	}
}
