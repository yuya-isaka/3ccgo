package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	// a := "b"
	// var lowerA rune = rune('a')
	// for _, v := range a {
	// 	fmt.Println((v - lowerA + 1) * 8)
	// }

	// fmt.Println("-------------------")

	// for _, a := range os.Args[1] {
	// 	fmt.Println(string([]rune{a}[0]))
	// 	fmt.Println(string(a))
	// 	fmt.Println("-------------------")
	// }

	var buffer bytes.Buffer
	// file, err := os.Create("test.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// io.WriteString(file, "Hello")
	buffer.Write([]byte("Hello World"))
	fmt.Println(buffer.String())
	io.WriteString(&buffer, "world")
	fmt.Println(buffer.String())

	request, err := http.NewRequest("GET", "http://ascii.jp", nil)
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("X-TEST", "test")
	request.Write(os.Stdout)
}
