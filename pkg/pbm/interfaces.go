package pbm

import (
	"time"
)

type PBM interface {
	Start(map[string]interface{}) error
	Post(claim []byte, header map[string][]string, timeout time.Duration) ([]byte, ErrorInfo)
	Test(claim []byte) ([]byte, ErrorInfo)
	GetStats() Stats
	Shutdown() error
}
