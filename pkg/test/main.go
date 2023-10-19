package main

import (
	"fmt"
	"go_app/pkg/test/helper"
)

func main() {
	fmt.Print("\n\ntest package order - generate fake order data\n\n")
	helper.PrintFakeData()
	fmt.Print("\n\ntest cache package\n\n")
	helper.CacheTest()
	fmt.Print("\n\ntest database package - generate fake order, insert to database, then select and print it\n\n")
	helper.DatabaseTest()
}
