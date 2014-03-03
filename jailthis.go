package main

import (
	"./jail"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"strconv"
)

func main() {
	var err error
	var username string

	c := jail.NewConfig()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] -- <command> [...args]:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&c.Root, "root", c.Root, "")
	flag.StringVar(&c.Work, "work", c.Work, "")
	flag.StringVar(&username, "user", c.User.Username, "only works when running as root (not in suid)")
	flag.Parse()

	c.Argv = flag.Args()

	_, err = strconv.Atoi(username)
	if err != nil {
		c.User, err = user.LookupId(username)
	} else {
		c.User, err = user.Lookup(username)
	}
	if err != nil {
		panic(err)
	}

	proc, err := jail.Run(c)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals)
	// Forward signals to the child process
	go func() {
		kill := false
		for {
			s := <-signals
			// Force-kill on second interrupt
			if s == os.Interrupt && kill {
				proc.Kill()
				return
			} else {
				if s == os.Interrupt {
					kill = true
				}
				proc.Signal(s)
			}
		}
	}()

	status, err := proc.Wait()
	if err != nil {
		panic(err)
	}

	signal.Stop(signals)

	os.Exit(status)
}

func isDir(path string) {
	os.Stat(path)
}
