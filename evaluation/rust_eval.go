package evaluation

// #cgo CFLAGS: -I${SRCDIR}/../rust
// #cgo LDFLAGS: -L${SRCDIR}/../rust/expr_tree/target/release -lexpr_tree
// #include "expr_tree/src/expr_tree.h"
// #include <stdlib.h>
import "C"
import (
	"reflect"
	"unsafe"
)

// EvalFromBytesRust passes the given flatbuffer to Rust for evaluation,
// and returns the result.
func EvalFromBytesInRust(bs []byte) (float64, error) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	data := (*C.char)(unsafe.Pointer(sh.Data))
	ln := C.int(sh.Len)
	ans := C.eval_from_c(data, ln)
	return float64(ans), nil
}

// GetBytesFromRust() returns a byte slice that contains a
// simple expression tree in a flatbuffer.
// The returned function must be called to free the memory
// when it is no longer needed.
func GetBytesFromRust() ([]byte, func(), error) {
	var ln C.int
	var offset C.int
	ptr := C.get_expr_tree(&ln, &offset)
	sh := new(reflect.SliceHeader)
	sh.Data = uintptr(unsafe.Pointer(ptr))
	sh.Len = int(ln)
	sh.Cap = int(ln)
	bs := *(*[]byte)(unsafe.Pointer(sh))
	bs = bs[offset:]
	freeFn := func() {
		C.free_expr_tree(ptr, ln)
	}
	return bs, freeFn, nil
}
