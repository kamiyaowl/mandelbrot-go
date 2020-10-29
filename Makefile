MAINFILE=./cmd/main.go
BINFILE=./main

all: test build
build:
	go build -o $(BINFILE) -v $(MAINFILE)
test:
	go test -v ./
run:
	go run -v $(MAINFILE)
run-wasm:
	wasmserve $(MAINFILE)
docs:
	godoc -http=:6060
init:
	go get golang.org/x/tools/cmd/godoc