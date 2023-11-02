# runeio

[![GoDoc](https://godoc.org/github.com/ianlewis/runeio?status.svg)](https://godoc.org/github.com/ianlewis/runeio)
[![Go Report Card](https://goreportcard.com/badge/github.com/ianlewis/runeio)](https://goreportcard.com/report/github.com/ianlewis/runeio)
[![tests](https://github.com/ianlewis/runeio/actions/workflows/pre-submit.units.yml/badge.svg)](https://github.com/ianlewis/runeio/actions/workflows/pre-submit.units.yml)
[![codecov](https://codecov.io/gh/ianlewis/runeio/graph/badge.svg?token=H2VXJL5MEV)](https://codecov.io/gh/ianlewis/runeio)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/ianlewis/runeio/badge)](https://api.securityscorecards.dev/projects/github.com/ianlewis/runeio)

`runeio` is a Go package that provides basic interfaces to I/O primitives for
runes. Its primary job is to wrap and supplement existing implementations of
such primitives, such as those in package `io`.

Because these interfaces and primitives wrap lower-level operations with various
implementations, unless otherwise informed clients should not assume they are
safe for parallel execution.

The package provides a buffered `RuneReader` that implements the
[`io.RuneScanner`](https://pkg.go.dev/io#RuneScanner) interface. Its primary
purpose is to buffer the reading of runes and allow users to `Peek` ahead by a
number of runes rather than a number of bytes as with the `bufio.Reader`.

## Examples

A single rune can be read with `ReadRune`.

```golang
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
```

`Read` can be used to read runes into a buffer.

```golang
r := runeio.NewReader(strings.NewReader("Hello World!"))

buf := make([]rune, 5)
_, _ = r.Read(buf)

fmt.Print(string(buf))

// Output: Hello
```

`Peek` can be used to look ahead into the stream without consuming the runes.

```golang
r := runeio.NewReader(strings.NewReader("Hello World!"))

buf := make([]rune, 6)
_, _ = r.Read(buf)
peeked, _ := r.Peek(6)

fmt.Print(string(peeked))

// Output: World!
```

## Related Projects

- [pelletier/go-buffruneio](https://github.com/pelletier/go-buffruneio): Another
  buffering implementation of `io.RuneScanner`.
  - Uses a custom `badRune` value for encoding errors instead of
    `unicode.ReplacementChar`
  - Uses a custom `EOF` rune value instead of just returning `io.EOF`.
  - The undo buffer grows unbounded until `Forget` is called.
  - `NewReader` also takes a `io.Reader` which is wrapped in a `bufio.Reader`
    when it could just take an `io.RuneReader` since it the implementation only
    requires`ReadRune` from the underlying reader.
- [SteelSeries/bufrr](https://github.com/SteelSeries/bufrr): Another buffering
  implementation of `io.RuneScanner`.
  - Implements only a static single rune buffer. This means that only one rune
    could be peeked at.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md)
