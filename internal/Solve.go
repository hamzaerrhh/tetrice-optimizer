package internal

import (
	"errors"
	"math"
	"sort"
	"strings"
)

type Piece struct {
	mask   uint64
	width  int
	height int
	letter rune
}

func Solve(input [][][]rune) (string, error) {
	if len(input) == 0 {
		return "", errors.New("no pieces provided")
	}

	pieces := make([]Piece, len(input))
	for i, grid := range input {
		p, err := extractPiece(grid, rune('A'+i))
		if err != nil {
			return "", err
		}
		pieces[i] = p
	}

	// Minimum possible square area
	minArea := len(pieces) * 4
	size := int(math.Ceil(math.Sqrt(float64(minArea))))

	for size <= 16 {
		board := make([]uint16, size) // Each row is a bitmask
		sort.Slice(pieces, func(i, j int) bool {
			return pieces[i].height*pieces[i].width > pieces[j].height*pieces[j].width
		})
		if backtrack(board, pieces, 0, size) {
			return renderBoard(board, pieces, size), nil
		}
		size++
	}

	return "", errors.New("cannot solve")
}

// Backtrack using "First Empty Slot" strategy
type state struct {
	p    Piece
	r, c int
}

var solution []state

func backtrack(board []uint16, pieces []Piece, index int, size int) bool {
	if index == len(pieces) {
		return true
	}

	p := pieces[index]

	for r := 0; r <= size-p.height; r++ {
		for c := 0; c <= size-p.width; c++ {
			if canPlace(board, p, r, c) {
				place(board, p, r, c)
				solution = append(solution, state{p, r, c})
				if backtrack(board, pieces, index+1, size) {
					return true
				}
				solution = solution[:len(solution)-1] // undo
				remove(board, p, r, c)
			}
		}
	}
	return false
}
func canPlace(board []uint16, p Piece, r, c int) bool {
	// Check if any of the 4 rows the piece occupies conflict
	pMask := p.mask
	for i := 0; i < p.height; i++ {
		rowMask := uint16((pMask>>(i*4))&0xF) << c
		if board[r+i]&rowMask != 0 {
			return false
		}
	}
	return true
}

func place(board []uint16, p Piece, r, c int) {
	for i := 0; i < p.height; i++ {
		rowMask := uint16((p.mask>>(i*4))&0xF) << c
		board[r+i] |= rowMask
	}
}

func remove(board []uint16, p Piece, r, c int) {
	for i := 0; i < p.height; i++ {
		rowMask := uint16((p.mask>>(i*4))&0xF) << c
		board[r+i] &= ^rowMask
	}
}

// extractPiece normalizes the # shape into a 4x4 bitmask
func extractPiece(grid [][]rune, letter rune) (Piece, error) {
	var points [][2]int
	minR, minC, maxR, maxC := 4, 4, -1, -1

	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if grid[r][c] == '#' {
				points = append(points, [2]int{r, c})
				if r < minR {
					minR = r
				}
				if r > maxR {
					maxR = r
				}
				if c < minC {
					minC = c
				}
				if c > maxC {
					maxC = c
				}
			}
		}
	}

	if len(points) != 4 {
		return Piece{}, errors.New("each piece must have exactly 4 '#' blocks")
	}

	var mask uint64
	for _, pt := range points {
		// Normalize to (0,0) and store in 4-bit row chunks
		r, c := pt[0]-minR, pt[1]-minC
		mask |= 1 << (r*4 + c)
	}

	return Piece{
		mask:   mask,
		width:  maxC - minC + 1,
		height: maxR - minR + 1,
		letter: letter,
	}, nil
}

func renderBoard(boardMasks []uint16, pieces []Piece, size int) string {
	// This is only called once at the end, so it doesn't need bitmask speed
	res := make([][]rune, size)
	for i := range res {
		res[i] = []rune(strings.Repeat(".", size))
	}

	// We need to re-run the placement logic to know which letter goes where
	// (Or store the solution state during backtrack)
	type state struct {
		p    Piece
		r, c int
	}
	solution := []state{}

	// Helper to find the solution coordinates again
	var solveFinal func(int) bool
	currentBoard := make([]uint16, size)
	solveFinal = func(idx int) bool {
		if idx == len(pieces) {
			return true
		}
		p := pieces[idx]
		for r := 0; r <= size-p.height; r++ {
			for c := 0; c <= size-p.width; c++ {
				if canPlace(currentBoard, p, r, c) {
					place(currentBoard, p, r, c)
					if solveFinal(idx + 1) {
						solution = append(solution, state{p, r, c})
						return true
					}
					remove(currentBoard, p, r, c)
				}
			}
		}
		return false
	}

	solveFinal(0)
	for _, s := range solution {
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				if (s.p.mask>>(i*4+j))&1 == 1 {
					res[s.r+i][s.c+j] = s.p.letter
				}
			}
		}
	}

	var sb strings.Builder
	for _, row := range res {
		sb.WriteString(string(row) + "\n")
	}
	return sb.String()
}
