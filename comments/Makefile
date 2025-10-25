BINDIR=${CURDIR}/bin
LOCAL_BIN=${CURDIR}/bin
PACKAGE=route256/comments

PROTO_FILES_PATH := "api"

include make/Makefile.proto
include make/Makefile.migrate

bindir:
	mkdir -p ${BINDIR}

build: bindir
	echo "build comments"


run-migrations:
	echo "run migrations"
