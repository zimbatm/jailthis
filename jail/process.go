package jail

import (
	"os"
)

type Process interface {
	Kill() error
	Signal(os.Signal) error
	// TODO: Build an abstraction on top of the syscall package
	Wait() (int, error)
}
