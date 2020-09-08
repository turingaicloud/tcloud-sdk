NAME := tcloud-sdk/cli
ARTIFACTS := \
		bin/tcloud

INSTALL_BIN := ${GOBIN}/tcloud

.PHONY : install build clean

build : ${ARTIFACTS}

${ARTIFACTS} : main.go
		mkdir -p bin
		go build -i -o bin/tcloud tcloud-sdk/cli
		# GOOS=linux		GOARCH=amd64 go build -i -o bin/tcloud-Linux-x86_64 tcloud-sdk/cli

${INSTALL_BIN} : ${ARTIFACTS}
		cp bin/tcloud ${GOBIN}/tcloud

# install : ${INSTALL_BIN}

# uninstall : 
# 		rm -f ${INSTALL_BIN}

clean :
		rm -rf bin