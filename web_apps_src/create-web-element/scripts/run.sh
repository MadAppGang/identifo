#!/bin/sh
if test -z "$1" 
then
    echo ERROR: Please specify a directory to start a new web element
    exit 1
fi
mkdir $1
cd $1
git init
git remote add upstream https://github.com/MadAppGang/identifo.git
git sparse-checkout init
git sparse-checkout set "web_apps_src/web-element"
git pull upstream master
shopt -s dotglob
mv web_apps_src/web-element/* ./
rm -r web_apps_src
rm -rf .git
git init
npm install @identifo/identifo-auth-js@latest --save