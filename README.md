# JamVote

JamVote is an web-application for managing GameJams and voting.

# Local Setup Guide

Install [Google App Engine Standard](https://cloud.google.com/appengine/docs/standard/go/).

```
cd $GOPATH/src/github.com/adinfinit
git clone git@github.com:adinfinit/jamvote.git
cd jamvote
go get -u ./...
dev_appserver.py appengine
```