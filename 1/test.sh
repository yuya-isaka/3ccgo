#!/bin/bash

test() {
	input="$1"
	expected="$2"

	./3ccgo "$input" > tmp.s || exit
	gcc -o tmp tmp.s
	./tmp
	actual="$?"

	if [ "$expected" = "$actual" ]; then
		echo "$input => $expected"
	else
		echo "$input => $expected expected, but got $actual"
	fi
}

go build -o 3ccgo main.go
echo Build OK

test '10' 10
test '10+20' 30
test '30-20' 10

echo OK