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

# linker options
LDFLAGS=-extldflags=-static -linkmode=external -s -w

# version control
VCS_COMMIT=$(git log --pretty=format:"%H" -1)

# compiler environment variables
GO_ENV= \
	CGO_ENABLED=${CGO_ENABLED} \
	CGO_CFLAGS=${CGO_CFLAGS} \
	CC=${CGO_CC} \
	CXX=${CGO_CXX} \
	CGO_LDFLAGS=${CGO_LDFLAGS} \
	GOOS=${GO_GOOS} \
	GOARCH=${GO_GOARCH}

all: nf deps

.PHONY: nf
nf:
	mkdir -p ./bin
	${GO_ENV} go build -ldflags '${LDFLAGS}' -o ./bin/nf main.go
