package main

import (
	"fmt"
  "github.com/jasonzou/gomarc21/pkg/marc"
)

func main() {
	tag, err := NewTagByStr("123")
	fmt.Print(tag)
}
