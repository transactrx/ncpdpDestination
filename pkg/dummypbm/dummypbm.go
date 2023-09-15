package dummypbm

import (
	"errors"
	"fmt"
	"ncpdpDestination/pkg/pbm"
	"ncpdpDestination/pkg/pbmerrors"
	"time"
)

type DummyPBM struct {
	URL        string
	statistics pbm.Stats
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

func (d *DummyPBM) GetStats() pbm.Stats {
	//TODO implement me
	return pbm.Stats{}
}

func (d *DummyPBM) Shutdown() error {
	return nil
}
