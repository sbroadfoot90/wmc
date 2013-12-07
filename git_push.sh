#!/bin/bash

args=("$@")

git pull
git add *
git commit -m "${args[0]}"
git push master origin
