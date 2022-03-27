package main

import (
	"fmt"
  "github.com/jasonzou/gomarc21"
)

func main() {
	tag, err := gomarc21.NewTagByStr("123")
  if err != nil {
    fmt.Print(err)
  }
	fmt.Print(tag)
}
