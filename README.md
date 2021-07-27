# safeulid
[![Test](https://github.com/vvvvv/safeulid/actions/workflows/test.yml/badge.svg)](https://github.com/vvvvv/safeulid/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/vvvvv/safeulid/branch/main/graph/badge.svg)](https://codecov.io/gh/vvvvv/safeulid) [![Go Reference](https://pkg.go.dev/badge/github.com/vvvvv/safeulid.svg)](https://pkg.go.dev/github.com/vvvvv/safeulid)

Small wrapper lib around [github.com/oklog/ulid/v2](https://github.com/oklog/ulid) that provides concurrency safety

## example
```
package main

import (
	"fmt"
	"sync"

	ulid "github.com/vvvvv/safeulid"
)

func main() {
	id := ulid.MustNew()
	fmt.Println(id)
}
```
