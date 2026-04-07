package main

import (
	"fmt"
	"os"
	"strings"
	"tetrice/internal"
)

func main() {

	args := os.Args
	if len(args) != 2 {
		fmt.Println("the args are not set propaly")
		return
	}
	files, err := os.ReadFile(args[1])
	if err != nil {
		fmt.Println("ERROR")
	}
	fileContent := string(files)
	fileContent = strings.ReplaceAll(fileContent, "\r\n", "\n")
	tetros, errValidation := internal.Validate(fileContent)
	if errValidation != nil {
		fmt.Println("ERROR")
		return
	}
	result, errSolving := internal.Solve(tetros)
	if errSolving != nil {
		fmt.Print("ERROR")
		return
	}
	fmt.Print(result)
}
