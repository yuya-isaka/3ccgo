#!/bin/bash

test() {
	input="$1"
	expected="$2"

	./3ccgo "$input" > tmp.s || exit
	gcc -o tmp tmp.s
	./tmp
	actual="$?"

	if [ "$actual" = "$expected" ]; then
		echo "$input => $expected"
	else
		echo "$input => $expected expected, but got $actual"
	fi
}

go build -o 3ccgo main.go
echo ----------------- Build Ok --------------------

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
test '10; 20;' 20
test '10+20-3; 10;' 10
test '30-27;' 3
test '10;' 10
test '10+20-3;' 27
test '10 + 30 - 3 - 3;' 34
test '(10 + 20) * 3;' 90
test '10 + 10 * 2;' 30
test '10 / 5 + 2 * (2 - 1);' 4
test '- - - -10;' 10
test '- 20 + 30 / 3 - - 10;' 0
test '1<10;' 1
test '10 < 10;' 0
test '9 < -10;' 0
test '1>10;' 0
test '10 > 10;' 0
test '9 > -10;' 1
test '9 > 1;' 1
test '10 <= 10;' 1
test '10 >= 10;' 1
test '10 <= 9;' 0
test '9 >= 10;' 0
test '10 <= -11;' 0
test '-11 >= 10;' 0
test 'a=3; b=4; a+b;' 7
test 'a=- - -3; b = 3 + -2 + 10 * 2 / (10/2); b-a;' 8
test 'a=10; z = 3; y = a*z; y;' 30
test 'a=10; (c=3); a+c;' 13

echo -------------------- Ok -----------------------