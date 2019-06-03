#!/bin/bash
# This script pulls latest version of the Identifo user panel and then builds it.

# Fetch and build source code.
wget https://github.com/MadAppGang/identifo-admin/archive/master.zip
tar xvf master.zip
cd identifo-admin-master
npm i
npm run build

# Update build directory content.
rm -rf ../build
mv build/ ../

# Clean up.
cd ../
rm -f develop.zip
rm -fr identifo-admin-master