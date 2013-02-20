package cuda

import (
	"code.google.com/p/mx3/data"
	"code.google.com/p/mx3/mag"
)

const DEFAULT_KERNEL_ACC = 6

func NewDemag(m *data.Quant) *Symm2D {
	k := mag.BruteKernel(m.Mesh(), DEFAULT_KERNEL_ACC)
	return NewConvolution(m, k)
}