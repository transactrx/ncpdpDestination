package dummypbm

import (
	"errors"
	"fmt"
	"ncpdpDestination/pkg/pbmerrors"
	"ncpdpDestination/pkg/pbmlib"
	"time"
)

type DummyPBM struct {
	URL        string
	statistics pbmlib.Stats
}

func (d *DummyPBM) Start(m map[string]interface{}) error {

	url, ok := m["URL"]
	if !ok {
		return errors.New("URL Required")
	}
	d.URL = fmt.Sprintf("%v", url)

	//TODO implement me
	return nil
}

func (d *DummyPBM) Post(claim []byte, header map[string][]string, timeout time.Duration) ([]byte, pbmerrors.PBMError) {
	//TODO implement me

	//tcp connect

	//update
	return nil, nil
}

func (d *DummyPBM) Test(claim []byte) ([]byte, pbmerrors.PBMError) {
	//TODO implement me
	return nil, nil
}

func (d *DummyPBM) GetStats() pbmlib.Stats {
	//TODO implement me
	return pbmlib.Stats{}
}

func (d *DummyPBM) Shutdown() error {
	return nil
}
