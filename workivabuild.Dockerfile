#! STAGE - Build Client Library
FROM golang:1.12-alpine3.9 as build_go_lib

WORKDIR /go/src/github.com/Workiva/eva-client-go
RUN apk add --update bash curl git gcc libc-dev openssh-client
ENV IS_SMITHY=1

# Install Go Tools
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN go get -u github.com/tebeka/go2xunit

# Verify Dependencies are Vendored
WORKDIR /go/src/github.com/Workiva/eva-client-go/
COPY ./edn ./edn
COPY ./eva ./eva
COPY ./vendor ./vendor
COPY ./Gopkg.toml ./Gopkg.toml
RUN dep ensure -dry-run

# Lint Code
COPY ./scripts ./scripts
RUN ./scripts/ci/gofmt.sh

# Unit-tests
ENV FULL_TESTS="true"
COPY ./test ./test
RUN go test -v -cover -coverprofile=coverage.txt -covermode=atomic ./... -ginkgo.noColor -ginkgo.succinct | tr -d '•' | tee test_reports.txt
RUN go2xunit -input test_reports.txt -output test_reports.xml
ARG BUILD_ARTIFACTS_TEST_REPORT=/go/src/github.com/Workiva/eva-client-go/test_reports.xml

# Code-Coverage Report
ARG GIT_BRANCH
RUN ./scripts/ci/codecov.sh

# Ensure Client Library is Build-able
RUN go build ./...

FROM scratch
