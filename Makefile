
jailthis: *.go jail/*.go
	go fmt .
	go fmt ./jail
	go build -o $@

clean:
	rm -f jailthis

.PHONY: clean

