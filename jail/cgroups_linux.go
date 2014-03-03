package jail

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	// FIXME: make that configurable
	CGROUP_ROOT_DIR = "/sys/fs/cgroup"

	CGROUP_BLKIO      = "blkio"
	CGROUP_CPU        = "cpu"
	CGROUP_CPUACCT    = "cpuacct"
	CGROUP_CPUSET     = "cpuset"
	CGROUP_DEVICES    = "devices"
	CGROUP_FREEZER    = "freezer"
	CGROUP_HUGETLB    = "hugetlb"
	CGROUP_MEMORY     = "memory"
	CGROUP_PERF_EVENT = "perf_event"
)

type Cgroup struct {
	id     string
	groups map[string]bool
}

func NewCgroup() (g *Cgroup) {
	g = &Cgroup{
		id:     randomId(),
		groups: make(map[string]bool),
	}
	return g
}

func randomId() string {
	uuid := make([]byte, 16)

	if _, err := rand.Read(uuid); err != nil {
		panic(err)
	}

	return hex.EncodeToString(uuid)
}
