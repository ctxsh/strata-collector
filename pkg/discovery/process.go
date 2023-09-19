package discovery

type Process struct {
	Pods      bool
	Services  bool
	Endpoints bool
}

func NewProcess(resources []string) Process {

	return Process{
		Pods:      hasResource(resources, "pods"),
		Services:  hasResource(resources, "services"),
		Endpoints: hasResource(resources, "endpoints"),
	}
}

func hasResource(rs []string, r string) bool {
	for _, v := range rs {
		if v == r {
			return true
		}
	}
	return false
}
