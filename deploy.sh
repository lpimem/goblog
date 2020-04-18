#! /bin/bash

# configurations
GOPATH=$HOME/go
INSTANCE_NAME=goblog_prod
SITE_TITLE="My Blog"
WORK_DIR=$HOME/$INSTANCE_NAME
DOC_DIR=$WORK_DIR/doc
BUILD_DIR=/tmp

# IP and port this blog will be served at
HOST=127.0.0.1
PORT=9000

MODE=prod
if [[ $# -gt 0 ]]; then
    MODE = $1
fi

# generate a random secret string to sign the cookies
# replace with your own secret key if you have one
APPSEC=`env LC_CTYPE=C tr -dc "a-zA-Z0-9-_\$\?" < /dev/urandom | fold -w 64 | head -n 1`

sudo apt-get -y install git mercurial

go get github.com/microcosm-cc/bluemonday
go get github.com/russross/blackfriday
go get github.com/bradfitz/slice
go get github.com/revel/revel
go get github.com/revel/cmd/revel
export PATH=$PATH:$GOPATH/bin

go get github.com/lpimem/goblog

cd $GOPATH/src
mv github.com/lpimem/goblog ./$INSTANCE_NAME
cd ./$INSTANCE_NAME
find . -name '*.go' | xargs sed -i s#\"goblog/#\"$INSTANCE_NAME/#
sed -i s#/goblog#/$INSTANCE_NAME# conf/app.conf
arg=(-vsite_title="$SITE_TITLE")
gawk "${arg[@]}" -i inplace '\
/goblog.base_dir = .+/ {print "goblog.base_dir = '$WORK_DIR'" ;next};\
/goblog.doc_base_dir = .+/ {print "goblog.doc_base_dir = '$DOC_DIR'" ;next};\
/app.secret =[ .]*/ {print "app.secret = '$APPSEC'" ;next}\
/goblog.site_title = [ .]*/ {print "goblog.site_title =", site_title;next}\
/http.addr = [ .]*/ {print "http.addr = '$HOST'" ;next}\
/http.port = [ .]*/ {print "http.port = '$PORT'" ;next}\
{print}' conf/app.conf

revel -X-v build -m $MODE $INSTANCE_NAME $BUILD_DIR/$INSTANCE_NAME

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

echo "To start the server: $BUILD_DIR/$INSTANCE_NAME/run.sh"