[![Go project version](https://badge.fury.io/go/github.com%2Flpimem%2Fgoblog.svg)](https://badge.fury.io/go/github.com%2Flpimem%2Fgoblog) 
[![Go Report Card](https://goreportcard.com/badge/github.com/lpimem/goblog)](https://goreportcard.com/report/github.com/lpimem/goblog)

----

## Introduction

**GoBlog** is a blogging server to publish you markdown notes as neatly rendered articles. 

## Prerequisite

`Go`: You can install it use this [automated script for Ubuntu](https://gist.github.com/lpimem/2b33bf3b5704aab6b56541c14157f80f).

## Easy Deployment

1. Download [deploy script](https://github.com/lpimem/goblog/blob/c459198cc6097ece6b62ef090d5d5b649de404ae/deploy.sh). Or follow instructions on [Go website](https://golang.org/doc/install).

2. Edit `deploy.sh`, update `configurations` part.

3. Execute `./deploy.sh`. 

## Configurations

`DOC_DIR` in deploy.sh: specify a directory for all articles to be served. `deploy.sh` will update `conf/app.conf` and genearte the final configuration file. 

`APPSEC` is the string used to sign the website's cookies. By default, `deploy.sh` will genereate a 64-byte random string using `/dev/urandom`. 

## Managing articles

Recommended way to manage your notes locally is to create a `notes` folder and put everything inside. You can put it inside a cloud storage folder to prevent data lose, e.g., `~/Dropbox`. 

Your `notes` folder's content should look like this:

```txt
notes
  |----img
  |     |----img1.jpg
  |     |----img2.png
  |----nots1.md
  |----note2.md
  |----...
```

Then you can sync notes to server for publishing with rsync: 

```bash
rsync -a --progress ~/Dropbox/notes/*.md $REMOTE_HOST:$DOC_DIR/
rsync -a --progress ~/Dropbox/notes/img/* $REMOTE_HOST:$REMOTE_GOPATH/src/$GOBLOG_INS/public/img
```

It doesn't matter which tool to use for syncing. Just make sure notes and images go to the correct folders. 

## Developing

1. Check out code:

  ```bash
  cd $GOPATH/src
  git clone git@github.com:lpimem/goblog.git
  ```
  
2. Install dependencies:
  
  ```bash
  cd goblog
  ./dependency.sh
  ```
