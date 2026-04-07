package internal

import (
	"errors"
	"fmt"
	"strings"
)

func Validate(text string) ([][][]rune, error) {

	tetros, err := parseTetros(text)
	if err != nil {
		return nil, err
	}

	return tetros, nil
}

func parseTetros(input string) ([][][]rune, error) {
	rawBlocks := strings.Split(strings.TrimSpace(input), "\n\n")

	var tetros [][][]rune
	var valid bool

	tetros, valid = isValidTetrimino(rawBlocks)
	if !valid {
		return nil, errors.New("invalid tetrimino found")
	}

	return tetros, nil
}

func isValidTetrimino(block []string) ([][][]rune, bool) {
	tetros := make([][][]rune, 0)
	for _, tetro := range block {
		tetro = strings.TrimSpace(tetro)
		lines := strings.Split(tetro, "\n")
		if len(lines) != 4 {

			return nil, false
		}
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) != 4 {

				return nil, false
			}
		}
		//convert the tetro to matrice
		matrix, err := blockToMatrix(strings.Join(lines, "\n"))
		if err != nil {
			return nil, false
		}
		tetros = append(tetros, matrix)

		count := 0
		connections := 0

		// Iterate over each cell
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {

				// Validate character
				if matrix[i][j] != '#' && matrix[i][j] != '.' {
					return nil, false
				}

				// Process only '#'
				if matrix[i][j] == '#' {
					count++

					if i > 0 && matrix[i-1][j] == '#' {
						connections++
					}
					if i < 3 && matrix[i+1][j] == '#' {
						connections++
					}
					if j > 0 && matrix[i][j-1] == '#' {
						connections++
					}
					if j < 3 && matrix[i][j+1] == '#' {
						connections++
					}
				}
			}
		}

		if !(count == 4 && (connections == 6 || connections == 8)) {
			return nil, false
		}

	}
	return tetros, true

}

func blockToMatrix(block string) ([][]rune, error) {
	// Trim extra spaces/newlines
	block = strings.TrimSpace(block)

	// Split into lines
	lines := strings.Split(block, "\n")
	if len(lines) != 4 {
		return nil, fmt.Errorf("expected 4 lines, got %d", len(lines))
	}

	matrix := make([][]rune, 4)
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) != 4 {
			return nil, fmt.Errorf("line %d has length %d, expected 4", i, len(line))
		}
		matrix[i] = []rune(line)
	}

	return matrix, nil
}
