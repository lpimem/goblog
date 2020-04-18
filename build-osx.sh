#! /bin/bash

################################################
## To install Dependencies:
##   brew install gnu-sed
##   brew install gawk
##
## To build for linux amd64:
##   GOOS=linux GOARCH=amd64 ./build-osx.sh
################################################

cd `dirname $(realpath ${0})`

if [[ ! -f configure.sh ]]; then
    echo 'Missing configure.sh'
    exit 1
fi

source configure.sh

BUILD_WORKDIR=$GOPATH/src/$INSTANCE_NAME
echo "Build workdir: $BUILD_WORKDIR"

mkdir -p $BUILD_WORKDIR
cp -R ./* $BUILD_WORKDIR/
cd $BUILD_WORKDIR

find . -name '*.go' | xargs gsed -i s#\"goblog/#\"$INSTANCE_NAME/#
gsed -i s#/goblog#/$INSTANCE_NAME# conf/app.conf
arg=(-vsite_title="$SITE_TITLE")
gawk "${arg[@]}" -i inplace '\
/goblog.base_dir = .+/ {print "goblog.base_dir = '$WORK_DIR'" ;next};\
/goblog.doc_base_dir = .+/ {print "goblog.doc_base_dir = '$DOC_DIR'" ;next};\
/app.secret =[ .]*/ {print "app.secret = '$APPSEC'" ;next}\
/goblog.site_title = [ .]*/ {print "goblog.site_title =", site_title;next}\
/http.addr = [ .]*/ {print "http.addr = '$HOST'" ;next}\
/http.port = [ .]*/ {print "http.port = '$PORT'" ;next}\
{print}' conf/app.conf

revel -X-v build -m $BUILD_MODE $INSTANCE_NAME $BUILD_DIR/$INSTANCE_NAME

rc=$?
if [[ $rc != 0 ]]; then
    echo "cannot build goblog"
    exit $rc
fi

DIR_OK=1
mkdir -p $WORK_DIR
rc=$?
if [[ $rc != 0 ]]; then
    echo "cannot make working dir, please assign a valid value to WORK_DIR"
    DIR_OK=0
fi

mkdir -p $DOC_DIR
rc=$?
if [[ $rc != 0 ]]; then
    echo "cannot make document dir, please assign a valid value to DOC_DIR"
    DIR_OK=0
fi

if [[ DIR_OK == 0 ]]; then
    exit 1;
fi
echo '- - - - - - - - - - - - - - - - - - - '
echo "Your blog articles will be scaned from
  $DOC_DIR"
echo "To start the server:
  $BUILD_DIR/$INSTANCE_NAME/run.sh"
echo '- - - - - - - - - - - - - - - - - - - '