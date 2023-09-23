package service

import "ctx.sh/strata-collector/pkg/resource"

type Collector interface {
	Start(<-chan resource.Resource)
	Stop()
	Lock()
	Unlock()
}
