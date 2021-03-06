package cuda

// GPU byte slice, used to store regions.

import (
	"github.com/barnex/cuda5/cu"
	"github.com/mumax/3/util"
	"log"
	"unsafe"
)

// 3D byte slice, used for region lookup.
type Bytes struct {
	Ptr unsafe.Pointer
	Len int
}

// Construct new byte slice with given length.
func NewBytes(Len int) *Bytes {
	ptr := cu.MemAlloc(int64(Len))
	cu.MemsetD8(cu.DevicePtr(ptr), 0, int64(Len))
	return &Bytes{unsafe.Pointer(uintptr(ptr)), Len}
}

// Upload src (host) to dst (gpu).
func (dst *Bytes) Upload(src []byte) {
	util.Argument(dst.Len == len(src))
	Sync()
	cu.MemcpyHtoD(cu.DevicePtr(uintptr(dst.Ptr)), unsafe.Pointer(&src[0]), int64(dst.Len))
	Sync()
}

// Copy on device: dst = src.
func (dst *Bytes) Copy(src *Bytes) {
	util.Argument(dst.Len == src.Len)
	Sync()
	cu.MemcpyDtoD(cu.DevicePtr(uintptr(dst.Ptr)), cu.DevicePtr(uintptr(src.Ptr)), int64(dst.Len))
	Sync()
}

// Copy to host: dst = src.
func (src *Bytes) Download(dst []byte) {
	util.Argument(src.Len == len(dst))
	Sync()
	cu.MemcpyDtoH(unsafe.Pointer(&dst[0]), cu.DevicePtr(uintptr(src.Ptr)), int64(src.Len))
	Sync()
}

// Set one element to value
func (dst *Bytes) Set(index int, value byte) {
	if index < 0 || index >= dst.Len {
		log.Panic("Bytes.Set: index out of range:", index)
	}
	src := value
	Sync()
	cu.MemcpyHtoD(cu.DevicePtr(uintptr(dst.Ptr)+uintptr(index)), unsafe.Pointer(&src), 1)
	Sync()
}

// Set one element to value
func (src *Bytes) Get(index int) byte {
	if index < 0 || index >= src.Len {
		log.Panic("Bytes.Set: index out of range:", index)
	}
	var dst byte
	Sync()
	cu.MemcpyDtoH(unsafe.Pointer(&dst), cu.DevicePtr(uintptr(src.Ptr)+uintptr(index)), 1)
	Sync()
	return dst
}

// Frees the GPU memory and disables the slice.
func (b *Bytes) Free() {
	if b.Ptr != nil {
		cu.MemFree(cu.DevicePtr(uintptr(b.Ptr)))
	}
	b.Ptr = nil
	b.Len = 0
}
