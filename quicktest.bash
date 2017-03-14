#!/usr/bin/env bash

CWD=${PWD##*/}
TODAY=`date +%Y%m%d`
STORE=`jq .datastore.type < ~/.gzr.json`
HASH=`git rev-parse HEAD`

echo $STORE
cd $GOPATH/src/github.com/bypasslane/gzr
go install .
cd -

gzr build .
RESULTS=`gzr image get $CWD`
if [ ! $? -eq 0 ] || ! `echo $RESULTS | grep -q $HASH`; then
    exit 1
fi
DELETED=`gzr image delete $CWD:$TODAY`
if [ ! $? -eq 0 ] || ! `echo $DELETED | grep -q "Deleted 1"`; then
    exit 1
fi
