
#!/bin/bash

set -e  #causes the whole script to fail if any commands fail

if [ -n "$IS_SMITHY" ]; then
  cd $GOPATH/src/github.com/Workiva/eva-client-go
fi

# Check go formatting
echo "Checking go code formatting..."
RESULT=$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))
if [ ! -z "$RESULT" ]; then
	echo "Improper go code formatting, need to run gofmt on these files:"
	echo "$RESULT"
	exit 1
else
	echo "Checking go code formatting...done"
fi
