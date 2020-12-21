#!/bin/bash

VERSION="0.1.0"

MAC_MAKEFILE="Makefile.mac"
LINUX_MAKEFILE="Makefile.linux"
WINDOWS_MAKEFILE="Makefile.windows"

QUICKSTART_PATH="/Users/xcwan/Desktop/Workspace/Project/quickstart"
TCLOUD_PATH=$PWD

cd cli

# Linux
make -f ${LINUX_MAKEFILE} clean && make -f ${LINUX_MAKEFILE} install
cp bin/tcloud ${QUICKSTART_PATH}/bin/linux-amd64-v${VERSION}

# Windows
make -f ${WINDOWS_MAKEFILE} clean && make -f ${WINDOWS_MAKEFILE} install
cp bin/tcloud ${QUICKSTART_PATH}/bin/windows-amd64-v${VERSION}

# Mac
make -f ${MAC_MAKEFILE} clean && make -f ${MAC_MAKEFILE} install
cp bin/tcloud ${QUICKSTART_PATH}/bin/macos-amd64-v${VERSION}

cd ${QUICKSTART_PATH}/bin
zip -r linux-amd64-v${VERSION}.zip linux-amd64-v${VERSION}
zip -r windows-amd64-v${VERSION}.zip windows-amd64-v${VERSION}
zip -r macos-amd64-v${VERSION}.zip macos-amd64-v${VERSION}

cd ${TCLOUD_PATH}
make clean