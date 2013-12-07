#!/bin/bash

args=("$@")

git pull
go fmt src/*.go
git add *
git commit -m "${args[0]}"
git push origin master
