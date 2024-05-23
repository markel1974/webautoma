#!/bin/sh

ROOT=$(pwd)
cd ./src/server/web/ui/react || exit
pwd
#nvm use
yarn run build
sed -i'.backup' 's/<!--PRODUCTION_BUILD-->/<script src="build.js"><\/script>/' ./dist/index.html
cd "${ROOT}" || exit
./build/webautoma -t web:./src/server/web/content:./src/server/web/ui/react/dist
