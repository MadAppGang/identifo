#!/bin/bash
# This script pulls latest version of the Identifo user panel and then builds it.

cd "$(dirname "$0")"
# Fetch and build source code.
wget https://github.com/MadAppGang/identifo-web-static/archive/main.zip
unzip main.zip
cd identifo-web-static-main
export BASE_URL=/web/
npm i

wget https://github.com/sokolovstas/identifo.js/archive/refs/heads/api-integration.zip
unzip api-integration.zip
rm -r identifo.js
mv identifo.js-api-integration identifo.js
cd identifo.js
npm i
npm run build

cd ../
npm run build

# Update build directory content.
rm -rf ../build
mv public/ ../build

# Clean up.
cd ../
rm -f main.zip
rm -fr identifo-web-static-main
