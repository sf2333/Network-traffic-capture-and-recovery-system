package ExtractFile

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -L${SRCDIR} -lextract -lstdc++
#include"extract_file.h"
*/
import "C"
import "unsafe"

func ExtractFile(data []byte, dataSize uint32) ([]byte, uint32, int) {
	var outSize uint32
	var outType int


	outSizePtr := unsafe.Pointer(&outSize)
	resultC := C.extractFile((*C.uchar)(C.CBytes(data)),
		C.ulong(dataSize),(*C.ulong)(outSizePtr),(*C.int)(unsafe.Pointer(&outType)))

	resultPointer := unsafe.Pointer(resultC)

	result := C.GoBytes(resultPointer,C.int(outSize))


	return result, outSize, outType
}
