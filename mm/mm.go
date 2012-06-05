package mm

import (
	"fmt"
	"log"
	. "nimble-cube/nc"
)

var (
	size [3]int // 3D geom size
	N    int    // product of size
	warp int    // buffer size for Range()

//	tChan      ScalarChan // Distributes time. Close means teardown listeners.
//	mChan      VecChan
//	alphaChan  Chan
//	hChan      VecChan
//	torqueChan VecChan
)

func Main() {

	initSize()

	tBox := new(TorqueBox)
	hBox := new(MeanFieldBox)
	eBox := new(EulerBox)
	eBox.dt = 0.01

	Connect3(&(hBox.m), &(eBox.m))
	Connect3(&(tBox.m), &(eBox.m))
	Connect3(&(tBox.h), &(hBox.h))
	Connect3(&(eBox.t), &(tBox.t))

	//	Probe3(&(eBox.m), "m")
	//	Probe3(&(hBox.h), "h")
	//	Probe3(&(tBox.t), "t")

	go tBox.Run()
	go hBox.Run()

	m0 := [3][]float32{make([]float32, N), make([]float32, N), make([]float32, N)}
	Memset3(m0, Vector{0.1, 0.99, 0})

	for i := 0; i < 1000; i++ {
		eBox.Run(m0, 10)
		fmt.Println(m0[X][0], m0[Y][0], m0[Z][0])
	}

}

func DefaultBufSize() int { return N / warp }

func Connect3(dst *FanOut3, src *FanIn3) {
	buf := DefaultBufSize()
	*dst = src.FanOut(buf)
}

//	-------
//	CONCEPT:
//
//	magnet struct{
//		m, h VecChan
//	}
//
//	TorqueBox struct{
//		m, h VecRecv
//		alpha Recv
//	}
//
//	t := TorqueBox
//	Read(&t.m, magnet.m)	
//	Read(&t.a, magnet.a)	
//	Read(&t.h, magnet.h)	
//	//Write(magnet.torque, t)
//	//Register(t.Start())
//
//	h := VecConst
//	Connect(	
//
//	// start all now
//	go torque.Run()	
//	go h.Run()	
//	go alpha.Run()	
//	...
//
//	-----

//	go RunGC()
//
//	go SendConst(tChan.Fanout(1), alphaChan, 0.05)
//	go SendVecConst(tChan.Fanout(1), hChan, Vector{0, 1, 0})
//	go RunTorque(tChan.Fanout(1), torqueChan, mChan.Fanout(1), hChan.Fanout(1), alphaChan.Fanout(1))
//
//	const dt = 0.001
//	const Nsteps = 10
//
//	// init buffered m channel with starting value.
//	mx := make([]float32, warp)
//	my := make([]float32, warp)
//	mz := make([]float32, warp)
//	Memset(mx, 1)
//	mInit := [3][]float32{mx, my, mz}
//
//	mOffline := MakeVecChan()
//	mOffRecv := mOffline.Fanout(N / warp)
//	for I := 0; I < N; I += warp {
//		mOffline.Send(mInit) // [I]
//	}
//
//	go func() {
//
//		mRecv := mChan.Fanout(N / warp)
//		torqueRecv := torqueChan.Fanout(1)
//
//		for I := 0; I < warp; I += warp {
//			mChan.Send(mOffRecv.Recv())
//		}
//
//		for step := 0; step < Nsteps; step++ {
//
//			// loop over blocks
//			for I := 0; I < N; I += warp {
//
//				newMList := VecBuffer()
//				mList := mRecv.Recv()
//				tqList := torqueRecv.Recv()
//
//				// loop inside blocks
//				for i := range newMList[X] {
//					newMList[X][i] = dt*tqList[X][i] + mList[X][i]
//					newMList[Y][i] = dt*tqList[Y][i] + mList[Y][i]
//					newMList[Z][i] = dt*tqList[Z][i] + mList[Z][i]
//
//				}
//				if step != Nsteps-1 {
//					mChan.Send(newMList)
//				} else {
//					mOffline.Send(newMList)
//				}
//
//			}
//		}
//
//	}()
//
//	go func(){for s:=0; s<Nsteps; s++{
//		log.Println("tick")
//		tChan.Send(1)
//	}}()
//
//	for I := 0; I < N; I += warp {
//		log.Println("m:", mOffRecv.Recv())
//	}

func initSize() {
	N0, N1, N2 := 1, 4, 8
	size := [3]int{N0, N1, N2}
	N = N0 * N1 * N2

	log.Println("size:", size)
	N := N0 * N1 * N2
	log.Println("N:", N)

	// Find some nice warp size
	warp = 8
	for N%warp != 0 {
		warp--
	}
	log.Println("warp:", warp)

}

//CONCEPT:

//func RunTorque(){
//	// replace by := notation
//	var recvm chan<- float32[] = recv(m) // engine inserts tee if needed 
//		// engine uses runtime.Caller to construct (purely informative) dependency graph:  torque <- RunTorque <- (m, h)
//	var recvh chan<- float32[] = recv(h)
//	var sendtorque <-chan[]float32 = send(torque)
//
//	for{
//		buf := <- getbuffer
//		m:=<-recvm
//		h:=<-recvh
//		torque(buf, m, h)
//		sendtorque <- torque
//		recycle <- m
//		recycle <- h
//	}
//}
