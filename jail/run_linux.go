package jail

import (
	"../cgroup"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	HOME_DIR = "/home"
	PROC_DIR = "/proc"
	// TODO: disable network ? | syscall.CLONE_NEWNET
	CLONE_FLAGS = syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC

	HOME_MOUNT_OPTS = syscall.MS_BIND | syscall.MS_NOEXEC
)

// A default non-secure implementation for common posix platforms
//
// It just behaves like the jailed version except that it's easy to escape.
func run(c *Config) (proc Process, err error) {
	var creds *syscall.Credential

	argv0 := lookPath(c.Argv[0], c.Root, c.Work, c.Env["PATH"])

	c.Env["HOME"] = HOME_DIR
	c.Env["PWD"] = HOME_DIR

	uid, _ := strconv.Atoi(c.User.Uid)
	gid, _ := strconv.Atoi(c.User.Gid)
	if uid != syscall.Getuid() || gid != syscall.Getgid() {
		creds = &syscall.Credential{
			uint32(uid),
			uint32(gid),
			nil,
		}
	}
	attr := &syscall.ProcAttr{
		Dir:   c.Work,
		Env:   c.Env.Cenv(),
		Files: []uintptr{0, 1, 2},
		Sys: &syscall.SysProcAttr{
			Chroot:     c.Root,
			Credential: creds,
			Ptrace:     true, // Used to pause execution to work around a concurrency issue
			Setsid:     false,
			Setpgid:    true,
			Setctty:    false,
			Noctty:     false,
			/* Linux only */
			// Ctty: int
			// TODO: How to cleanup sub-processes ?
			Pdeathsig: syscall.SIGTERM,
			// Cloneflags: uintptr
		},
	}

	// New filesystem namespace
	if err = syscall.Unshare(CLONE_FLAGS); err != nil {
		return
	}

	// Mount the work directory
	if err = syscall.Mount(c.Work, filepath.Join(c.Root, HOME_DIR[1:]), "none", HOME_MOUNT_OPTS, ""); err != nil {
		return
	}

	// Mount the proc filesystem
	procDir := filepath.Join(c.Root, PROC_DIR[1:])
	if isDir(procDir) && c.Root != "/" {
		if err = syscall.Mount(c.Work, procDir, "proc", 0, ""); err != nil {
			return
		}
	}

	g, err := cgroup.New(cgroup.CPUACCT)
	if err != nil {
		return
	}
	defer g.Teardown()

	pid, err := syscall.ForkExec(argv0, c.Argv, attr)
	if err != nil {
		return
	}

	if err = g.Add(pid); err != nil {
		return
	}

	if err = syscall.PtraceDetach(pid); err != nil {
		return
	}

	return &LinuxProcess{pid, g}, nil
}

type LinuxProcess struct {
	pid    int
	cgroup cgroup.GroupMap
}

func (self *LinuxProcess) Kill() error {
	self.cgroup.Teardown()
	return nil
}

func (self *LinuxProcess) Signal(s os.Signal) error {
	return syscall.Kill(self.pid, s.(syscall.Signal))
}

func (self *LinuxProcess) Wait() (int, error) {
	var ws syscall.WaitStatus
	//var rs syscall.Rusage
	var err error
	_, err = syscall.Wait4(self.pid, &ws, 0, nil)
	// Cleanup
	self.Kill()
	return ws.ExitStatus(), err
}

func lookPath(command string, root string, work string, path string) string {
	command = filepath.Clean(command)

	if command[0:0] == "/" {
		return command[1:]
	}

	if strings.Contains(command, "/") {
		return filepath.Join(work, command)
	}

	for _, p := range filepath.SplitList(path) {
		x := filepath.Join(root, p[1:], command)
		if isExecutable(x) {
			return filepath.Join(p, command)
		}
	}

	// Not found
	return command
}
