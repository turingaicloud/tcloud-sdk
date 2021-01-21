#!/bin/bash

VERSION="0.1.1"

MAC_MAKEFILE="Makefile.mac"
LINUX_MAKEFILE="Makefile.linux"
WINDOWS_MAKEFILE="Makefile.windows"

QUICKSTART_PATH="/Users/xcwan/Desktop/Workspace/Project/quickstart"
TCLOUD_PATH=$PWD

# Makedir
mkdir -p ${QUICKSTART_PATH}/bin/linux-amd64-v${VERSION}
mkdir -p ${QUICKSTART_PATH}/bin/windows-amd64-v${VERSION}
mkdir -p ${QUICKSTART_PATH}/bin/macos-amd64-v${VERSION}

cd cli

cat <<EOF >>setup.sh
#!/bin/bash

DIR=\$PWD

mkdir -p \$HOME/.tcloud

if [ ! -f \$BASH ]; then
    touch \$BASH
fi

echo "Remember to execute the following command:"
echo "export PATH=\$DIR:\$PATH"
EOF

# Linux
make -f ${LINUX_MAKEFILE} clean && make -f ${LINUX_MAKEFILE} install
cp bin/tcloud ${QUICKSTART_PATH}/bin/linux-amd64-v${VERSION}
cp setup.sh ${QUICKSTART_PATH}/bin/linux-amd64-v${VERSION}

# Windows
make -f ${WINDOWS_MAKEFILE} clean && make -f ${WINDOWS_MAKEFILE} install
cp bin/tcloud ${QUICKSTART_PATH}/bin/windows-amd64-v${VERSION}

# Mac
make -f ${MAC_MAKEFILE} clean && make -f ${MAC_MAKEFILE} install
cp bin/tcloud ${QUICKSTART_PATH}/bin/macos-amd64-v${VERSION}
cp setup.sh ${QUICKSTART_PATH}/bin/macos-amd64-v${VERSION}

rm setup.sh

cd ${QUICKSTART_PATH}/bin
zip -r linux-amd64-v${VERSION}.zip linux-amd64-v${VERSION}
zip -r windows-amd64-v${VERSION}.zip windows-amd64-v${VERSION}
zip -r macos-amd64-v${VERSION}.zip macos-amd64-v${VERSION}

cd ${TCLOUD_PATH}
make clean
