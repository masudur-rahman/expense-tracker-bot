#!/bin/bash

set -eou pipefail

export CGO_ENABLED=0
export GO111MODULE=on
export GOFLAGS="-mod=vendor"

TARGETS="$@"
COMPANY_PREFIXES="github.com/masudur-rahman"
IMPORTS_ORDER="std,project,company,general"

if [ -n "$TARGETS" ]; then
    echo "Running goimports:"
    set cmd
#    cmd="goimports-reviser -recursive -imports-order=${IMPORTS_ORDER} -format ${TARGETS}"
    cmd="goimports-reviser -recursive -company-prefixes=${COMPANY_PREFIXES} -imports-order=${IMPORTS_ORDER} -format ./..."
    echo "$cmd"
    $cmd
    echo

#    echo "Running gofmt:"
#    cmd="gofmt -s -w ${TARGETS}"
#    echo "$cmd"
#    $cmd
#    echo
fi

#echo "Running shfmt:"
#cmd="find . -path ./vendor -prune -o -name '*.sh' -exec shfmt -l -w -ci -i 4 {} \;"
#echo "$cmd"
#eval "$cmd" # xref: https://stackoverflow.com/a/5615748/244009
#echo
