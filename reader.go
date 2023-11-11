// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runeio

import (
	"errors"
	"io"
	"unicode/utf8"
)

// defaultBufSize is the default buffer size.
const defaultBufSize = 1024

var (
	// ErrBufferFull indicates that the current buffer size cannot support the operation.
	ErrBufferFull = errors.New("buffer full")

	// ErrInvalidUnreadRune is returned when a rune cannot be unread.
	ErrInvalidUnreadRune = errors.New("invalid use of UnreadRune")

	// ErrNegativeCount is returned when a negative size is given.
	ErrNegativeCount = errors.New("negative count")
)

// RuneReader implements buffered look-ahead for a io.RuneReader.
type RuneReader struct {
	// rd is the underlying rune reader.
	rd io.RuneReader

	// buf is the rune lookahead buffer.
	buf []rune

	// lastRune is the last read rune.
	lastRune rune

	// r is the read index into the buffer.
	r int

	// e is the end index of the buffer.
	e int

	// err is the last read error that occurred.
	err error
}

// NewReader returns a new RuneReader with whose buffer has the default size
// of 1024 runes. NewReader is provided an underlying io.RuneReader (such as
// bufio.Reader, or strings.Reader).
func NewReader(r io.RuneReader) *RuneReader {
	return NewReaderSize(r, defaultBufSize)
}

// NewReaderSize returns a new RuneReader with whose buffer has the
// specified size.
func NewReaderSize(r io.RuneReader, size int) *RuneReader {
	buf := make([]rune, size)
	rd := new(RuneReader)
	rd.reset(buf, r)
	return rd
}

// Read reads runes into p and returns the number of runes read. If the number
// of runes read is different than the size of p, then an error is returned
// explaining the reason.
func (r *RuneReader) Read(p []rune) (int, error) {
	i := 0
	var err error
	for err == nil && i < len(p) {
		var rn rune
		rn, _, err = r.ReadRune()
		if err != nil {
			return i, err
		}

		p[i] = rn
		i++
	}
	return i, err
}

// ReadRune reads a single UTF-8 encoded Unicode character and returns the
// rune and its size in bytes.
func (r *RuneReader) ReadRune() (rune, int, error) {
	r.fill(1)

	if r.Buffered() == 0 {
		return 0, 0, r.readErr()
	}

	rn := r.buf[r.r]
	r.r++
	r.lastRune = rn
	return rn, utf8.RuneLen(rn), nil
}

// Discard attempts to discard n runes and returns the number actually
// discarded. If the number of runes discarded is different than n, then an
// error is returned explaining the reason.
//
// Calling Discard prevents an UnreadRune call from succeeding until the next
// read operation.
func (r *RuneReader) Discard(n int) (int, error) {
	if n < 0 {
		return 0, ErrNegativeCount
	}

	for i := 0; i < n; i++ {
		_, _, err := r.ReadRune()
		if err != nil {
			return i, err
		}
	}

	r.lastRune = -1

	return n, nil
}

// Peek returns the next n runes from the buffer without advancing the reader.
// The runes stop being valid at the next read call. If Peek returns fewer than
// n runes, it also returns an error indicating why the read is short.
// ErrBufferFull is returned if n is larger than the reader's buffer size.
//
// Calling Peek prevents an UnreadRune call from succeeding until the next
// read operation.
func (r *RuneReader) Peek(n int) ([]rune, error) {
	if n < 0 {
		return nil, ErrNegativeCount
	}

	if n > len(r.buf) {
		return nil, ErrBufferFull
	}

	r.lastRune = -1

	if n > r.Buffered() {
		r.fill(n)
	}

	if n > r.Buffered() {
		n = r.Buffered()
	}

	return r.buf[r.r : r.r+n], r.readErr()
}

// Reset discards any buffered data, resets all state, and switches the
// buffered reader to read from rd. Calling Reset on the zero value of Reader
// initializes the internal buffer to the default size. Calling r.Reset(r) (that
// is, resetting a Reader to itself) does nothing.
func (r *RuneReader) Reset(rd io.RuneReader) {
	if r == rd {
		return
	}
	if r.buf == nil {
		r.buf = make([]rune, defaultBufSize)
	}
	r.reset(r.buf, rd)
}

func (r *RuneReader) reset(buf []rune, rd io.RuneReader) {
	*r = RuneReader{
		rd:       rd,
		buf:      buf,
		lastRune: -1,
	}
}

// Size returns the size of the underlying buffer in number of runes.
func (r *RuneReader) Size() int {
	return len(r.buf)
}

// UnreadRune unreads the last rune. Only the most recently read rune can be unread.
//
// UnreadRune returns ErrInvalidUnreadRune if the most recent method called on
// the RuneReader was not a read operation. Notably, Peek, and Discard are not
// considered read operations.
func (r *RuneReader) UnreadRune() error {
	if r.lastRune < 0 {
		return ErrInvalidUnreadRune
	}
	r.r--
	r.buf[r.r] = r.lastRune
	r.lastRune = -1
	return nil
}

// Buffered returns the number of runes that can be read from the buffer.
func (r *RuneReader) Buffered() int {
	return r.e - r.r
}

// fill fills the RuneReader's buffer so that it contains n runes.
func (r *RuneReader) fill(n int) {
	if r.r > 0 && r.e-r.r < n {
		copy(r.buf, r.buf[r.r:r.e])
		r.e -= r.r
		r.r = 0
	}

	// Fill the rest of the buffer.
	for ; r.e < n; r.e++ {
		rn, _, err := r.rd.ReadRune()
		if err != nil {
			r.err = err
			break
		}
		r.buf[r.e] = rn
	}
}

// readErr returns the last readErr and clears it.
func (r *RuneReader) readErr() error {
	err := r.err
	r.err = nil
	return err
}
