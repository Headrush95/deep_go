package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	startMemPtr := unsafe.Pointer(&memory[0])
	inUse := make(map[unsafe.Pointer]struct{})
	for _, ptr := range pointers {
		inUse[ptr] = struct{}{}
	}

	for i := 0; i < len(pointers); i++ {
		// у нас же нет гарантии, что pointers отсортирован
		if uintptr(startMemPtr) >= uintptr(pointers[i]) {
			continue
		}
		for {
			if _, occupied := inUse[startMemPtr]; occupied {
				startMemPtr = unsafe.Add(startMemPtr, 1)
				continue
			}
			break
		}
		*(*byte)(startMemPtr), *(*byte)(pointers[i]) = *(*byte)(pointers[i]), *(*byte)(startMemPtr)
		pointers[i] = startMemPtr
		delete(inUse, pointers[i])
		startMemPtr = unsafe.Add(startMemPtr, 1)
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
