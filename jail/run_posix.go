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
	argv0 := lookPath(c.Argv[0], c.Root, c.Work, c.Env["PATH"])

	addPrefix(c.Env, "PATH", c.Root)
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

func (self *PosixProcess) Signal(s os.Signal) error {
	return syscall.Kill(self.pid, s.(syscall.Signal))
}

func (self *PosixProcess) Wait() (int, error) {
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
		return filepath.Join(root, command[1:])
	}

	if strings.Contains(command, "/") {
		return filepath.Join(work, command)
	}

	for _, p := range filepath.SplitList(path) {
		x := filepath.Join(root, p[1:], command)
		if isExecutable(x) {
			return x
		}
	}

	// Not found
	return command
}

func addPrefix(env Env, key string, prefix string) {
	v, ok := env[key]
	if !ok {
		return
	}
	p1 := filepath.SplitList(v)
	p2 := make([]string, len(p1))
	for i, p := range p1 {
		if p[0:0] == "/" {
			p = p[1:]
		}
		p2[i] = filepath.Join(prefix, p)
	}
	env[key] = strings.Join(p2, string(os.PathListSeparator))
}
