.PHONY:all build clean help check test vctl kubeadm nervctl
BIN=bin
PATH=./cmd
ADM=kubeadm
SERVER=apiserver
VCTL=vctl
NCTL=nervctl
CTLM=kube-controller-manager
CMD=app/cmd
GO=$(shell which go)
CLEAN=$(shell rm -rf ${BIN})
FOLD=$(shell if [ -d "./$(BIN)/" ]; then echo "$(BIN) exits"; else mkdir $(BIN);echo "make $(BIN) folder"; fi)
all: check build
build:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${ADM}" "${PATH}/${ADM}/${ADM}.go"
	$(GO) build -o "${BIN}/${VCTL}" "${PATH}/${VCTL}/${VCTL}.go"
	$(GO) build -o "${BIN}/${NCTL}" "${PATH}/${NCTL}/${NCTL}.go"
	$(GO) build -o "${BIN}/${SERVER}" "${PATH}/${SERVER}/${SERVER}.go"
	$(GO) build -o "${BIN}/${CTLM}" "${PATH}/${CTLM}/${CTLM}.go"
clean:
	$(GO) clean
	$(CLEAN)
help:
	@echo "make (all) -- format go code and generate binary in bin dir"
	@echo "make build -- compile the code and generate binary in bin dir"
	@echo "make clean -- clean all file generated by build"
	@echo "make check -- format go code and do some static check"
	@echo "make test  -- test for some go function code"
check:
	$(GO) fmt $(PATH)/$(ADM)
	$(GO) fmt $(PATH)/$(VCTL)
	$(GO) vet $(PATH)/$(ADM)
	$(GO) vet $(PATH)/$(VCTL)

test:
	$(GO) test $(PATH)/$(ADM)/$(CMD)
kubeadm:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${ADM}" "${PATH}/${ADM}/${ADM}.go"
vctl:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${VCTL}" "${PATH}/${VCTL}/${VCTL}.go"
nervctl:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${NCTL}" "${PATH}/${NCTL}/${NCTL}.go"
apiserver:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${SERVER}" "${PATH}/${SERVER}/${SERVER}.go"
kube-controller-manager:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${CTLM}" "${PATH}/${CTLM}/${CTLM}.go"