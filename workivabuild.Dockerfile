FROM drydock-prod.workiva.net/workiva/smithy-runner-generator:296096

ARG GIT_BRANCH

ARG BUILD_ARTIFACTS_TEST_REPORT=/go/src/github.com/Workiva/eva-client-go/unit_tests.xml

ENV IS_SMITHY=1
ENV REPO_NAME=eva-client-go

WORKDIR /go/src/github.com/Workiva/eva-client-go

COPY . /go/src/github.com/Workiva/eva-client-go
RUN ./scripts/ci/print_env.sh

RUN echo "Running go format check." && \
    ./scripts/go/check_code.sh && \
    echo "Done."

RUN echo "Running tests." && \
    ./scripts/go/test.sh && \
    echo "Done."

RUN echo "Getting code coverage report." && \
    ./scripts/codecov/report.sh && \
    echo "Done."

RUN echo "Building eva-client-go." && \
    go build . && \
    echo "Done."

FROM scratch
