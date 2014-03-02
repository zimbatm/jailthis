package jail

import (
	"syscall"
)

type Process interface {
	Kill() error
	Signal(syscall.Signal) error
	Wait() (*syscall.WaitStatus, error)
}
