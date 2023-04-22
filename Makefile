.PHONY:all build clean help check test
PWD=$(shell pwd)
BIN=build
PATH=cmd
SERVER=apiserver
GO=$(shell which go)
CLEAN=$(shell rm -rf ${BIN})
FOLD=$(shell if [ -d "./$(BIN)/" ]; then echo "$(BIN) exits"; else mkdir $(BIN);echo "make $(BIN) folder"; fi)

build:
	@echo "$(FOLD)"
	$(GO) build -o $(PWD)/$(BIN)/$(SERVER) $(PWD)/$(PATH)/$(SERVER)/$(SERVER).go

clean:
	$(GO) clean
	$(CLEAN)