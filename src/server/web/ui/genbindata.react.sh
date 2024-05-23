#!/bin/sh
#BASEDIR=$(dirname $(realpath $0))

BASEDIR="."

cd "${BASEDIR}/ui.react" || exit
yarn run build
sed -i'.backup' 's/<!--PRODUCTION_BUILD-->/<script src="build.js"><\/script>/' ./dist/index.html
cd .. || exit

./go-bindata -pkg "web" -prefix "ui.react/dist/" -o static_assets.go ./ui.react/dist/

#VER=$(git rev-parse --short HEAD)
#VER="cb7b899"
#sed -i '' -E "s/(build\.js\?v=).{7}/\1${VER}/g" "${BASEDIR}/ui.react/index.html"
#TODO in ui/index.html --> version.Version

#webpack://injector/./node_modules/@elastic/eui/es/components/icon/assets/_lazy_^\.\/.*$_chunkName:_icon.%5Brequest%5D_namespace_object?:
