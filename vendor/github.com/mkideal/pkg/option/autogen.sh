#!/bin/bash

set -e

types='
bool
byte
rune
string
int
int8
int16
int32
int64
uint
uint8
uint16
uint32
uint64
'

cat >option.go<<EOF
package option
EOF

for T in $types
do
	upperT=`echo $T | awk '{for(i=1;i<=NF;i++) {printf "%s%s ", toupper(substr($i,1,1)),substr($i,2)};printf ORS}'`
	Ts="$T"s
	cat >>option.go<<EOF

func $upperT(dft $T, $Ts ...$T) $T {
	if len($Ts) == 0 {
		return dft
	}
	return $Ts[0]
}
EOF
done

cat >>option.go<<EOF

func Interface(dft interface{}, interfaces ...interface{}) interface{} {
	if len(interfaces) == 0 {
		return dft
	}
	return interfaces[0]
}
EOF

gofmt -s -w option.go
