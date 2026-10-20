//go:build wasip1

package main

import "unsafe"

const heapSize = 1 << 20
var heap [heapSize]byte
var heapOff uint32

//go:wasmexport reset
func reset() { heapOff = 0 }

//go:wasmexport allocate
func allocate(size uint32) uint32 {
	if size == 0 || size > uint32(len(heap)) || heapOff > uint32(len(heap))-size {
		return 0
	}
	ptr := uint32(uintptr(unsafe.Pointer(&heap[heapOff])))
	heapOff += size
	return ptr
}

//go:wasmexport apply
func apply(statePtr, stateLen, intentPtr, intentLen, ctxPtr, ctxLen uint32) uint32 {
	out := []byte(`{"ok":true,"msg":"state_updated"}`)
	total := uint32(4 + len(out))
	
	if heapOff > uint32(len(heap))-total {
		return 0
	}
	
	dst := heap[heapOff : heapOff+total]
	ptr := uint32(uintptr(unsafe.Pointer(&heap[heapOff])))
	heapOff += total
	
	dst[0] = byte(len(out))
	dst[1] = byte(len(out) >> 8)
	dst[2] = byte(len(out) >> 16)
	dst[3] = byte(len(out) >> 24)
	copy(dst[4:], out)
	
	return ptr
}

func main() {}
