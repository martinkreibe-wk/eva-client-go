#!/bin/bash

GREEN='\e[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

dirs=$(go list -f {{.Dir}} ./... | grep -v /vendor/ | grep -v /gen-go/)
if [ "$(for d in $dirs; do gofmt -l $d/*.go | tee /dev/stderr; done)" ]; then
    printf "${RED}The above files were not gofmt'ed. Run \"make fmt-go-lib\"${NC}"
    exit 1
fi

printf "${GREEN}All files successfully gofmt'ed${NC}\n"
