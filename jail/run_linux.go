package jail

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	HOME_DIR = "/home"
)

// A default non-secure implementation for common posix platforms
//
// It just behaves like the jailed version except that it's easy to escape.
func run(c *Config) (proc Process, err error) {
	argv0 := lookPath(c.Argv[0], c.Root, c.Work, c.Env["PATH"])

	c.Env["HOME"] = HOME_DIR
	c.Env["PWD"] = HOME_DIR

	uid, _ := strconv.Atoi(c.User.Uid)
	gid, _ := strconv.Atoi(c.User.Gid)
	attr := &syscall.ProcAttr{
		Dir:   c.Work,
		Env:   c.Env.Cenv(),
		Files: []uintptr{0, 1, 2},
		Sys: &syscall.SysProcAttr{
			Chroot: c.Root,
			Credential: &syscall.Credential{
				uint32(uid),
				uint32(gid),
				nil,
			},
			Ptrace:  true, // Used to pause execution to work around a concurrency issue
			Setsid:  false,
			Setpgid: true,
			Setctty: false,
			Noctty:  false,
			/* Linux only */
			// Ctty: int
			// TODO: How to cleanup sub-processes ?
			Pdeathsig: syscall.SIGTERM,
			// Cloneflags: uintptr
		},
	}

	pid, err := syscall.ForkExec(argv0, c.Argv, attr)
	if err != nil {
		return
	}

	proc = &LinuxProcess{pid}

	return
}

type LinuxProcess struct {
	pid int
}

func (self *LinuxProcess) Kill() error {
	// Send the kill to the whole process group
	return syscall.Kill(-self.pid, syscall.SIGKILL)
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
