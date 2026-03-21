package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func Validate(text string) [][]int {

	tetrimonios,err:=parseTetros(text)
if err!=nil{
	fmt.Println("err:",err)
	os.Exit(0)
}
if !ValidateTetros(tetrimonios){
fmt.Println("error not validate")
	os.Exit(0)
}

return nil
}

func parseTetros(input string) ([][]string, error) {
	rawBlocks := strings.Split(strings.TrimSpace(input), "\n\n")
	var result [][]string

	for _, block := range rawBlocks {
		lines := strings.Split(block, "\n")

		if len(lines) != 4 {
			return nil, errors.New("invalid block: must have 4 lines")
		}

		// Check each line length
		for _, line := range lines {
			if len(line) != 4 {
				return nil, errors.New("invalid block: each line must be 4 characters")
			}
			//check the caractere in lines
			for _,char:=range line{
				if !(char=='#' || char =='.'){
									return nil, errors.New("invalid charactere: the char must be . or # ")

				}
			}

		}

		result = append(result, lines)
	}

	return result, nil
}
func ValidateTetros(blocks [][]string)bool{
	for _, block := range blocks {
		if !isValidTetrimino(block){
return false
		}
		
	}


return true
}
func isValidTetrimino(block []string) bool {
	count := 0
	connections := 0

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if block[i][j] == '#' {
				count++

				// check neighbors
				if i > 0 && block[i-1][j] == '#' {
					connections++
				}
				if i < 3 && block[i+1][j] == '#' {
					connections++
				}
				if j > 0 && block[i][j-1] == '#' {
					connections++
				}
				if j < 3 && block[i][j+1] == '#' {
					connections++
				}
			}
		}
	}

	// must have exactly 4 blocks
	if count != 4 {
		return false
	}

	// valid connection counts
	return connections == 6 || connections == 8
}
