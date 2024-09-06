NAME := tcloud-sdk/cli
ARTIFACTS := \
		bin/tcloud

GOBIN := ${GOPATH}/src/tcloud-sdk/cli/bin

INSTALL_BIN := ${GOBIN}/tcloud

.PHONY : install build clean

build : ${ARTIFACTS}

${ARTIFACTS} : main.go
		mkdir -p bin
		go build -o bin/tcloud tcloud-sdk/cli

${INSTALL_BIN} : ${ARTIFACTS}

install : ${INSTALL_BIN}
		cp ${INSTALL_BIN} ${GOPATH}/bin/

# uninstall : 
# 		rm -f ${INSTALL_BIN}

clean :
		rm -rf bin
		rm -f ${GOPATH}/bin/tcloud

		cd ${GOPATH}/src/tcloud-sdk/example && \
		for x in MXNet PyTorch TensorFlow helloworld; do \
			rm -rf "$$x"/configurations; \
			rm -f "$$x"/run.sh; \
		done
