package jail

import (
	"os"
	"os/user"
)

type Config struct {
	Root string
	Work string
	User *user.User
	Env  Env
	Argv []string
}

func NewConfig() *Config {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "/"
	}
	user, _ := user.Current()
	return &Config{
		Root: "/",
		Work: cwd,
		User: user,
		Env: Env{
			"PATH":   "/bin:/usr/bin",
			"LC_ALL": "C",
		},
	}
}
