package pbmlib

import (
	"time"
)

type PBM interface {
	Name() string
	Start() error
	Post(clm Claim, header map[string][]string, timeout time.Duration, privateMessage bool) ([]byte, map[string][]string, ErrorInfo)
	Test(claim []byte) ([]byte, ErrorInfo)
	Shutdown() error
}
