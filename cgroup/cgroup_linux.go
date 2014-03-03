package cgroup

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	BLKIO      = "blkio"
	CPU        = "cpu"
	CPUACCT    = "cpuacct"
	CPUSET     = "cpuset"
	DEVICES    = "devices"
	FREEZER    = "freezer"
	HUGETLB    = "hugetlb"
	MEMORY     = "memory"
	PERF_EVENT = "perf_event"
)

// This can change from system to system
var CGROUP_ROOT_DIR = "/sys/fs/cgroup"

func SetRootDir(dir string) {
	CGROUP_ROOT_DIR = dir
}

type GroupMap map[string]*Group

func New(groups ...string) (gm GroupMap, err error) {
	var g *Group

	name := randomId()
	gm = make(GroupMap, len(groups))

	for _, kind := range groups {
		if g, err = newGroup(kind, name); err != nil {
			break
		}
		gm[kind] = g
	}

	if err != nil {
		for _, g = range gm {
			g.teardown()
		}
	}

	return
}

func (self GroupMap) Add(pid int) (err error) {
	for _, g := range self {
		if err = g.add(pid); err != nil {
			return
		}
	}
	return
}

func (self GroupMap) Teardown() {
	for _, g := range self {
		g.teardown()
	}
	return
}

type Group struct {
	kind string
	dir  string
}

func (self *Group) Write(key string, data string) error {
	return ioutil.WriteFile(self.path(key), []byte(data), 0644)
}

func (self *Group) Read(key string) (string, error) {
	data, err := ioutil.ReadFile(self.path(key))
	return string(data), err
}

func newGroup(kind string, name string) (g *Group, err error) {
	dir := filepath.Join(CGROUP_ROOT_DIR, kind, name)
	g = &Group{kind, dir}

	if err = os.Mkdir(dir, 0755); err != nil {
		return
	}
	// We don't want any release agents since we're managing the group
	if err = g.Write("notify_on_release", "0"); err != nil {
		os.Remove(dir)
	}

	return
}

func (self *Group) add(pid int) error {
	return self.Write("tasks", strconv.Itoa(pid))
}

func (self *Group) path(key string) string {
	return filepath.Join(self.dir, key)
}

func (self *Group) teardown() {
	dir := self.dir

	for isDir(dir) {
		if err := syscall.Rmdir(dir); err != nil {
			self.killAll()
		}
	}
}

func (self *Group) killAll() (err error) {
	var d string
	var pid int

	if d, err = self.Read("tasks"); err != nil {
		return
	}

	for _, l := range strings.Split(d, "\n") {
		if pid, _ = strconv.Atoi(strings.TrimSpace(l)); pid > 1 {
			syscall.Kill(pid, syscall.SIGKILL)
		}
	}

	return nil
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
