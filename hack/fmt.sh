#!/bin/bash

set -eou pipefail

export CGO_ENABLED=0
export GO111MODULE=on
export GOFLAGS="-mod=vendor"

#TARGETS="$@"
EXCLUDES="$@"
: ${EXCLUDES:=vendor}
COMPANY_PREFIXES="github.com/masudur-rahman"
IMPORTS_ORDER="std,project,company,general,blanked"

format_go_files() {
    echo "Running goimports-reviser:"
    set cmd
    cmd="goimports-reviser -recursive -company-prefixes=${COMPANY_PREFIXES} -imports-order=${IMPORTS_ORDER} -format -excludes ${EXCLUDES} ./..."
    echo "$cmd"
    $cmd
}

format_script_files() {
    echo "Running shfmt:"
    cmd="find . -path ./vendor -prune -o -name '*.sh' -exec shfmt -l -w -ci -i 4 {} \;"
    echo "$cmd"
    eval "$cmd" # xref: https://stackoverflow.com/a/5615748/244009
    echo
}

format_go_files
format_script_files
