package check

import (
	"math"
	"reflect"
	"testing"
)

var Eps = 1e-6

func Eq(t *testing.T, a, b interface{}) {
	t.Helper()
	eq := reflect.DeepEqual(a, b)
	if !eq {
		ua, uaOK := toUint64(a)
		ub, ubOK := toUint64(b)
		if uaOK && ubOK {
			eq = ua == ub
		} else {
			fa, faOK := toFloat64(a)
			fb, fbOK := toFloat64(b)
			if faOK && fbOK {
				eq = math.IsInf(fa, 1) && math.IsInf(fb, 1) ||
					math.IsInf(fa, -1) && math.IsInf(fb, -1) ||
					math.IsNaN(fa) && math.IsNaN(fb) ||
					abs(fa-fb) < Eps
			}
		}
	}
	if !eq {
		t.Errorf("%v != %v", a, b)
	}
}

func toUint64(x interface{}) (uint64, bool) {
	switch x.(type) {
	case int:
		return uint64(x.(int)), true
	case int8:
		return uint64(x.(int8)), true
	case int16:
		return uint64(x.(int16)), true
	case int32:
		return uint64(x.(int32)), true
	case int64:
		return uint64(x.(int64)), true
	case uint:
		return uint64(x.(uint)), true
	case uint8:
		return uint64(x.(uint8)), true
	case uint16:
		return uint64(x.(uint16)), true
	case uint32:
		return uint64(x.(uint32)), true
	case uint64:
		return uint64(x.(uint64)), true
	default:
		return 0, false
	}
}

func toFloat64(x interface{}) (float64, bool) {
	switch x.(type) {
	case float32:
		return float64(x.(float32)), true
	case float64:
		return x.(float64), true
	default:
		if u, ok := toUint64(x); ok {
			return float64(u), true
		}
		return 0, false
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
