package main

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
)

func main() {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage("./test/sample.jpg")
	text, _ := client.Text()
	fmt.Println(text)
}
