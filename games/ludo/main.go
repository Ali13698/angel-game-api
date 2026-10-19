package main

import (
	"encoding/json"
	"unsafe"
)

const heapSize = 1 << 20
var heap [heapSize]byte
var heapOff uint32

//go:wasmexport reset
func reset() { heapOff = 0 }

//go:wasmexport allocate
func allocate(size uint32) uint32 {
	if size == 0 || heapOff > uint32(len(heap))-size { return 0 }
	ptr := uint32(uintptr(unsafe.Pointer(&heap[heapOff])))
	heapOff += size
	return ptr
}

type wasmResult struct {
	State  State   `json:"state"`
	Events []Event `json:"events"`
}

//go:wasmexport apply
func apply(statePtr, stateLen, intentPtr, intentLen, ctxPtr, ctxLen uint32) uint32 {
	stateB := readBytes(statePtr, stateLen)
	intentB := readBytes(intentPtr, intentLen)
	_ = readBytes(ctxPtr, ctxLen) // ctx is not used for now

	var state State
	if err := json.Unmarshal(stateB, &state); err != nil { return 0 }

	var intent Intent
	if err := json.Unmarshal(intentB, &intent); err != nil { return 0 }

	newState, events, err := ApplyCommand(state, intent.From, intent.Cmd.Type, intent.Cmd.Payload)
	if err != nil { return 0 } // In a production environment, we should return the error payload

	res := wasmResult{State: newState, Events: events}
	resB, _ := json.Marshal(res)

	return writeResult(resB)
}

func writeResult(data []byte) uint32 {
	total := uint32(4 + len(data))
	ptr := allocate(total)
	if ptr == 0 { return 0 }
	dst := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), total)
	dst[0] = byte(len(data))
	dst[1] = byte(len(data) >> 8)
	dst[2] = byte(len(data) >> 16)
	dst[3] = byte(len(data) >> 24)
	copy(dst[4:], data)
	return ptr
}

func readBytes(ptr, n uint32) []byte {
	if ptr == 0 || n == 0 { return nil }
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), n)
}

func main() {}
