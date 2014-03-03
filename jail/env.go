package jail

import (
	"errors"
	"strings"
)

type Env map[string]string

func (self Env) Cenv() (env []string) {
	env = make([]string, len(self))
	i := 0
	for k, v := range self {
		env[i] = strings.Join([]string{k, v}, "=")
		i += 1
	}
	return
}

func (self Env) String() string {
	return ""
}

func (self Env) Set(str string) error {
	kv := strings.SplitN(str, "=", 2)
	if len(kv) != 2 {
		return errors.New("env is not in form key=value")
	}
	self[kv[0]] = kv[1]
	return nil
}
