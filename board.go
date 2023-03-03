package main

import (
	"fmt"
	"reflect"

	validator "gopkg.in/validator.v2"
)

// board value key
const ( // iota starts at 0
	UNKNOWN   = iota
	VAL_ONE   = iota
	VAL_TWO   = iota
	VAL_THREE = iota
	VOLTORB   = -1
)

// board constraints
const (
	num_rows              = 5
	num_cols              = 5
	num_inputs            = 2
	value_sum_maximum     = 5 * 3
	voltorb_count_maximum = 5 * 1
)

var starting_poss = []int{VOLTORB, VAL_ONE, VAL_TWO, VAL_THREE}

type Tile struct {
	Row   int `validate:"min=0,max=5"`
	Col   int `validate:"min=0,max=5"`
	Value int `validate:"min=-1,max=3"`
	Poss  []int
}

// Custom to string
func (t Tile) String() string {
	if verbose {
		return fmt.Sprintf("%d %v", t.Value, t.Poss)
	}
	return fmt.Sprintf("%d", t.Value)
}

type Board struct {
	Grid       [num_rows][num_cols]Tile
	Vertical   [num_rows][num_inputs]int `validate:"valid_input"`
	Horizontal [num_cols][num_inputs]int `validate:"valid_input"`
	RowStats   [num_rows]Stats
	ColStats   [num_cols]Stats
}

type Stats struct {
	Sum      int `validate:"min=0"`
	Vcount   int `validate:"min=0,max=5"`
	Unkcount int `validate:"min=0,max=5"`
}

// Custom to string
func (b Board) String() string {
	output := ""
	for i := range b.Grid {
		for _, val := range b.Grid[i] {
			if verbose {
				output += fmt.Sprintf("[%v]    ", val)
			} else {
				output += fmt.Sprintf("[%v]\t", val)
			}
		}
		if verbose {
			output += fmt.Sprint(b.Vertical[i])
			output += fmt.Sprint("  ", b.RowStats[i])
		}
		output += fmt.Sprintln("")
	}
	if verbose {
		for _, val := range b.Horizontal {
			output += fmt.Sprintf("[%v]    ", val)
		}
		output += fmt.Sprintln("")
		for _, val := range b.ColStats {
			output += fmt.Sprintf("[%v]  ", val)
		}
	}
	return fmt.Sprint(output)
}

func validate_inputs(v interface{}, param string) error {
	input := reflect.ValueOf(v)
	if input.Kind() != reflect.Array {
		return validator.ErrUnsupported
	}
	var input_arr [num_rows][num_inputs]int = v.([num_rows][num_inputs]int)
	for _, rc := range input_arr {
		sum := rc[0]
		vcount := rc[1]
		if sum < 0 || sum > value_sum_maximum {
			return fmt.Errorf("provided value sum \"%d\" is not within the accepted range: %d - %d", sum, 0, value_sum_maximum)
		} else if vcount < 0 || vcount > voltorb_count_maximum {
			return fmt.Errorf("provided voltorb count \"%d\" is not within the accepted range: %d - %d", vcount, 0, voltorb_count_maximum)
		}
	}
	return nil
}

func generate_board() Board {
	b := Board{
		Grid:       [num_rows][num_cols]Tile{},
		Vertical:   [num_rows][num_inputs]int{},
		Horizontal: [num_cols][num_inputs]int{},
		RowStats:   [num_rows]Stats{},
		ColStats:   [num_cols]Stats{},
	}

	for i := range b.Grid {
		for j := range b.Grid[i] {
			b.Grid[i][j].Row = i
			b.Grid[i][j].Col = j
			b.Grid[i][j].Poss = append([]int{}, starting_poss...)
		}
	}

	return b
}

func (b Board) validate() error {
	for _, i := range b.Grid {
		for _, j := range i {
			if errs := validator.Validate(j); errs != nil {
				return errs
			}
		}
	}

	if errs := validator.Validate(b); errs != nil {
		return errs
	}

	return nil
}
