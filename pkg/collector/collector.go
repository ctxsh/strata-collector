package collector

import "ctx.sh/strata-collector/pkg/resource"

type Collector interface {
	SendChan() chan<- resource.Resource
	Start()
	Stop()
	Lock()
	Unlock()
}
