package main

import (
	"fmt"
	"os"
	"tetrice/internal"
)

func main() {

	fmt.Println("read the data and check args")
	args:= os.Args
	if len(args)!=2{
		fmt.Println("the args are not set propaly")
		return
	}
	files,err:=os.ReadFile(args[1])
	if err!=nil{
fmt.Println("error reading files")
	}
	fileContent:=string(files)
	tetrice:=internal.Validate(fileContent)
fmt.Println("tet",tetrice)
	fmt.Println("validate the tetrice ")
	fmt.Println("optimize the tetrice")
}