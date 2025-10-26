BINDIR=${CURDIR}/bin
LOCAL_BIN=${CURDIR}/bin
PACKAGE=route256/comments

PROTO_FILES_PATH := "api"

MAKEFILE_DIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))

include $(MAKEFILE_DIR)/make/Makefile.proto
include $(MAKEFILE_DIR)/make/Makefile.migrate

bindir:
	mkdir -p ${BINDIR}

build: bindir
	echo "build comments"

test:
	go test ./internal/... -count=1

test-cover:
	go test ./internal/... -count=1 -cover

test-race:
	go test ./internal/... -count=1 -race
	
minimock-gen:
	mkdir -p mocks
	$(MAKE) minimock-gen-common

minimock-gen-win:
	mkdir mocks 2>nul || exit 0
	$(MAKE) minimock-gen-common

minimock-gen-common:
	minimock -i route256/comments/internal/service.commentRepository -o ./mocks/ -s "_mock.go"
