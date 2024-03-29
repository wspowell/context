//go:build !release
// +build !release

package context

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Duplicated from stdlib net/http2/gotrack.go, because the Go Authors know
// what's best for us, and this is the only portable way to get the goroutine
// ID, even though this number is readily available in the go internals as
// a real integer.

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"

	"github.com/wspowell/errors"
)

// ID is the ID number of a goroutine.
type goroutineId uint64

func (self goroutineId) isSameGoroutine() bool {
	return self == curID()
}

const goroutineSpace = "goroutine "

// CurID gets the ID number of the current goroutine.
//
// According to the Go Authors, using this number will cause, among other
// things: sponaneous combustion, incurable insanity, and rapid acute
// cardiac lithomophosis (RACL).
//
// That said, it turns out that this functionality is very important for
// implemening things such as the errors.Annotate functionality.
//
// If, at some point, the Go Authors deign to provide this functionality via the
// runtime package, this method and type will be deleted immediately in favor of
// that.
//
// (if you are a Go Author, please, please, PLEASE provide this as part of the
// stdlib...).
func curID() goroutineId {
	bp, ok := littleBuf.Get().(*[]byte)
	if !ok {
		return goroutineId(0)
	}
	defer littleBuf.Put(bp)
	b := *bp
	b = b[:runtime.Stack(b, false)]
	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, []byte(goroutineSpace))
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := parseUintBytes(b, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}

	return goroutineId(n)
}

// nolint:gochecknoglobals // reason: sync pool
var littleBuf = sync.Pool{
	New: func() any {
		buf := make([]byte, 64)

		return &buf
	},
}

// parseUintBytes is like strconv.ParseUint, but using a []byte.
func parseUintBytes(s []byte, base int, bitSize int) (n uint64, err error) {
	var cutoff, maxVal uint64

	if bitSize == 0 {
		bitSize = int(strconv.IntSize)
	}

	s0 := s
	switch {
	case len(s) < 1:
		err = strconv.ErrSyntax

		goto Error

	case base >= 2 && base <= 36:
		// valid base; nothing to do

	case base == 0:
		// Look for octal, hex prefix.
		switch {
		case s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X'):
			base = 16
			s = s[2:]
			if len(s) < 1 {
				err = strconv.ErrSyntax

				goto Error
			}
		case s[0] == '0':
			base = 8
		default:
			base = 10
		}

	default:
		err = errors.New("allyourbase", "invalid base "+strconv.Itoa(base))

		goto Error
	}

	n = 0
	cutoff = cutoff64(base)
	maxVal = 1<<uint(bitSize) - 1

	for i := 0; i < len(s); i++ {
		var v byte
		d := s[i]
		switch {
		// nolint:gocritic // reason: easier to read
		case '0' <= d && d <= '9':
			v = d - '0'
		// nolint:gocritic // reason: easier to read
		case 'a' <= d && d <= 'z':
			v = d - 'a' + 10
		// nolint:gocritic // reason: easier to read
		case 'A' <= d && d <= 'Z':
			v = d - 'A' + 10
		default:
			n = 0
			err = strconv.ErrSyntax

			goto Error
		}
		if int(v) >= base {
			n = 0
			err = strconv.ErrSyntax

			goto Error
		}

		if n >= cutoff {
			// n*base overflows
			n = 1<<64 - 1
			err = strconv.ErrRange

			goto Error
		}
		n *= uint64(base)

		n1 := n + uint64(v)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			n = 1<<64 - 1
			err = strconv.ErrRange

			goto Error
		}
		n = n1
	}

	return n, nil

Error:
	return n, &strconv.NumError{Func: "ParseUint", Num: string(s0), Err: err}
}

// Return the first number n such that n*base >= 1<<64.
func cutoff64(base int) uint64 {
	if base < 2 {
		return 0
	}

	return (1<<64-1)/uint64(base) + 1
}
