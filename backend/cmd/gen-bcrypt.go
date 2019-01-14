package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var password = flag.String("pass", "password", "The password to create a hash from")

func main() {
	flag.Parse()
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hash))
}
