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

type expectedSize struct {
	expectedSize int
}

func (e *expectedSize) expect(t *testing.T, r *RuneReader) {
	t.Helper()
	if got, want := r.Size(), e.expectedSize; got != want {
		t.Errorf("Size: got: %v, want: %v", got, want)
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
