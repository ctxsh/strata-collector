CRD_OPTIONS ?= "crd:maxDescLen=0,generateEmbeddedObjectMeta=true"
RBAC_OPTIONS ?= "rbac:roleName=strata-role"
WEBHOOK_OPTIONS ?= "webhook"
OUTPUT_OPTIONS ?= "output:artifacts:config=config/base/crd"

deps:
	@echo "Checking dependencies"
	@GO111MODULE=off go get sigs.k8s.io/controller-tools/cmd/controller-gen

ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go install sigs.k8s.io/controller-tools/cmd/controller-gen
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(shell go env GOPATH)/bin/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

codegen:
	@./k8s/update-codegen.sh

manifests:
	$(CONTROLLER_GEN) $(CRD_OPTIONS) $(RBAC_OPTIONS) $(WEBHOOK_OPTIONS) paths="./pkg/apis/..." $(OUTPUT_OPTIONS)

