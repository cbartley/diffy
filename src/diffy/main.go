package main

import (
	"fmt"
	"diffy/diff"
)

func main() {
	fmt.Println("diffy!")
	_ = diff.LevenshteinDistance_v2("foo", "bar")
}
