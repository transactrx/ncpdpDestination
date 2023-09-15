package pbmlib

import (
	"time"
)

type PBM interface {
	Start(map[string]interface{}) error
	Post(clm Claim, header map[string][]string, timeout time.Duration) ([]byte, map[string][]string, ErrorInfo)
	Test(claim []byte) ([]byte, ErrorInfo)
	GetStats() Stats
	Shutdown() error
}
