# product version control
VERSION=0.1.0
VCS_COMMIT=$(shell git --no-pager log --pretty=format:"%H" -1)
BUILD_TIME=$(shell date --rfc-3339=seconds)
VCS_DIRTY=$(shell if [ "$$(git status --porcelain | wc -l)" -gt 0 ]; then echo -dirty; fi;)

# third party libraries
LIBPCAP_PATH=${HOME}/cgo/linux/amd64
LIBPCAP_INCLUDE=${LIBPCAP_PATH}/include
LIBPCAP_LIB=${LIBPCAP_PATH}/lib

# cgo options
CGO_CFLAGS=-I${LIBPCAP_INCLUDE}
CGO_LDFLAGS=-L${LIBPCAP_LIB}
CGO_CC=gcc
CGO_CXX=g++
CGO_ENABLED=1

# go cross compile options
GO_GOOS=linux
GO_GOARCH=amd64

# enable static link
CGO_STATIC_LINK=1
STATIC_LD_FLAGS=

ifdef CGO_STATIC_LINK
STATIC_LD_FLAGS=\
	-extldflags=-static \
	-linkmode=external
endif

# linker options
LDFLAGS=\
	${STATIC_LD_FLAGS} \
	-s \
	-w \
	-X "github.com/nnnewb/nf/internal/constants.BUILD_COMMIT=${VCS_COMMIT}${VCS_DIRTY}" \
	-X "github.com/nnnewb/nf/internal/constants.BUILD_TIME=${BUILD_TIME}" \
	-X "github.com/nnnewb/nf/internal/constants.BUILD_STATIC=${CGO_STATIC_LINK}" \
	-X "github.com/nnnewb/nf/internal/constants.VERSION=${VERSION}"


# compiler environment variables
GO_ENV= \
	CGO_ENABLED=${CGO_ENABLED} \
	CGO_CFLAGS=${CGO_CFLAGS} \
	CC=${CGO_CC} \
	CXX=${CGO_CXX} \
	CGO_LDFLAGS=${CGO_LDFLAGS} \
	GOOS=${GO_GOOS} \
	GOARCH=${GO_GOARCH}

all: nf

.PHONY: nf
nf:
	mkdir -p ./bin
	${GO_ENV} go build -ldflags '${LDFLAGS}' -o ./bin/nf main.go
