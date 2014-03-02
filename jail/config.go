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
		Path: []string{"/bin", "/usr/bin"},
		Uid:  os.Getuid(),
		Env:  make(Env),
	}
}
