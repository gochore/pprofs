## pprofs

Auto capture profiles.

[![Build Status](https://travis-ci.com/gochore/pprofs.svg?branch=master)](https://travis-ci.com/gochore/pprofs)
[![codecov](https://codecov.io/gh/gochore/pprofs/branch/master/graph/badge.svg)](https://codecov.io/gh/gochore/pprofs)
[![Go Report Card](https://goreportcard.com/badge/github.com/gochore/pprofs)](https://goreportcard.com/report/github.com/gochore/pprofs)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gochore/pprofs)](https://github.com/gochore/pprofs/blob/master/go.mod)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/gochore/pprofs)](https://github.com/gochore/pprofs/releases)

## Quick start

```go
package main

import "github.com/gochore/pprofs"

func main() {
	if err := pprofs.EnableCapture(); err != nil {
		panic(err)
	}
	// ...
}
```

It will auto capture profiles and store files in temp dir.

You can specify the dir and the ttl of files by environment variables:

```bash
export PPROF_DIR="~/somewhere"
export PPROF_TTL="1h"
```

Or use options:

```go
func main() {
	if err := pprofs.EnableCapture(
		pprofs.WithStorage(pprofs.NewFileStorage("prefix", "~/somewhere", time.Hour)),
	); err != nil {
		panic(err)
	}
	// ...
}
```

See [more examples](https://github.com/gochore/pprofs/tree/master/_example).

## Reference

- Inspired by [autopprof](https://github.com/rakyll/autopprof).
- Learned a lot from [Continuous Profiling of Go programs](https://medium.com/google-cloud/continuous-profiling-of-go-programs-96d4416af77b).