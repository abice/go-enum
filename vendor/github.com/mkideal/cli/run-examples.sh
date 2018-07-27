#!/bin/bash

set -e

CWD=`pwd`
EXAMPLE_DIR="./_examples"
EXAMPLES=`ls $EXAMPLE_DIR`

for APP in $EXAMPLES
do
	cd $EXAMPLE_DIR/$APP
	SCRIPT="./run.sh"
	if [ -f "$SCRIPT" ]; then
		echo ">>>> exmaple: $APP"
		chmod +x $SCRIPT
		$SCRIPT
	fi
	cd $CWD
done
