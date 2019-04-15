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


# Rationale

Package check implements easy to use functions to write your tests in a concise
way.

In the past I have used if-statements to compare values for equality, for
example:

```
func Test(t *testing.T) {
	sum := add(1, 2)
	if sum != 3 {
		t.Errorf("1+2 = 3 but got %v", sum)
	}
}
```

This would soon get tedious and copy paste errors for the error messages would
creep in. Thus I would create a helper function to compare two integers, like
this:

```
func checkInts(t *testing.T, msg string, have, want int) {
	if have != want {
		t.Errorf("%s: %v != %v", msg, have, want)
	}
}
```

Since there are no generics in Go, I would have to write a helper for all types
of values that I wanted to compare: int, uint, byte, float32, []byte, map, etc.
The functions to implement comparison for equality of two values can get more
complex depeding on the type.
Floating point values can usually not be compared with == because of rounding
errors. You want to have an epsilon by which the values might diverge from each
other and still be considered equal. Furthermore, what about INF and NAN?
Comparing slices and maps means you must compare their lengths and iterate over
them to check each item for equality.
Nested structs and interfaces can get hairy pretty quickly.

To get over these problems once and for all I created this package. It aims at a
minimal API with maximum usability. You can only check for equality or
non-equality with the Eq and Neq functions.

The above example becomes:

```
func Test(t *testing.T) {
	sum := add(1, 2)
	check.Eq(t, sum, 3, "1+2")
}
```

It does not matter whether the add function returns an int, a uint32, a byte or
a float64. Eq and Neq compare values in a deep way while handling different
integer types, floating point accuracy, INF and NAN and comparison between
string, []byte and []rune.

This package will not solve all your testing needs but probably 95% of it. You
can still write if-statements or special helpers for the cases where simple
equality of values does not fit your needs.


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