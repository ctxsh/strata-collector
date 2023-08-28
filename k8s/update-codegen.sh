#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

REPO=ctx.sh/strata-collector
SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

source "${CODEGEN_PKG}/kube_codegen.sh"

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.

# Fix up output base
kube::codegen::gen_helpers \
    --input-pkg-root ${REPO}/pkg \
    --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
    --boilerplate "${SCRIPT_ROOT}/k8s/boilerplate.go.txt"

kube::codegen::gen_client \
    --with-watch \
    --input-pkg-root ${REPO}/pkg \
    --output-pkg-root ${REPO}/pkg/client \
    --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
    --boilerplate "${SCRIPT_ROOT}/k8s/boilerplate.go.txt"
