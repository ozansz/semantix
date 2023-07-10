package main

import (
	"fmt"

	"github.com/oklog/ulid/v2"
)

func main() {
	fmt.Println(ulid.Make())
}
