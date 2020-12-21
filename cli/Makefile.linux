NAME := tcloud-sdk/cli
ARTIFACTS := \
		bin/tcloud

GOBIN := ${GOPATH}/src/tcloud-sdk/cli/bin

INSTALL_BIN := ${GOBIN}/tcloud

.PHONY : install build clean

build : ${ARTIFACTS}

${ARTIFACTS} : main.go
		mkdir -p bin
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -i -o bin/tcloud tcloud-sdk/cli

${INSTALL_BIN} : ${ARTIFACTS}

install : ${INSTALL_BIN}
		cp ${INSTALL_BIN} ${GOPATH}/bin/

# uninstall : 
# 		rm -f ${INSTALL_BIN}

clean :
		rm -rf bin
		rm -f ${GOPATH}/bin/tcloud
		rm -rf configurations
		rm -f run.sh
		rm -rf ${GOPATH}/src/tcloud-sdk/examples/TuXiv_example/configurations
