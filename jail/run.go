package jail

import (
	"fmt"
	"os"
	"path/filepath"
)

func Run(c *Config) (proc Process, err error) {
	if c.Argv == nil || len(c.Argv) == 0 {
		err = fmt.Errorf("command missing")
		return
	}

	if !isDir(c.Root) {
		err = fmt.Errorf("root is not a directory")
		return
	}

	if !isDir(c.Work) {
		err = fmt.Errorf("work is not a directory")
		return
	}

	if c.Root, err = filepath.Abs(c.Root); err != nil {
		return
	}

	if c.Work, err = filepath.Abs(c.Work); err != nil {
		return
	}

	// Make sure processes have a language flag
	if _, ok := c.Env["LC_ALL"]; !ok {
		c.Env["LC_ALL"] = "C"
	}

	return run(c)
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// FIXME: Also take the uid/gid into account
// FIXME: Also check the executable bit maybe :p
func isExecutable(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}
