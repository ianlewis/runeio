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
	"strings"
	"testing"
	"unicode"
)

// sliceEqual returns true if the two slices are equal.
func sliceEqual(l, r []rune) bool {
	if len(l) != len(r) {
		return false
	}

	for i := range l {
		if l[i] != r[i] {
			return false
		}
	}

	return true
}

type expectation interface {
	expect(*testing.T, *RuneReader)
}

type expectedPeek struct {
	size          int
	expectedRunes []rune
	expectedErr   error
}

func (e *expectedPeek) expect(t *testing.T, r *RuneReader) {
	t.Helper()
	p, err := r.Peek(e.size)
	if got, want := err, e.expectedErr; !errors.Is(got, want) {
		t.Errorf("expected error: got: %v, want: %v", got, want)
	}

	if got, want := p, e.expectedRunes; !sliceEqual(got, want) {
		t.Errorf("Peek: got: %q, want: %q", string(got), string(want))
	}
}

type expectedDiscard struct {
	size        int
	expected    int
	expectedErr error
}

func (e *expectedDiscard) expect(t *testing.T, r *RuneReader) {
	t.Helper()
	n, err := r.Discard(e.size)
	if got, want := err, e.expectedErr; !errors.Is(got, want) {
		t.Errorf("expected error: got: %v, want: %v", got, want)
	}

	if got, want := n, e.expected; got != want {
		t.Errorf("Discard: got: %d, want: %d", got, want)
	}
}

type expectedRead struct {
	size          int
	expectedNum   int
	expectedErr   error
	expectedRunes []rune
}

func (e *expectedRead) expect(t *testing.T, r *RuneReader) {
	t.Helper()
	p := make([]rune, e.size)
	n, err := r.Read(p)
	if got, want := err, e.expectedErr; !errors.Is(got, want) {
		t.Errorf("expected error: got: %v, want: %v", got, want)
	}

	if got, want := n, e.expectedNum; got != want {
		t.Errorf("Read: got: %v, want: %v", got, want)
	}
	if got, want := p[:n], e.expectedRunes; !sliceEqual(got, want) {
		t.Errorf("Read: got: %q, want: %q", string(got), string(want))
	}
}

type expectedReset struct {
	str       string
	selfReset bool
}

func (e *expectedReset) expect(t *testing.T, r *RuneReader) {
	t.Helper()
	if e.selfReset {
		r.Reset(r)
	} else {
		r.Reset(strings.NewReader(e.str))
	}
}

type expectedReadRune struct {
	expectedRune rune
	expectedSize int
	expectedErr  error
}

func (e *expectedReadRune) expect(t *testing.T, r *RuneReader) {
	t.Helper()
	rn, size, err := r.ReadRune()
	if got, want := err, e.expectedErr; !errors.Is(got, want) {
		t.Errorf("expected error: got: %v, want: %v", got, want)
	}

	if got, want := size, e.expectedSize; got != want {
		t.Errorf("ReadRune size: got: %v, want: %v", got, want)
	}
	if got, want := rn, e.expectedRune; got != want {
		t.Errorf("ReadRune rune: got: %v, want: %v", got, want)
	}
}

type expectedSize struct {
	expectedSize int
}

func (e *expectedSize) expect(t *testing.T, r *RuneReader) {
	t.Helper()
	if got, want := r.Size(), e.expectedSize; got != want {
		t.Errorf("Size: got: %v, want: %v", got, want)
	}
}

type expectedUnreadRune struct {
	expectedErr error
}

func (e *expectedUnreadRune) expect(t *testing.T, r *RuneReader) {
	t.Helper()
	err := r.UnreadRune()
	if got, want := err, e.expectedErr; !errors.Is(got, want) {
		t.Errorf("expected error: got: %v, want: %v", got, want)
	}
}

