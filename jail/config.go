package jail

import (
	"os"
)

type Config struct {
	Root string
	Work string
	Path []string
	Uid  int
	Env  Env
	Argv []string
}

func NewConfig() *Config {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "/"
	}
	return &Config{
		Root: "/",
		Work: cwd,
		Uid:  os.Getuid(),
		Env: Env{
			"PATH":   "/bin:/usr/bin",
			"LC_ALL": "C",
		},
	}
}
