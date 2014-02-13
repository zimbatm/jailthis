jailthis: Runs a command in a jail (or not)
===========================================

Status: WIP

`jailthis --root /some/path --work /tmp -- bash -c "echo > foo"`

* bash must be available at /some/path/bin/bash
* foo will be created at /tmp/foo

On Linux, the jail is strong, featuring cgroups and unshare (TBD).

On MacOS and other POSIX platforms there is no jail. PATH is changed to
/some/path/bin and PWD to /tmp. This is to allow running tests locally before
going to production.


TODO:

* POSIX fake-jail
* unshare and mount syscalls are unmapped
* look into fork/exec

DONE:

* Command-line parsing
* Modularization
