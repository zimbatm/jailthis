jailthis: command-line jail
===========================

Status: WIP

`jailthis --root /some/path --work /tmp -- bash -c "echo > foo"`

* bash must be available at /some/path/bin/bash
* foo will be created at /tmp/foo

On Linux, the jail is strong, featuring cgroups and unshare (TBD).

On MacOS and other POSIX platforms there is no jail. PATH is changed to
/some/path/bin and PWD to /tmp. This is to allow running tests locally before
going to production. In some ways it operates more like
[Homebrew's superenv](https://github.com/Homebrew/homebrew/wiki/Homebrew-0.9.3).


TODO:

* POSIX fake-jail
* OSX's [sandbox(7)](https://developer.apple.com/library/mac/documentation/Darwin/Reference/ManPages/man7/sandbox.7.html)
* FreeBSD's jail
* unshare and mount syscalls are unmapped
* look into fork/exec
* add a --[no-]network argument

DONE:

* Command-line parsing
* Modularization

Works well with
---------------

* [debootstrap(8)](http://man.cx/debootstrap) and similar tools
* [busybox](http://www.busybox.net/)

See also
--------

* the [venerable chroot(1)](http://man.cx/chroot)
* [systemd-nspawn](http://www.freedesktop.org/software/systemd/man/systemd-nspawn.html)
* [LXC](http://linuxcontainers.org/)
* [docker](http://docker.io)
* [jailkit](http://olivier.sessink.nl/jailkit/)
* [FreeBSD jail](http://www.freebsd.org/cgi/man.cgi?query=jail&format=html)

