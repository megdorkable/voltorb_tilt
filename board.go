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
	num_poss              = 4 // depreciated
	value_sum_maximum     = 5 * 3
	voltorb_count_maximum = 5 * 1
)

var starting_poss = []int{VOLTORB, VAL_ONE, VAL_TWO, VAL_THREE}

type Tile struct {
	Row   int `validate:"min=0,max=5"`
	Col   int `validate:"min=0,max=5"`
	Value int `validate:"min=-1,max=3"`
	// Poss  [5]bool
	Poss []int
}

// Custom to string
func (t Tile) String() string {
	return fmt.Sprintf("%d %v", t.Value, t.Poss)
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
			output += fmt.Sprintf("[%v]  ", val)
		}
		output += fmt.Sprintln(b.Vertical[i])
	}
	for _, inp := range b.Horizontal {
		output += fmt.Sprint(inp)
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

func (b *Board) update() bool {
	updated := false

	b.refresh_stats()
	for row := range b.Grid {
		updated = b.update_row_poss(row) || updated
	}
	for col := range b.Grid[0] {
		updated = b.update_col_poss(col) || updated
	}

	return updated
}

func (b *Board) refresh_stats() {
	b.RowStats = [num_rows]Stats{}
	b.ColStats = [num_cols]Stats{}

	for row := range b.Grid {
		for col, val := range b.Grid[row] {
			if val.Value == -1 {
				b.RowStats[row].Vcount += 1
				b.ColStats[col].Vcount += 1
			} else if val.Value == 0 {
				b.RowStats[row].Unkcount += 1
				b.ColStats[col].Unkcount += 1
			} else {
				b.RowStats[row].Sum += val.Value
				b.ColStats[col].Sum += val.Value
			}
		}
	}
}

func (b *Board) update_row_poss(row int) bool {
	updated := false

	row_total_sum := b.Vertical[row][0]
	row_total_vcount := b.Vertical[row][1]
	row_rem_sum := row_total_sum - b.RowStats[row].Sum
	row_rem_vcount := row_total_vcount - b.RowStats[row].Vcount
	row_rem_unknown := b.RowStats[row].Unkcount

	// if remaining sum + remaining number voltorb == remaining number unknown, must be VOLTORB OR VAL_ONE
	if row_rem_sum+row_rem_vcount == row_rem_unknown {
		for col := range b.Grid[row] {
			if b.Grid[row][col].Value == UNKNOWN {
				new_poss := []int{}
				for _, pval := range b.Grid[row][col].Poss {
					if pval == VOLTORB || pval == VAL_ONE {
						new_poss = append(new_poss, pval)
					}
				}
				(*b).Grid[row][col].Poss = new_poss
				updated = true
			}
		}
	}

	fmt.Println(b)

	// if no remaining voltorbs, can't be VOLTORB
	if row_rem_vcount == 0 {
		for col := range b.Grid[row] {
			if b.Grid[row][col].Value == UNKNOWN {
				new_poss := []int{}
				for _, pval := range b.Grid[row][col].Poss {
					if pval != -1 {
						new_poss = append(new_poss, pval)
					}
				}
				(*b).Grid[row][col].Poss = new_poss
				updated = true
			}
		}
	}

	return updated
}

func (b *Board) update_col_poss(col int) bool {
	updated := false

	col_total_sum := b.Horizontal[col][0]
	col_total_vcount := b.Horizontal[col][1]
	col_rem_sum := col_total_sum - b.ColStats[col].Sum
	col_rem_vcount := col_total_vcount - b.ColStats[col].Vcount
	col_rem_unknown := b.ColStats[col].Unkcount

	if col_rem_sum+col_rem_vcount == col_rem_unknown {
		for row := range b.Grid {
			if b.Grid[row][col].Value == UNKNOWN {
				new_poss := []int{}
				for _, pval := range b.Grid[row][col].Poss {
					if pval == VOLTORB || pval == VAL_ONE {
						new_poss = append(new_poss, pval)
					}
				}
				(*b).Grid[row][col].Poss = new_poss
				updated = true
			}
		}
	}

	return updated
}

func remove(s *[]int, i int) {
	fmt.Println(s)
	(*s)[i] = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
}

func main() {
	validator.SetValidationFunc("valid_input", validate_inputs)

	b := generate_board()
	b.Vertical = example_input_vertical()
	b.Horizontal = example_input_horizontal()
	// b.Vertical = example_input_horizontal()
	// b.Horizontal = example_input_vertical()

	if errs := b.validate(); errs != nil {
		fmt.Println(errs)
	}
	// else {
	// 	fmt.Println(b)
	// }

	fmt.Println(b.update())
	fmt.Println(b.RowStats)
	fmt.Println(b.ColStats)
	fmt.Println("")

	fmt.Println(b)
}
