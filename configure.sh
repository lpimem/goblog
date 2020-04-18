INSTANCE_NAME=goblog_prod
SITE_TITLE="Lei's Blog"
WORK_DIR=/usr/local/goblog
DOC_DIR=$WORK_DIR/doc
BUILD_DIR=/tmp
# IP and port this blog will be served at
HOST=127.0.0.1
PORT=9000
BUILD_MODE=prod

# generate a random secret string to sign the cookies
# replace with your own secret key if you have one
APPSEC=`env LC_CTYPE=C tr -dc "a-zA-Z0-9-_\$\?" < /dev/urandom | fold -w 64 | head -n 1`