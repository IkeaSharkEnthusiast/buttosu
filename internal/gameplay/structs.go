package gameplay

import "buttplugosu/pkg/mem"

type dynamicAddresses struct {
	IsReady bool
}

var (
	processes []mem.Process
	process   mem.Process
	procerr   error

	previousHits     = 0
	DynamicAddresses = dynamicAddresses{}
)
