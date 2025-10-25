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

