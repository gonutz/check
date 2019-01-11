package check

import (
	"math"
	"reflect"
	"testing"
	"unsafe"
)

// Eq compares a and b and calls Errorf on t if they differ. Values are compared
// in a deep way, similar to reflect.DeepEqual, only that float and complex
// values are compared using an epsilon of 1e-6.
func Eq(t *testing.T, a, b interface{}) {
	t.Helper()
	EqEps(t, a, b, 1e-6)
}

// EqExact compares a and b and calls Errorf on t if they differ. Values are
// compared in a deep way, similar to reflect.DeepEqual, float and complex
// values must match exactly.
func EqExact(t *testing.T, a, b interface{}) {
	t.Helper()
	EqEps(t, a, b, 0)
}

// EqEps compares a and b and calls Errorf on t if they differ. Values are
// compared in a deep way, similar to reflect.DeepEqual, only that float and
// complex values are compared using epsilon. Values are considered equal if
// their absolute difference is less than or equal to epsilon. Set epsilon to
// zero to compare for exact equality (or use EqExact).
func EqEps(t *testing.T, a, b interface{}, epsilon float64) {
	t.Helper()
	if !deepEqual(a, b, epsilon) {
		t.Errorf("%#v != %#v", a, b)
	}
}

// deepEqual is a modified version of reflect.DeepEqual. deepEqual compares
// float and complex values using epsilon.
func deepEqual(x, y interface{}, epsilon float64) bool {
	if x == nil || y == nil {
		return x == y
	}
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		return false
	}
	return deepValueEqual(v1, v2, epsilon, make(map[visit]bool), 0)
}

type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

func deepValueEqual(v1, v2 reflect.Value, eps float64, visited map[visit]bool, depth int) bool {
	if !v1.IsValid() || !v2.IsValid() {
		return v1.IsValid() == v2.IsValid()
	}
	if v1.Type() != v2.Type() {
		return false
	}

	// We want to avoid putting more in the visited map than we need to.
	// For any possible reference cycle that might be encountered,
	// hard(t) needs to return true for at least one of the types in the cycle.
	hard := func(k reflect.Kind) bool {
		switch k {
		case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
			return true
		}
		return false
	}

	if v1.CanAddr() && v2.CanAddr() && hard(v1.Kind()) {
		addr1 := unsafe.Pointer(v1.UnsafeAddr())
		addr2 := unsafe.Pointer(v2.UnsafeAddr())
		if uintptr(addr1) > uintptr(addr2) {
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true
		}

		// Remember for later.
		visited[v] = true
	}

	switch v1.Kind() {
	case reflect.Array:
		for i := 0; i < v1.Len(); i++ {
			if !deepValueEqual(v1.Index(i), v2.Index(i), eps, visited, depth+1) {
				return false
			}
		}
		return true
	case reflect.Slice:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for i := 0; i < v1.Len(); i++ {
			if !deepValueEqual(v1.Index(i), v2.Index(i), eps, visited, depth+1) {
				return false
			}
		}
		return true
	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() == v2.IsNil()
		}
		return deepValueEqual(v1.Elem(), v2.Elem(), eps, visited, depth+1)
	case reflect.Ptr:
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		return deepValueEqual(v1.Elem(), v2.Elem(), eps, visited, depth+1)
	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			if !deepValueEqual(v1.Field(i), v2.Field(i), eps, visited, depth+1) {
				return false
			}
		}
		return true
	case reflect.Map:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for _, k := range v1.MapKeys() {
			val1 := v1.MapIndex(k)
			val2 := v2.MapIndex(k)
			if !val1.IsValid() ||
				!val2.IsValid() ||
				!deepValueEqual(v1.MapIndex(k), v2.MapIndex(k), eps, visited, depth+1) {
				return false
			}
		}
		return true
	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		// Can't do better than this:
		return false
	case reflect.Bool:
		return v2.Kind() == reflect.Bool && v1.Bool() == v2.Bool()
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		k2 := v2.Kind()
		return (k2 == reflect.Uint ||
			k2 == reflect.Uint8 ||
			k2 == reflect.Uint16 ||
			k2 == reflect.Uint32 ||
			k2 == reflect.Uint64 ||
			k2 == reflect.Uintptr) && v1.Uint() == v2.Uint()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		k2 := v2.Kind()
		return (k2 == reflect.Int ||
			k2 == reflect.Int8 ||
			k2 == reflect.Int16 ||
			k2 == reflect.Int32 ||
			k2 == reflect.Int64) && v1.Int() == v2.Int()
	case reflect.Float32, reflect.Float64:
		k2 := v2.Kind()
		return (k2 == reflect.Float32 || k2 == reflect.Float64) &&
			floatEq(v1.Float(), v2.Float(), eps)
	case reflect.Complex64, reflect.Complex128:
		k2 := v2.Kind()
		return (k2 == reflect.Complex64 || k2 == reflect.Complex128) &&
			floatEq(real(v1.Complex()), real(v2.Complex()), eps) &&
			floatEq(imag(v1.Complex()), imag(v2.Complex()), eps)
	case reflect.String:
		return v2.Kind() == reflect.String && v1.String() == v2.String()
	case reflect.UnsafePointer:
		return v2.Kind() == reflect.UnsafePointer &&
			v1.UnsafeAddr() == v2.UnsafeAddr()
	default:
		return false
	}
}

func floatEq(a, b, eps float64) bool {
	return math.IsInf(a, 1) && math.IsInf(b, 1) ||
		math.IsInf(a, -1) && math.IsInf(b, -1) ||
		math.IsNaN(a) && math.IsNaN(b) ||
		abs(a-b) <= eps
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
