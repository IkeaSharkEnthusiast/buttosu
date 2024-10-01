package gameplay

import "buttplugosu/pkg/mem"

type dynamicAddresses struct {
	IsReady bool
}

var (
	processes        []mem.Process
	process          mem.Process
	previousHits     = 0
	DynamicAddresses = dynamicAddresses{}
)
