package rapidyenc

/*
#cgo CFLAGS: -I${SRCDIR}/src
#cgo darwin LDFLAGS: ${SRCDIR}/librapidyenc_darwin.a -lstdc++
#cgo windows,amd64 LDFLAGS: ${SRCDIR}/librapidyenc_windows_amd64.a -lstdc++
#cgo linux,amd64 LDFLAGS: ${SRCDIR}/librapidyenc_linux_amd64.a -lstdc++
#cgo linux,arm64 LDFLAGS: ${SRCDIR}/librapidyenc_linux_arm64.a -lstdc++
#include "rapidyenc.h"
*/
import "C"
import (
	"bytes"
	"compress/gzip"
	"sync"
	"unsafe"
)

func MaxLength(length, lineLength int) int {
	return int(C.rapidyenc_encode_max_length(C.size_t(length), C.int(lineLength)))
}

type Encoder struct {
	LineLength int
}

func NewEncoder() *Encoder {
	return &Encoder{
		LineLength: 128,
	}
}

var encodeInitOnce sync.Once

func (e *Encoder) Encode(src []byte) []byte {
	encodeInitOnce.Do(func() {
		C.rapidyenc_encode_init()
	})

	// Compress the source data using gzip
	var compressed bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressed)
	_, err := gzipWriter.Write(src)
	if err != nil {
		panic(err)
	}
	gzipWriter.Close()

	compressedData := compressed.Bytes()
	dst := make([]byte, MaxLength(len(compressedData), e.LineLength))

	length := C.rapidyenc_encode(
		unsafe.Pointer(&compressedData[0]),
		unsafe.Pointer(&dst[0]),
		C.size_t(len(compressedData)),
	)

	return dst[:length]
}
