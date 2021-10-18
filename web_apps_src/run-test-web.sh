#!/bin/bash
# This script pulls latest version of the Identifo user panel and then builds it.

cd "$(dirname "$0")"
./update-web.sh
cd ..
make run_boltdb &
IDENTIFO_PID=$!
cd web_apps_src/web-element
npx cypress run --record --key f1169950-dcab-42ef-ad47-f6849179dd71
kill $IDENTIFO_PID