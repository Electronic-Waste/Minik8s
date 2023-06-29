.PHONY:all build clean help check test vctl kubeadm nervctl kubelet listener scheduler jobserver
BIN=bin
PATH=./cmd
ADM=kubeadm
SERVER=apiserver
VCTL=vctl
NCTL=nervctl
NATIVE=knative
KUBELET=kubelet
LISTENER=listener
CTLM=kube-controller-manager
CTL=kubectl
CMD=app/cmd
SCH=scheduler
JOB=jobserver
GO=$(shell which go)
CLEAN=$(shell rm -rf ${BIN})
GPUDIR=scripts/gpuscripts
FOLD=$(shell if [ -d "./$(BIN)/" ]; then echo "$(BIN) exits"; else mkdir $(BIN);echo "make $(BIN) folder"; fi)
all: check build
build:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${ADM}" "${PATH}/${ADM}/${ADM}.go"
	$(GO) build -o "${BIN}/${VCTL}" "${PATH}/${VCTL}/${VCTL}.go"
	$(GO) build -o "${BIN}/${NCTL}" "${PATH}/${NCTL}/${NCTL}.go"
	$(GO) build -o "${BIN}/${SERVER}" "${PATH}/${SERVER}/${SERVER}.go"
	$(GO) build -o "${BIN}/${KUBELET}" "${PATH}/${KUBELET}/${KUBELET}.go"
	$(GO) build -o "${BIN}/${CTLM}" "${PATH}/${CTLM}/${CTLM}.go"
	$(GO) build -o "${BIN}/${CTL}" "${PATH}/${CTL}/${CTL}.go"
	$(GO) build -o "${BIN}/${SCH}" "${PATH}/${SCH}/${SCH}.go"
	$(GO) build -o "${BIN}/${LISTENER}" "${PATH}/${LISTENER}/${LISTENER}.go"
	$(GO) build -o "${BIN}/${NATIVE}" "$(PATH)/${NATIVE}/${NATIVE}.go"
	$(GO) build -o "${GPUDIR}/${JOB}" "${PATH}/${JOB}/${JOB}.go"
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
kubelet:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${KUBELET}" "${PATH}/${KUBELET}/${KUBELET}.go"
listener:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${LISTENER}" "${PATH}/${LISTENER}/${LISTENER}.go"
kube-controller-manager:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${CTLM}" "${PATH}/${CTLM}/${CTLM}.go"
kubectl:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${CTL}" "${PATH}/${CTL}/${CTL}.go"
scheduler:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${SCH}" "${PATH}/${SCH}/${SCH}.go"
knative:
	@echo "$(FOLD)"
	$(GO) build -o "${BIN}/${NATIVE}" "$(PATH)/${NATIVE}/${NATIVE}.go"
jobserver:
	@echo "$(FOLD)"
	$(GO) build -o "${GPUDIR}/${JOB}" "${PATH}/${JOB}/${JOB}.go"
