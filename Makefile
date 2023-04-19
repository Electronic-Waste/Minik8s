.PHONY:all build clean help check
BIN=bin
PATH=./cmd
ADM=kubeadm
GO=$(shell which go)
CLEAN=$(shell rm -rf ${BIN})
FOLD=$(shell if [ -d "./bin/" ]; then echo "bin exits"; else mkdir bin;echo "make bin folder"; fi)
all: check build
build:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${ADM}" "${PATH}/${ADM}/${ADM}.go"
clean:
	$(GO) clean
	$(CLEAN)
help:
	@echo "make (all) -- format go code and generate binary in bin dir" 
	@echo "make build -- compile the code and generate binary in bin dir" 
	@echo "make clean -- clean all file generated by build" 
	@echo "make check -- format go code and do some static check" 
check:
	$(GO) fmt $(PATH)/$(ADM)
	$(GO) vet $(PATH)/$(ADM)