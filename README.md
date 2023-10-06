# runeio

[![GoDoc](https://godoc.org/github.com/ianlewis/runeio?status.svg)](https://godoc.org/github.com/ianlewis/runeio)
[![Go Report Card](https://goreportcard.com/badge/github.com/ianlewis/runeio)](https://goreportcard.com/report/github.com/ianlewis/runeio)
[![tests](https://github.com/ianlewis/runeio/actions/workflows/pre-submit.units.yml/badge.svg)](https://github.com/ianlewis/runeio/actions/workflows/pre-submit.units.yml)
[![codecov](https://codecov.io/gh/ianlewis/runeio/graph/badge.svg?token=H2VXJL5MEV)](https://codecov.io/gh/ianlewis/runeio)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/ianlewis/runeio/badge)](https://api.securityscorecards.dev/projects/github.com/ianlewis/runeio)

`runeio` is a Go package that provides basic interfaces to I/O primitives for
runes. Its primary job is to wrap existing implementations of such primitives,
such as those in package io.

Because these interfaces and primitives wrap lower-level operations with various
implementations, unless otherwise informed clients should not assume they are
safe for parallel execution.
