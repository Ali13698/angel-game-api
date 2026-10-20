package main

import (
	"encoding/json"
	"unsafe"
)

const heapSize = 1 << 20

var heap [heapSize]byte
var heapOff uint32

// heapAddr is the only place a pointer is converted to uintptr.
func heapAddr() uintptr { return uintptr(unsafe.Pointer(&heap[0])) }

//go:wasmexport reset
func reset() { heapOff = 0 }

//go:wasmexport allocate
func allocate(size uint32) uint32 {
	if size == 0 || size > uint32(len(heap)) || heapOff > uint32(len(heap))-size {
		return 0
	}
	ptr := uint32(heapAddr()) + heapOff
	heapOff += size
	return ptr
}

type wasmResult struct {
	State  State   `json:"state"`
	Events []Event `json:"events"`
}

func sliceAt(addr, n uint32) []byte {
	if addr == 0 || n == 0 {
		return nil
	}
	off := addr - uint32(heapAddr())
	if off > uint32(len(heap)) || n > uint32(len(heap))-off {
		return nil
	}
	return heap[off : off+n]
}

//go:wasmexport apply
func apply(statePtr, stateLen, intentPtr, intentLen, ctxPtr, ctxLen uint32) uint32 {
	stateB := sliceAt(statePtr, stateLen)
	intentB := sliceAt(intentPtr, intentLen)
	_ = sliceAt(ctxPtr, ctxLen)

	var state State
	if err := json.Unmarshal(stateB, &state); err != nil {
		return 0
	}
	var intent Intent
	if err := json.Unmarshal(intentB, &intent); err != nil {
		return 0
	}

	newState, events, err := ApplyCommand(state, intent.From, intent.Cmd.Type, intent.Cmd.Payload)
	if err != nil {
		return 0
	}

	resB, err := json.Marshal(wasmResult{State: newState, Events: events})
	if err != nil {
		return 0
	}

	total := uint32(4 + len(resB))
	ptr := allocate(total)
	if ptr == 0 {
		return 0
	}

	dst := sliceAt(ptr, total)
	if dst == nil {
		return 0
	}
	dst[0] = byte(len(resB))
	dst[1] = byte(len(resB) >> 8)
	dst[2] = byte(len(resB) >> 16)
	dst[3] = byte(len(resB) >> 24)
	copy(dst[4:], resB)
	return ptr
}

func main() {}
