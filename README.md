# Usage

`func Eq(t *testing.T, a, b interface{})`

Eq compares a and b and calls Errorf on t if they differ. Values are compared in a deep way, similar to reflect.DeepEqual, only that float and complex values are compared using an epsilon of 1e-6. 


`func EqEps(t *testing.T, a, b interface{}, epsilon float64)`

EqEps compares a and b and calls Errorf on t if they differ. Values are compared in a deep way, similar to reflect.DeepEqual, only that float and complex values are compared using epsilon. Values are considered equal if their absolute difference is less than or equal to epsilon. Set epsilon to zero to compare for exact equality (or use EqExact).


`func EqExact(t *testing.T, a, b interface{})`

EqExact compares a and b and calls Errorf on t if they differ. Values are compared in a deep way, similar to reflect.DeepEqual, float and complex values must match exactly.


# Installation

Run `go get github.com/gonutz/check` to install this library.


# Note

If the name `check` conflicts with something you are already using, you might rename your import to something else, like `is` or `must`:

```
package main

import (
	"testing"
	must "github.com/gonutz/check"
)

func Test(t *testing.T) {
	must.Eq(t, 1, 1)
}
```