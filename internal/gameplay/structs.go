package gameplay

import "buttplugosu/pkg/memory"

type dynamicAddresses struct {
	IsReady bool
}

var (
	processes []memory.Process
	process   memory.Process
	procerr   error

	previousHits     = 0
	DynamicAddresses = dynamicAddresses{}
)
