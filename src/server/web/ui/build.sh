#!/bin/sh
#BASEDIR=$(dirname $(realpath $0))

BASEDIR="."

cd "${BASEDIR}/react" || exit
pwd
nvm use
#nvm use
yarn run build
sed -i'.backup' 's/<!--PRODUCTION_BUILD-->/<script src="build.js"><\/script>/' ./dist/index.html
cd .. || exit

#webpack://injector/./node_modules/@elastic/eui/es/components/icon/assets/_lazy_^\.\/.*$_chunkName:_icon.%5Brequest%5D_namespace_object?:
