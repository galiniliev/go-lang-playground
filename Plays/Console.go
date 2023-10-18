package main

// Import resty into your code and refer it as `resty`.
import (
	"flag"
	"fmt"
)

func main() {

	wordPtr := flag.String("word", "default value", "a string for description")
	flag.Parse()
	fmt.Println("word:", *wordPtr)

}
