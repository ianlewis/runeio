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

package runeio_test

import (
	"fmt"
	"strings"

	"github.com/ianlewis/runeio"
)

func ExampleRuneReader() {
	r := runeio.NewReader(strings.NewReader("Hello World!"))

	var runes []rune
	for {
		rn, _, err := r.ReadRune()
		if err != nil {
			break
		}
		runes = append(runes, rn)
	}

	fmt.Print(string(runes))

	// Output: Hello World!
}

func ExampleRuneReader_Read() {
	r := runeio.NewReader(strings.NewReader("Hello World!"))

	buf := make([]rune, 5)
	_, _ = r.Read(buf)

	fmt.Print(string(buf))

	// Output: Hello
}

func ExampleRuneReader_Peek() {
	r := runeio.NewReader(strings.NewReader("Hello World!"))

	buf := make([]rune, 6)
	_, _ = r.Read(buf)
	peeked, _ := r.Peek(6)

	fmt.Print(string(peeked))

	// Output: World!
}
