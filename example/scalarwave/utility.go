package main

import (
	"fmt"
)

var pr = true

/* some helpers to see which function is called */
func trace(s string) string {
	if pr == true {
		fmt.Println("entering:", s)
	} else {
		fmt.Print("")
	}
	return s
}

func un(s string) {
	if pr == true {
		fmt.Println("leaving: ", s)
	} else {
		fmt.Print("")
	}
}
