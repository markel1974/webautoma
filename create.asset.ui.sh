#!/bin/sh

export NVM_DIR=$HOME/.nvm
# shellcheck disable=SC2039
source "${NVM_DIR}/nvm.sh"

ROOT=$(pwd)
cd ./src/server/web/ui/react || exit
nvm use

yarn run build
sed -i'.backup' 's/<!--PRODUCTION_BUILD-->/<script src="build.js"><\/script>/' ./dist/index.html
cd "${ROOT}" || exit
./build/webautoma -t web:./src/server/web/content:./src/server/web/ui/react/dist
