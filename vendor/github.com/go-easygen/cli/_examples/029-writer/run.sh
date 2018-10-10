#!/bin/bash

set -e

go build -o app

# write to stdout
./app

# write to file
./app -w 1.txt
cat 1.txt

rm app
rm 1.txt
