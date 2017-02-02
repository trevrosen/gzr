.PHONY: build clean doc gen run test vet

default: build

build:
	go build -o gzr

clean:
	rm gzr

run:
	go build -o gzr && ./gzr

test:
	go test -v ./..

local_install:
	go install `go list | grep -v /vendor/`

install:
	glide install

doc:
	godoc -http=:8080 -index

vet:
	go vet ./...
