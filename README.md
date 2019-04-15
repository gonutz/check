# Usage

This is a copy of the [godoc](https://godoc.org/github.com/gonutz/check) for this package.

`func Eq(t Tester, a, b interface{}, msg ...interface{})`

Eq compares a and b and calls Errorf on t if they differ. Values are compared in a deep way, similar to reflect.DeepEqual, only that float and complex values are compared using an epsilon of 1e-6. If there are any msg parameters, they are printed in concatenation before the error message, e.g. if you pass ["input ", 5] as msg, errors will be printed as: "input 5: <error>".


`func EqEps(t Tester, a, b interface{}, epsilon float64, msg ...interface{})`

EqEps compares a and b and calls Errorf on t if they differ. Values are compared in a deep way, similar to reflect.DeepEqual, only that float and complex values are compared using epsilon. Values are considered equal if their absolute difference is less than or equal to epsilon. Set epsilon to zero to compare for exact equality (or use EqExact). If there are any msg parameters, they are printed in concatenation before the error message, e.g. if you pass ["input ", 5] as msg, errors will be printed as: "input 5: <error>".


`func EqExact(t Tester, a, b interface{}, msg ...interface{})`

EqExact compares a and b and calls Errorf on t if they differ. Values are compared in a deep way, similar to reflect.DeepEqual, float and complex values must match exactly. If there are any msg parameters, they are printed in concatenation before the error message, e.g. if you pass ["input ", 5] as msg, errors will be printed as: "input 5: <error>".


`func Neq(t Tester, a, b interface{}, msg ...interface{})`

Neq compares a and b and calls Errorf on t if they are equal. Values are compared in a deep way, similar to reflect.DeepEqual, only that float and complex values are compared using an epsilon of 1e-6. If there are any msg parameters, they are printed in concatenation before the error message, e.g. if you pass ["input ", 5] as msg, errors will be printed as: "input 5: <error>".


`func NeqEps(t Tester, a, b interface{}, epsilon float64, msg ...interface{})`

NeqEps compares a and b and calls Errorf on t if they are equal. Values are compared in a deep way, similar to reflect.DeepEqual, only that float and complex values are compared using epsilon. Values are considered equal if their absolute difference is less than or equal to epsilon. Set epsilon to zero to compare for exact equality (or use EqExact). If there are any msg parameters, they are printed in concatenation before the error message, e.g. if you pass ["input ", 5] as msg, errors will be printed as: "input 5: <error>".


`func NeqExact(t Tester, a, b interface{}, msg ...interface{})`

NeqExact compares a and b and calls Errorf on t if they are equal. Values are compared in a deep way, similar to reflect.DeepEqual, float and complex values must match exactly. If there are any msg parameters, they are printed in concatenation before the error message, e.g. if you pass ["input ", 5] as msg, errors will be printed as: "input 5: <error>".


Use your `*testing.T` for the `Tester` parameter.


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