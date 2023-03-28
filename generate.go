package main

import (
	"fmt"
	"math/rand"

	validator "gopkg.in/validator.v2"
)

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func generate_random_board() []Board {
	validator.SetValidationFunc("valid_input", validate_inputs)

	solution := [5][5]int{}
	vertical := [5][2]int{}
	horizontal := [5][2]int{}
	for row := range solution {
		for col := range solution[row] {
			val := starting_poss[randInt(0, 3)]
			solution[row][col] = val
			if val == VOLTORB {
				vertical[row][1] += 1
				horizontal[col][1] += 1
			} else {
				vertical[row][0] += val
				horizontal[col][0] += val
			}
		}
	}

	b := generate_board()
	b.Vertical = vertical
	b.Horizontal = horizontal

	bs := generate_board()
	bs.Vertical = vertical
	bs.Horizontal = horizontal

	for row := range bs.Grid {
		for col := range bs.Grid[row] {
			bs.Grid[row][col].Value = solution[row][col]
		}
	}

	if errs := b.validate(); errs != nil {
		fmt.Println(errs)
	}

	return []Board{b, bs}
}
