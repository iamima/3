package nc

import (
	"log"
)

// ScalarChan is like a chan float32, but with fan-out
// (replicate data over multiple output channels).
// It can only have one input side though.
type FanInScalar []chan float32

// Add a new fanout and return it.
// All fanouts should be created before using the channel.
//func (v *FanInScalar) Fanout(buf int) FanOutScalar {
//	v.fanout = append(v.fanout, make(chan float32, buf))
//	return v.fanout[len(v.fanout)-1]
//}

// Send operator.
func (v FanInScalar) Send(data float32) {
	if Debug && len(v) == 0 {
		log.Println("[WARNING] FanInScalar.Send: no fanout")
	}
	for i := range v {
		v[i] <- data
	}
}

func (f FanInScalar) Close() {
	for i := range f {
		close(f[i])
	}
}
