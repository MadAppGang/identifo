#!/bin/bash
# This script pulls latest version of the Identifo user panel and then builds it.

cd "$(dirname "$0")"
cd admin
npm i
export BASE_URL=/adminpanel/ && export ASSETS_PATH=/adminpanel/ # Needed for build.
npm run build
rm -rf ../../static/admin_panel
cp -r ./build ../../static/admin_panel