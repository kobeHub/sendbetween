# The binary to build
BIN ?= sendbetween
PKG := github.com/kobeHub/sendbetween

# Get arch info from `go env` if not spedficed
local : ARCH ?= $(shell go env GOOS)-$(shell go env GARCH)
ARCH ?= linux-amd64	

# env variables
platform_temp = $(subst -, ,${ARCH})
GOOS = $(word 1, ${platform_temp})
GOARCH = $(word 2, ${platform_temp})
GOPROXY ?= "https://goproxy.io,direct"                                                                                                                                                                                
MKDIR_P = mkdir -p 

# build output
OUTPUT_DIR = $(shell pwd)/bin/${GOOS}/${GOARCH}
OUTPUT = ${OUTPUT_DIR}/${BIN}
ifeq ("${GOOS}", "windows")
	OUTPUT := ${OUTPUT}.exe
endif

symbol_file := $(shell pwd)/${BIN}


.PHONY: hello build clean 

hello:
	echo ${GOOS}, ${GOARCH}, ${GOPROXY}, ${BIN}, ${OUTPUT}

build:
	@$(MAKE) clean-all
	@$(MAKE) test
	go build -o ${OUTPUT} \
		${PKG}/cmd/${BIN}
	if [ "${GOOS}" = "linux" ]; then rm -f ${symbol_file} && ln -s ${OUTPUT} ${symbol_file}; fi

test:
	go test -v ./... -short

clean:
	rm -rf $(shell pwd)/bin
	rm ${symbol_file}

clean-all:
	go clean -i ./...

fmt:
	go fmt ./...

run:
	@$(MAKE) build && ${OUTPUT}
	
