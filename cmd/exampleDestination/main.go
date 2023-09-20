package main

import (
	"github.com/transactrx/ncpdpDestination/pkg/dummypbm"
	"github.com/transactrx/ncpdpDestination/pkg/pbmlib"
	"time"
)

func main() {
	dummyPBM := dummypbm.DummyPBM{
		PBMName: "DummyPBM1",
		Latency: time.Millisecond * 10,
	}
	dummyPBM2 := dummypbm.DummyPBM{
		PBMName: "DummyPBM2",
		Latency: time.Millisecond * 50,
	}
	pbmHandler, err := pbmlib.NewPBMHandler()
	if err != nil {
		panic(err)
	}

	err = pbmHandler.HandlePBMS([]pbmlib.PBM{&dummyPBM, &dummyPBM2}, []string{"880"})
	if err != nil {
		panic(err)
	}

	pbmHandler.HandleShutdownAndDrainNats()
}