func TestRuneReader(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		str      string
		bufSize  int
		expected []expectation
	}{
		{
			name: "single read",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedRead{
					size:          9,
					expectedNum:   9,
					expectedRunes: []rune("Hello, 世界"),
				},
			},
		},
		{
			name: "single read not exact",
			str:  "Hello, 世界",
			expected: []expectation{
				&expectedRead{
					size:          12,
					expectedNum:   9,
					expectedRunes: []rune("Hello, 世界"),
					expectedErr:   io.EOF,
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "multiple reads exact",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedRead{
					size:          3,
					expectedNum:   3,
					expectedRunes: []rune("Hel"),
				},
				&expectedRead{
					size:          5,
					expectedNum:   5,
					expectedRunes: []rune("lo, 世"),
				},
				&expectedRead{
					size:          1,
					expectedNum:   1,
					expectedRunes: []rune("界"),
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "multiple reads not exact",
			str:  "Hello, 世界",
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedRead{
					size:          3,
					expectedNum:   3,
					expectedRunes: []rune("Hel"),
				},
				&expectedRead{
					size:          5,
					expectedNum:   5,
					expectedRunes: []rune("lo, 世"),
				},
				&expectedRead{
					size:          8,
					expectedNum:   1,
					expectedRunes: []rune("界"),
					expectedErr:   io.EOF,
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "single peek",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedPeek{
					size:          3,
					expectedRunes: []rune("Hel"),
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "multiple peek",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedPeek{
					size:          3,
					expectedRunes: []rune("Hel"),
				},
				&expectedPeek{
					size:          5,
					expectedRunes: []rune("Hello"),
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "read and peek",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedRead{
					size:          3,
					expectedNum:   3,
					expectedRunes: []rune("Hel"),
				},
				&expectedPeek{
					size:          5,
					expectedRunes: []rune("lo, 世"),
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "read and peek multi",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedRead{
					size:          3,
					expectedNum:   3,
					expectedRunes: []rune("Hel"),
				},
				&expectedPeek{
					size:          5,
					expectedRunes: []rune("lo, 世"),
				},
				&expectedRead{
					size:          4,
					expectedNum:   4,
					expectedRunes: []rune("lo, "),
				},
				&expectedPeek{
					size:          1,
					expectedRunes: []rune("世"),
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "peek exact",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedPeek{
					size:          9,
					expectedRunes: []rune("Hello, 世界"),
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "peek not exact",
			str:  "Hello, 世界",
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedPeek{
					size:          11,
					expectedRunes: []rune("Hello, 世界"),
					expectedErr:   io.EOF,
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "peek neg count",
			str:  "Hello, 世界",
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedPeek{
					size:        -1,
					expectedErr: ErrNegativeCount,
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "discard exact",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedDiscard{
					size:     9,
					expected: 9,
				},
				&expectedPeek{
					size:          6,
					expectedRunes: []rune("/World"),
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "discard not exact",
			str:  "Hello, 世界",
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedDiscard{
					size:        11,
					expected:    9,
					expectedErr: io.EOF,
				},
				&expectedPeek{
					size:        5,
					expectedErr: io.EOF,
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "discard neg count",
			str:  "Hello, 世界",
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedDiscard{
					size:        -1,
					expected:    0,
					expectedErr: ErrNegativeCount,
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name:    "peek larger than buffer size",
			str:     "Hello, 世界",
			bufSize: 5,
			expected: []expectation{
				&expectedSize{
					expectedSize: 5,
				},
				&expectedPeek{
					size:        6,
					expectedErr: ErrBufferFull,
				},
				&expectedSize{
					expectedSize: 5,
				},
			},
		},
		{
			name: "utf-8 replacement character",
			str:  string([]byte{0xef, 0xbf, 0xbd}),
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedRead{
					size:          1,
					expectedNum:   1,
					expectedRunes: []rune{unicode.ReplacementChar},
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "utf-8 two continuation bytes",
			str:  string([]byte{0x80, 0x80}),
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedRead{
					size:          1,
					expectedNum:   1,
					expectedRunes: []rune{unicode.ReplacementChar},
				},
				&expectedRead{
					size:          1,
					expectedNum:   1,
					expectedRunes: []rune{unicode.ReplacementChar},
				},
				&expectedRead{
					size:        1,
					expectedErr: io.EOF,
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "utf-8 invalid bytes",
			str:  string([]byte{0xfe, 0xfe, 0xff, 0xff}),
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedRead{
					size:          1,
					expectedNum:   1,
					expectedRunes: []rune{unicode.ReplacementChar},
				},
				&expectedRead{
					size:          2,
					expectedNum:   2,
					expectedRunes: []rune{unicode.ReplacementChar, unicode.ReplacementChar},
				},
				&expectedRead{
					size:          2,
					expectedNum:   1,
					expectedRunes: []rune{unicode.ReplacementChar},
					expectedErr:   io.EOF,
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "utf-8 valid and invalid mixed",
			str:  string([]byte{0xfe, 0xfe, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c, 0xff, 0xff}),
			expected: []expectation{
				&expectedSize{
					expectedSize: defaultBufSize,
				},
				&expectedRead{
					size:          2,
					expectedNum:   2,
					expectedRunes: []rune{unicode.ReplacementChar, unicode.ReplacementChar},
				},
				&expectedRead{
					size:          2,
					expectedNum:   2,
					expectedRunes: []rune("世界"),
				},
				&expectedRead{
					size:          3,
					expectedNum:   2,
					expectedRunes: []rune{unicode.ReplacementChar, unicode.ReplacementChar},
					expectedErr:   io.EOF,
				},
				&expectedSize{
					expectedSize: defaultBufSize,
				},
			},
		},
		{
			name: "unread no read",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedUnreadRune{
					expectedErr: ErrInvalidUnreadRune,
				},
			},
		},
		{
			name: "single unread",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedRead{
					size:          9,
					expectedNum:   9,
					expectedRunes: []rune("Hello, 世界"),
				},
				&expectedUnreadRune{},
				&expectedPeek{
					size:          1,
					expectedRunes: []rune("界"),
				},
			},
		},
		{
			name: "multi-unread",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedRead{
					size:          9,
					expectedNum:   9,
					expectedRunes: []rune("Hello, 世界"),
				},
				&expectedUnreadRune{},
				&expectedUnreadRune{
					expectedErr: ErrInvalidUnreadRune,
				},
				&expectedPeek{
					size:          1,
					expectedRunes: []rune("界"),
				},
			},
		},
		{
			name: "peek-unread",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedRead{
					size:          9,
					expectedNum:   9,
					expectedRunes: []rune("Hello, 世界"),
				},
				&expectedPeek{
					size:          3,
					expectedRunes: []rune("/Wo"),
				},
				&expectedUnreadRune{
					expectedErr: ErrInvalidUnreadRune,
				},
				&expectedPeek{
					size:          3,
					expectedRunes: []rune("/Wo"),
				},
			},
		},
		{
			name: "discard-unread",
			str:  "Hello, 世界/World/Universe!",
			expected: []expectation{
				&expectedRead{
					size:          9,
					expectedNum:   9,
					expectedRunes: []rune("Hello, 世界"),
				},
				&expectedDiscard{
					size:     3,
					expected: 3,
				},
				&expectedUnreadRune{
					expectedErr: ErrInvalidUnreadRune,
				},
				&expectedPeek{
					size:          3,
					expectedRunes: []rune("rld"),
				},
			},
		},
		{
			name: "readrune",
			str:  "Hello, 世界!",
			expected: []expectation{
				&expectedReadRune{
					expectedRune: 'H',
					expectedSize: 1,
				},
				&expectedReadRune{
					expectedRune: 'e',
					expectedSize: 1,
				},
				&expectedReadRune{
					expectedRune: 'l',
					expectedSize: 1,
				},
				&expectedReadRune{
					expectedRune: 'l',
					expectedSize: 1,
				},
				&expectedReadRune{
					expectedRune: 'o',
					expectedSize: 1,
				},
				&expectedReadRune{
					expectedRune: ',',
					expectedSize: 1,
				},
				&expectedReadRune{
					expectedRune: ' ',
					expectedSize: 1,
				},
				&expectedReadRune{
					expectedRune: '世',
					expectedSize: 3,
				},
				&expectedReadRune{
					expectedRune: '界',
					expectedSize: 3,
				},
				&expectedReadRune{
					expectedRune: '!',
					expectedSize: 1,
				},
				&expectedReadRune{
					expectedErr: io.EOF,
				},
			},
		},
		{
			name: "reset",
			str:  "Hello, 世界!",
			expected: []expectation{
				&expectedRead{
					size:          9,
					expectedNum:   9,
					expectedRunes: []rune("Hello, 世界"),
				},
				&expectedReset{
					str: "Hello World!",
				},
				&expectedRead{
					size:          9,
					expectedNum:   9,
					expectedRunes: []rune("Hello Wor"),
				},
			},
		},
		{
			name: "reset-noop",
			str:  "Hello, 世界!",
			expected: []expectation{
				&expectedRead{
					size:          7,
					expectedNum:   7,
					expectedRunes: []rune("Hello, "),
				},
				&expectedReset{
					selfReset: true,
				},
				&expectedRead{
					size:          3,
					expectedNum:   3,
					expectedRunes: []rune("世界!"),
				},
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		b := strings.NewReader(c.str)

		var r *RuneReader
		if c.bufSize != 0 {
			r = NewReaderSize(b, c.bufSize)
		} else {
			r = NewReader(b)
		}

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			for _, e := range c.expected {
				e.expect(t, r)
			}
		})
	}
}

func TestRuneReader_Reset_ZeroValue(t *testing.T) {
	t.Parallel()

	r := new(RuneReader)
	r.Reset(strings.NewReader("foo"))
	b, err := r.Peek(2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got, want := string(b), "fo"; got != want {
		t.Errorf("Peek: got: %q, want: %q", got, want)
	}
}

func BenchmarkReadRune(b *testing.B) {
	s := strings.Repeat("x", 512+1)
	rs := strings.NewReader(s)
	rr := NewReader(rs)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//nolint:errcheck // errors not checked in benchmarks.
		_, _, _ = rr.ReadRune()
		rs.Reset(s)
		rr.Reset(rs)
	}
}

func BenchmarkReadSmall(b *testing.B) {
	s := strings.Repeat("x", 512+1)
	rs := strings.NewReader(s)
	rr := NewReader(rs)
	buf := make([]rune, 512)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//nolint:errcheck // errors not checked in benchmarks.
		_, _ = rr.Read(buf)
		rs.Reset(s)
		rr.Reset(rs)
	}
}

func BenchmarkReadLarge(b *testing.B) {
	s := strings.Repeat("x", (32*1024)+1)
	rs := strings.NewReader(s)
	rr := NewReader(rs)
	buf := make([]rune, 32*1024)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//nolint:errcheck // errors not checked in benchmarks.
		_, _ = rr.Read(buf)
		rs.Reset(s)
		rr.Reset(rs)
	}
}

func BenchmarkNoCopySmall(b *testing.B) {
	n := 512
	s := strings.Repeat("x", n+1)
	rs := strings.NewReader(s)
	rr := NewReader(rs)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var read int
		for read < n {
			n := rr.Buffered()
			if n == 0 {
				n = 1
			}
			//nolint:errcheck // errors not checked in benchmarks.
			buf, _ := rr.Peek(n)
			p := len(buf)
			//nolint:errcheck // errors not checked in benchmarks.
			_, _ = rr.Discard(p)
			read += p
		}

		rs.Reset(s)
		rr.Reset(rs)
	}
}

func BenchmarkNoCopyLarge(b *testing.B) {
	n := 32 * 1024
	s := strings.Repeat("x", n+1)
	rs := strings.NewReader(s)
	rr := NewReader(rs)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var read int
		for read < n {
			n := rr.Buffered()
			if n == 0 {
				n = 1
			}
			//nolint:errcheck // errors not checked in benchmarks.
			buf, _ := rr.Peek(n)
			p := len(buf)
			//nolint:errcheck // errors not checked in benchmarks.
			_, _ = rr.Discard(p)
			read += p
		}

		rs.Reset(s)
		rr.Reset(rs)
	}
}

func BenchmarkPeekSmall(b *testing.B) {
	n := 512
	s := strings.Repeat("x", n+1)
	rs := strings.NewReader(s)
	rr := NewReader(rs)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//nolint:errcheck // errors not checked in benchmarks.
		_, _ = rr.Peek(n)
		rs.Reset(s)
		rr.Reset(rs)
	}
}

func BenchmarkPeekLarge(b *testing.B) {
	n := 32 * 1024
	s := strings.Repeat("x", n+1)
	rs := strings.NewReader(s)
	rr := NewReaderSize(rs, n+1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//nolint:errcheck // errors not checked in benchmarks.
		_, _ = rr.Peek(n)
		rs.Reset(s)
		rr.Reset(rs)
	}
}

func BenchmarkDiscardSmall(b *testing.B) {
	n := 512
	s := strings.Repeat("x", n+1)
	rs := strings.NewReader(s)
	rr := NewReader(rs)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//nolint:errcheck // errors not checked in benchmarks.
		_, _ = rr.Discard(n)
		rs.Reset(s)
		rr.Reset(rs)
	}
}

func BenchmarkDiscardLarge(b *testing.B) {
	n := 32 * 1024
	s := strings.Repeat("x", n+1)
	rs := strings.NewReader(s)
	rr := NewReaderSize(rs, n+1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//nolint:errcheck // errors not checked in benchmarks.
		_, _ = rr.Discard(n)
		rs.Reset(s)
		rr.Reset(rs)
	}
}
