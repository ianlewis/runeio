# runeio

[![GoDoc](https://godoc.org/github.com/ianlewis/runeio?status.svg)](https://godoc.org/github.com/ianlewis/runeio/runeio)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)

`runeio` is a Go package that provides basic interfaces to I/O primitives for
runes. Its primary job is to wrap existing implementations of such primitives,
such as those in package io.

Because these interfaces and primitives wrap lower-level operations with various
implementations, unless otherwise informed clients should not assume they are
safe for parallel execution.
