#!/bin/bash
# This script pulls latest version of the Identifo user panel and then builds it.

cd "$(dirname "$0")"
# Fetch and build source code.
wget https://github.com/MadAppGang/identifo-web-elements/archive/main.zip
unzip main.zip
cd identifo-web-elements-main
npm i
npm run build

# Update build directory content.
cd ../
rm -rf ./build/element;
cp -r ./identifo-web-elements-main/dist ./build/element


# Clean up.
rm -f main.zip
rm -fr identifo-web-elements-main