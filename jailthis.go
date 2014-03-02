package main

import (
	"./jail"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strconv"
)

func main() {
	var err error
	var u *user.User
	var username string

	c := jail.NewConfig()

	u, err = user.Current()
	if err != nil {
		username = strconv.Itoa(c.Uid)
	}
	username = u.Username

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] -- <command> [...args]:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&c.Root, "root", c.Root, "")
	flag.StringVar(&c.Work, "work", c.Work, "")
	flag.StringVar(&username, "user", username, "only works when running as root (not in suid)")
	flag.Parse()

	c.Argv = flag.Args()

	uid, err := strconv.Atoi(username)
	if err != nil {
		c.Uid = uid
	} else {
		u, err = user.Lookup(username)
		if err != nil {
			panic(err)
		}
		c.Uid, err = strconv.Atoi(u.Uid)
		if err != nil {
			panic(err)
		}
	}

	proc, err := jail.Run(c)
	if err != nil {
		panic(err)
	}

	proc.Wait()

	fmt.Println("WOOT")
}

func isDir(path string) {
	os.Stat(path)
}
