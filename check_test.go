package check_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/gonutz/check"
)

type mockTester struct {
	err string
}

func (t *mockTester) Helper() {}
func (t *mockTester) Errorf(format string, a ...interface{}) {
	t.err = fmt.Sprintf(format, a...)
}

func TestEq(t *testing.T) {
	eq := func(a, b interface{}) {
		t.Helper()
		{
			// check Eq(a,b) prints no error
			var tt mockTester
			check.Eq(&tt, a, b)
			if tt.err != "" {
				t.Errorf("%v == %v but error for Eq was %q", a, b, tt.err)
			}
		}
		{
			// check Eq(b,a) prints no error
			var tt mockTester
			check.Eq(&tt, b, a)
			if tt.err != "" {
				t.Errorf("%v == %v but error for Eq was %q", b, a, tt.err)
			}
		}
		{
			// check Neq(a,b) prints an error
			var tt mockTester
			check.Neq(&tt, a, b)
			if tt.err == "" {
				t.Errorf("%v != %v but have no error for Neq", a, b)
			}
		}
		{
			// check Neq(b,a) prints an error
			var tt mockTester
			check.Neq(&tt, b, a)
			if tt.err == "" {
				t.Errorf("%v != %v but have no error for Neq", b, a)
			}
		}
	}
	neq := func(a, b interface{}) {
		t.Helper()
		{
			// check Eq(a,b) prints an error
			var tt mockTester
			check.Eq(&tt, a, b)
			if tt.err == "" {
				t.Errorf("%v != %v but have no error for Eq", a, b)
			}
		}
		{
			// check Eq(b,a) prints an error
			var tt mockTester
			check.Eq(&tt, b, a)
			if tt.err == "" {
				t.Errorf("%v != %v but have no error for Eq", b, a)
			}
		}
		{
			// check Neq(a,b) prints no error
			var tt mockTester
			check.Neq(&tt, a, b)
			if tt.err != "" {
				t.Errorf("%v == %v but error for Neq was %q", a, b, tt.err)
			}
		}
		{
			// check Neq(b,a) prints no error
			var tt mockTester
			check.Neq(&tt, b, a)
			if tt.err != "" {
				t.Errorf("%v == %v but error for Neq was %q", b, a, tt.err)
			}
		}
	}

	eq(0, 0)
	neq(1, 0)
	eq(nil, nil)
	eq(int8(2), uint64(2))
	eq(float32(1.2), float64(1.2))
	eq(complex(1, 2), complex(1, 2))
	neq(complex(1, 2), complex(1, 456))
	eq(complex64(complex(1, 2)), complex128(complex(1, 2)))
	eq(4, 4.0)
	eq(-4.0, -4)
	neq(1.2, "1.2")
	var i, j int
	eq(&i, &i)
	eq(&i, &j)
	eq(1.2, 1.20000001)
	neq(uint64(0xFFFFFFFFFFFFFFFF), int64(-1))
	eq(int32(5), int8(5))
	eq(uint32(5), uint64(5))
	eq(uint64(5), uint64(5))
	neq(uint64(999), uint64(5))
	eq(2.0, uint64(2))
	eq("abc", "abc")
	neq("abc", "ABC")

	var s1, s2 struct{ a int }
	eq(s1, s1)
	eq(s1, s2)
	s2.a = 111
	neq(s1, s2)

	p1 := unsafe.Pointer(&s1)
	p2 := unsafe.Pointer(&s2)
	eq(p1, p1)
	neq(p1, p2)
	neq(p1, &p1)
	eq(&p1, &p1)

	eq(true, true)
	neq(true, false)

	eq(eq, eq)
	neq(eq, neq)
	var nilF func()
	eq(nilF, nilF)

	eq([]int{1, 2, 3}, []int{1, 2, 3})
	neq([]int{1, 2, 3}, []int{1, 2})
	neq([]int{1, 2, 3}, []int{1, 2, 4})
	neq([]int{1, 2, 3}, []int(nil))
	slice := []int{1, 2, 3}
	eq(slice, slice)

	var ints2 [2]int
	var ints3 [3]int
	eq(ints2, ints2)
	neq(ints2, ints3)
	neq([2]int{1, 2}, [2]int{1, 3})

	var nilMap map[int]string
	m1 := map[int]string{1: "abc"}
	m2 := map[int]string{1: "abc", 2: "def"}
	m22 := map[int]string{1: "abc", 2: "def"}
	m3 := map[int]string{1: "abc", 2: "DEF"}
	eq(m1, m1)
	neq(m1, m2)
	neq(nilMap, m1)
	eq(nilMap, nilMap)
	neq(m2, m3)
	eq(m2, m22)

	type bbb struct{}
	type aaa struct {
		bbb
		u unsafe.Pointer
		i interface {
			a()
		}
	}
	var aa1, aa2 aaa
	eq(aa1, aa2)
	aa1.u = unsafe.Pointer(&aa1)
	aa2.u = unsafe.Pointer(&aa2)
	eq(aa1, aa1)
	neq(aa1, aa2)

	aa2.u = aa1.u
	eq(aa1, aa2)

	aa1.i = aer{i: 1}
	aa2.i = aer{i: 1}
	eq(aa1, aa2)

	aa1.i = aer{i: 1}
	aa2.i = aer{i: 2}
	neq(aa1, aa2)

	neq([]interface{}{""}, []interface{}{5})
	eq([]interface{}{7.0, 6.0, 5.0}, []interface{}{7, 6, 5})

	eq("abc", "abc")
	eq("abc", []byte("abc"))
	eq("abc", []rune("abc"))
	neq("XXX", "abc")
	neq("XXX", []byte("abc"))
	neq("XXX", []rune("abc"))
}

type aer struct{ i int }

func (aer) a() {}

func TestEqHasMessage(t *testing.T) {
	var tt mockTester
	check.Eq(&tt, 1, 2, "message")
	if tt.err != "message: 1 != 2" {
		t.Error(tt.err)
	}

	tt.err = ""
	check.Neq(&tt, 1, 1, "wat")
	if tt.err != "wat: 1 == 1" {
		t.Error(tt.err)
	}
}

func TestEqExact(t *testing.T) {
	var tt mockTester
	check.EqExact(&tt, 1.0, 1.0)
	if tt.err != "" {
		t.Error("no error expected")
	}

	tt.err = ""
	check.EqExact(&tt, 1.0, 1.00000001)
	if tt.err == "" {
		t.Error("error expected")
	}

	tt.err = ""
	check.NeqExact(&tt, 1.0, 1.0)
	if tt.err == "" {
		t.Error("error expected")
	}

	tt.err = ""
	check.NeqExact(&tt, 1.0, 1.00000001)
	if tt.err != "" {
		t.Error("no error expected", tt.err)
	}
}
