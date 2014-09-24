DEPENDENCIES := github.com/droundy/goopt \
				github.com/stretchr/testify \
				github.com/mattbaird/elastigo 

all: build test

build:
		go build gostats.go 

test:
		go test -v $(PACKAGES)

format:
		go fmt $(PACKAGES)

deps:
		go get $(DEPENDENCIES)
