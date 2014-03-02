package jail

import "fmt"

type Env map[string]string

func (self Env) Cenv() (env []string) {
	env = make([]string, len(self))
	i := 0
	for k, v := range self {
		env[i] = fmt.Sprintf("%s=%s", k, v)
		i += 1
	}
	return
}
