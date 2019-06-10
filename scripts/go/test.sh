#!/bin/bash
set -e

if [ -n "$IS_SMITHY" ]; then
  cd $GOPATH/src/github.com/Workiva/eva-client-go
fi

# Set the outfile
outfile=gotest.out

export FULL_TESTS="true"
go test -v ./... -ginkgo.noColor -ginkgo.succinct | tr -d '•' | tee $outfile
export FULL_TESTS="false"

# Get go2xunit
which go2xunit > /dev/null || {
    go get github.com/tebeka/go2xunit
}

# Convert the out file to xml so it can be reported to Smithy
go2xunit -input $outfile -output unit_tests.xml

# Get codecov report
pip install goverge
goverge --project_import github.com/Workiva/eva-client-go
