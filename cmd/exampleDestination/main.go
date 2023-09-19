package main

import (
	"github.com/transactrx/ncpdpDestination/pkg/dummypbm"
	"github.com/transactrx/ncpdpDestination/pkg/pbmlib"
	"time"
)

func main() {
	dummyPBM := dummypbm.DummyPBM{
		Name:    "DummyPBM1",
		Latency: time.Second * 2,
	}
	dummyPBM2 := dummypbm.DummyPBM{
		Name:    "DummyPBM2",
		Latency: time.Second * 1,
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
