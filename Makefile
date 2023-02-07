UNAME := $(shell uname)
CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
LINTVER=latest
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}

run-upload-client:
	go run cmd/uploadclient/main.go

run-listclient:
	go run cmd/listclient/main.go

run-downloadclient:
	go run cmd/downloadclient/main.go

run-server:
	go run cmd/server/main.go

generate-install:
	go install \
	github.com/golang/protobuf/protoc-gen-go \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	protoc -I ./api \
	--go_out ./pkg --go_opt paths=source_relative \
	--go-grpc_out ./pkg --go-grpc_opt paths=source_relative \
	./api/fileserver/fileserver.proto

lint: install-lint
	${LINTBIN} run


install-lint: bindir
	test -f ${LINTBIN} || \
  (GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
  mv ${BINDIR}/golangci-lint ${LINTBIN})

bindir:
	mkdir -p ${BINDIR}

lintfix:
    ifeq ($(UNAME), Linux)
		find . \( -path './cmd/*' -or -path './internal/*' -or -path './pkg/*' -or -path './e2e/*' \) \
		-type f -name '*.go' -print0 | \
		xargs -0  sed -i '/import (/,/)/{/^\s*$$/d;}'
    endif
    ifeq ($(UNAME), Darwin)
		find . \( -path './cmd/*' -or -path './internal/*' -or -path './pkg/*' -or -path './e2e/*' \) \
		-type f -name '*.go' -print0 | \
		xargs -0  sed -i '' '/import (/,/)/{/^\s*$$/d;}'
    endif
	goimports -local=github.com/SerjZimmer/testovoe1 -w ./cmd ./internal ./pkg

result: run-upload-client run-listclient