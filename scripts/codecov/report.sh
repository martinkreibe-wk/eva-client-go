#!/bin/bash
set -e

if [ -n "$IS_SMITHY" ]; then
  REPORT_FILE="$GOPATH/src/github.com/Workiva/eva-client-go/test_coverage.txt"
else
  REPORT_FILE="test_coverage.txt"
fi

if [ -z "$GIT_BRANCH" ]; then
	echo "GIT_BRANCH environment variable not set, skipping codecov push"
elif [ ! -f "$REPORT_FILE" ]; then
	echo "Cannot find test coverage report file, skipping codecov push"
else
	# upload report
	bash <(curl -s https://codecov.workiva.net/bash) \
	  -u https://codecov.workiva.net \
	  -B $GIT_BRANCH \
	  -r workiva/eva-client-go \
	  -X coveragepy \
	  -f $REPORT_FILE \
	  || echo "ERROR: Codecov failed to upload reports."
fi
