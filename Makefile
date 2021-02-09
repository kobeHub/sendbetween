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

symbol_file := $(shell pwd)/app

hello:
	echo ${GOOS}, ${GOARCH}, ${GOPROXY}, ${BIN}, ${OUTPUT}

build: 
	go build -o ${OUTPUT} \
		${PKG}/cmd/${BIN}
	rm -f ${symbol_file}
	ln -s ${OUTPUT} ${symbol_file} 

clean:
	rm -rf $(shell pwd)/bin
	rm ${symbol_file}

run:
	@$(MAKE) build && ${OUTPUT}
