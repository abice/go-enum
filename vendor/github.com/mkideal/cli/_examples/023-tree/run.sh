#!/bin/bash

go build -o app
./app cmd1
./app cmd1 cmd11
./app cmd1 cmd12

./app cmd2
./app cmd2 cmd21
./app cmd2 cmd22
./app cmd2 cmd22 cmd221
./app cmd2 cmd22 cmd222
./app cmd2 cmd22 cmd223
rm app
