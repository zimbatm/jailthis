
jailthis: *.go cgroup/*.go jail/*.go
	go fmt .
	go fmt ./cgroup/
	go fmt ./jail/
	go build -o $@

clean:
	rm -f jailthis

.PHONY: clean

