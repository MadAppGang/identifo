#!/bin/bash
# This script pulls latest version of the Identifo user panel and then builds it.

cd "$(dirname "$0")"
cd identifo.js
npm i
npm run build
npm link
cd ../web-element
npm i
npm link @identifo/identifo-auth-js
npm run build
rm -rf ../../static/web/element
cp -r ./dist ../../static/web/element