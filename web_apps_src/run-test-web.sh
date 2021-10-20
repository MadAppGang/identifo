#!/bin/bash
# This script pulls latest version of the Identifo user panel and then builds it.

cd "$(dirname "$0")"
./update-web.sh
cd ..
make run_boltdb &
IDENTIFO_PID=$!
cd web_apps_src/web-element
npx cypress run
status=$?
kill $IDENTIFO_PID
exit "$status"