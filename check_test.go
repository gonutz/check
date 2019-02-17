package check

import "testing"

func TestEq(t *testing.T) {
	var (
		i    int
		i8   int8
		i16  int16
		i32  int32
		i64  int64
		u    uint
		u8   uint8
		u16  uint16
		u32  uint32
		u64  uint64
		uptr uintptr
		f32  float32
		f64  float64
		c64  complex64
		c128 complex128
	)
	Eq(t, 0, 0)
	Eq(t, i, 0)
	Eq(t, i16, i8)
	Eq(t, i64, i32)
	Eq(t, i64, i)
	Eq(t, u64, u)
	Eq(t, u8, u32)
	Eq(t, uptr, u16)
	Eq(t, 0.0, f32)
	Eq(t, 0, f32)
	Eq(t, f64, f32)
	Eq(t, c64, c128)
	Eq(t, "abc", "abc")
	Eq(t, make(map[int]string), make(map[int]string))
	Eq(t, struct{}{}, struct{}{})
}
