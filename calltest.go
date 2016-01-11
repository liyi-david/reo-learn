package main

import reo "./lib/reo"
import "fmt"
import "sync"

func alternator(A, B, C reo.Port) {
	// definition of ports
	M0 := reo.MakePort()
	M1 := reo.MakePort()
	M2 := reo.MakePort()
	M3 := reo.MakePort()
	M4 := reo.MakePort()
	M5 := reo.MakePort()
	M6 := reo.MakePort()
	StopFlag := reo.GenerateStopPort(6)

	// definition of channels
	go reo.ReplicatorChannel(A, M0, M1, StopFlag[0])
	go reo.ReplicatorChannel(B, M2, M3, StopFlag[1])
	go reo.MergerChannel(M4, M5, M6, StopFlag[2])

	go reo.SyncdrainChannel(M1, M2, StopFlag[3])
	go reo.SyncChannel(M0, M4, StopFlag[4])
	go reo.FifoChannel(M3, M5, StopFlag[5])
}

func sender(port reo.Port, msg string) {
	for {
		port.SyncWrite(msg)
	}
}

func monitor(port reo.Port) {
	for {
		fmt.Println("[MONITOR]", port.SyncRead())
	}
}

func main() {
	// configurations
	//reo.CloseLog()

	var wg sync.WaitGroup
	A := reo.MakePort()
	B := reo.MakePort()
	C := reo.MakePort()

	// connector startup
	alternator(A, B, C)
	// running components
	wg.Add(1)
	go sender(A, "MSG A")
	go sender(B, "MSG B")
	go monitor(C)
	wg.Wait()
}
