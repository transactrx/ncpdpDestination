package dummypbm

import (
	"github.com/transactrx/ncpdpDestination/pkg/pbmlib"
	"log"
	"time"
)

//	Start(map[string]interface{}) error
//	Post(clm Claim, header map[string][]string, timeout time.Duration) ([]byte, ErrorInfo)
//	Test(claim []byte) ([]byte, ErrorInfo)
//	GetStats() Stats
//	Shutdown() error

type DummyPBM struct {
	PBMName    string
	Latency    time.Duration
	statistics pbmlib.Stats

	lastResponseGood bool
}

func (d *DummyPBM) Start() error {
	return nil
}
func (d *DummyPBM) Name() string {
	return d.PBMName
}

func (d *DummyPBM) Post(claim pbmlib.Claim, header map[string][]string, timeout time.Duration, privateMessage bool) ([]byte, map[string][]string, pbmlib.ErrorInfo) {

	log.Printf("Got Message: %s", d.Name())
	time.Sleep(d.Latency)
	good := true
	if d.lastResponseGood {
		good = false
		d.lastResponseGood = false
	} else {
		d.lastResponseGood = true
	}
	if good {
		return []byte("GOOD RESPONSE"), nil, pbmlib.ErrorCode.TRX00
	}
	return []byte(""), nil, pbmlib.ErrorCode.TRX01
}

func (d *DummyPBM) Test(claim []byte) ([]byte, pbmlib.ErrorInfo) {
	//TODO implement me
	return nil, pbmlib.ErrorCode.TRX00
}

func (d *DummyPBM) GetStats() pbmlib.Stats {
	//TODO implement me
	return pbmlib.Stats{}
}

func (d *DummyPBM) Shutdown() error {
	return nil
}
