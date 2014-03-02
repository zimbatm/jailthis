// +build !linux

package jail

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// A default non-secure implementation for common posix platforms
//
// It just behaves like the jailed version except that it's easy to escape.
func run(c *Config) (proc Process, err error) {
	argv0 := lookPath(c.Argv[0], c.Root, c.Work, c.Path)

	paths := make([]string, len(c.Path))
	for i, p := range c.Path {
		paths[i] = filepath.Join(c.Root, p)
	}
	c.Env["PATH"] = strings.Join(paths, string(os.PathListSeparator))
	c.Env["HOME"] = c.Work
	c.Env["PWD"] = c.Work

	attr := &syscall.ProcAttr{
		Dir:   c.Work,
		Env:   c.Env.Cenv(),
		Files: []uintptr{0, 1, 2},
		Sys: &syscall.SysProcAttr{
			Chroot:     "",
			Credential: nil,
			Ptrace:     false,
			Setsid:     false,
			Setpgid:    true, // Create a new process group
			Setctty:    false,
			Noctty:     false,
		},
	}

	pid, err := syscall.ForkExec(argv0, c.Argv, attr)
	if err != nil {
		return
	}

	proc = &PosixProcess{pid}

	return
}

type PosixProcess struct {
	pid int
}

func (self *PosixProcess) Kill() error {
	// Send the kill to the whole process group
	return syscall.Kill(-self.pid, syscall.SIGKILL)
}

func (self *PosixProcess) Signal(s syscall.Signal) error {
	return syscall.Kill(self.pid, s)
}

func (self *PosixProcess) Wait() (ws *syscall.WaitStatus, err error) {
	_, err = syscall.Wait4(self.pid, ws, 0, nil)
	return
}

func lookPath(command string, root string, work string, paths []string) string {
	command = filepath.Clean(command)

	if command[0:0] == "/" {
		return filepath.Join(root, command[1:])
	}

	if strings.Contains(command, "/") {
		return filepath.Join(work, command)
	}

	for _, path := range paths {
		x := filepath.Join(root, path[1:], command)
		if isExecutable(x) {
			return x
		}
	}

	return command
}
