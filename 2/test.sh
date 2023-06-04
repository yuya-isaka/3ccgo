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

test '  10 ;' 10
test '10  +   20;' 30
test '30  -  20 +3 - 3 ;   ' 10
test '   10 -3  +20   + 30-2;' 55
test '   10 *3  +20   + 30-2;' 78
test '10/2 * 3;' 15
test '3 + 2 * 3;' 9
test '10 /2 + 3 * 3 - 3-1/1;' 10
test '10; 20;' 20
test '10+20-3; 10-3;' 7
test '-3+10;' 7
test '- - - -10;' 10
test '-3 + 20 - 3 +4-3---4; ++10 - - 3 - 10 + + +20 - - - - 4;' 27
test '1 == 1;' 1
test '1 == 0;' 0
test '1 != 0;' 1
test '1 != 1;' 0
test '1 < 10;' 1
test '1>10;' 0
test '10 < 10;' 0
test '10 > 10;' 0
test '1 <= 10;' 1
test '10 <= 10;' 1
test '11 <= 10;' 0
test '1 >= 10;' 0
test '10 >= 10;' 1
test '11 >= 10;' 1
test '(1 > 0) + 1;' 2
test '(3 + 2) * 5;' 25
test '(2+2) * ( 3-1) + (5-3) *2;' 12
test 'a=3; a;' 3
test 'a=3; b=4; a+b;' 7
test 'a=3; b=4; c=5; a+b-c;' 2

echo OK