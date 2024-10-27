// Package pkg1 is testdata for multichecker.
package main

import (
	"os"
)

//lint:ignore U1000 Ignore unused function
func main() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	os.Exit(1) // want "call to os.Exit in main function is forbidden"

}
