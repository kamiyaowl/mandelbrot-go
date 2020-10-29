MAINFILE=./cmd/main.go
BINFILE=./main

all: build-wasm build-arch
build-arch:
	go build -o $(BINFILE) -v $(MAINFILE)
build-wasm:
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js web/
	GOOS=js GOARCH=wasm go build -o web/$(BINFILE).wasm $(MAINFILE)
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