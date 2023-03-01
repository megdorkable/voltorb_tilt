package main

import "fmt"

// board value key
const ( // iota starts at 0
	unknown   = iota
	val_one   = iota
	val_two   = iota
	val_three = iota
	voltorb   = -1
)

var bkey = map[int]string{voltorb: "v", unknown: "?", val_one: "1", val_two: "2", val_three: "3"}
var bposskey = map[int]int{voltorb: 0, val_one: 1, val_two: 2, val_three: 3}

// input index key
const ( // iota starts at 0
	value_sum_index     = iota
	voltorb_count_index = iota
)

// board constraints
const (
	num_rows              = 5
	num_cols              = 5
	num_inputs            = 2
	num_poss              = 4
	value_sum_maximum     = 5 * 3
	voltorb_count_maximum = 5 * 1
)

func main() {
	var board [num_rows][num_cols]int
	var input_vertical [num_rows][num_inputs]int
	var input_horizontal [num_cols][num_inputs]int
	var board_poss [num_rows][num_cols][num_poss]bool
	for row := range board_poss {
		for col := range board_poss[row] {
			board_poss[row][col] = [num_poss]bool{true, true, true, true}
		}
	}

	board = example_board_valid_incomplete()
	input_vertical = example_input_vertical()
	input_horizontal = example_input_horizontal()

	fmt.Println("\nPrinting Board..")
	print_board(board, input_vertical, input_horizontal)
	fmt.Println("\nValidating..")
	fmt.Println(validate_board(board, input_vertical, input_horizontal))

	fmt.Println("\nResetting Board..")
	reset_board(&board, &input_vertical, &input_horizontal)

	fmt.Println("\nSetting Board Values..")
	fmt.Println(set_board_val(&board, 0, 0, val_one))
	fmt.Println(set_input(&input_vertical, 0, value_sum_index, 1))
	fmt.Println(set_input(&input_vertical, 3, voltorb_count_index, 2))
	fmt.Println(get_board_val(&board, 3, 0))

	fmt.Println("\nPrinting Board..")
	print_board(board, input_vertical, input_horizontal)

	fmt.Println("\nUpdating Board Poss..")
	update_poss(board, &board_poss)
	print_board_poss(board_poss)
}

func print_board(board [num_rows][num_cols]int, input_vertical [num_rows][num_inputs]int, input_horizontal [num_cols][num_inputs]int) {
	var board_str [num_rows][num_cols]string
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			board_str[i][j] = bkey[board[i][j]]
		}
		fmt.Printf("%v  %v\n", board_str[i], input_vertical[i])
	}
	fmt.Println(input_horizontal)
}

func print_board_poss(board_poss [num_rows][num_cols][num_poss]bool) {
	for row := range board_poss {
		fmt.Println(board_poss[row])
	}
}

func reset_board(board *[num_rows][num_cols]int, input_vertical *[num_rows][num_inputs]int, input_horizontal *[num_cols][num_inputs]int) {
	*board = [num_rows][num_cols]int{}
	*input_vertical = [num_rows][num_inputs]int{}
	*input_horizontal = [num_cols][num_inputs]int{}
}

/*
row between 0 and 4 inclusive
col between 0 and 4 inclusive
*/
func validate_board_index(board [num_rows][num_cols]int, row int, col int) (bool, string) {
	// validate index
	if row < 0 || row >= len(board) ||
		col < 0 || col >= len(board[0]) {
		return false, "Index out of bounds"
	}

	return true, ""
}

/*
i between -1 and 3 inclusive
*/
func validate_board_val(i int) (bool, string) {
	// validate i
	if i < -1 || i > 3 {
		return false, "Invalid value"
	}

	return true, ""
}

func get_board_val(board *[num_rows][num_cols]int, row int, col int) (bool, string) {
	var i int

	valid, err := validate_board_index(*board, row, col)
	if !valid {
		return valid, err
	}

	fmt.Printf("Flip the tile at [%d][%d], and type the result: ", row, col)
	fmt.Scan(&i)

	return set_board_val(board, row, col, i)
}

func set_board_val(board *[num_rows][num_cols]int, row int, col int, i int) (bool, string) {
	valid, err := validate_board_index(*board, row, col)
	if !valid {
		return valid, err
	}
	valid, err = validate_board_val(i)
	if !valid {
		return valid, err
	}

	// set
	(*board)[row][col] = i
	return true, ""
}

/*
row_col between 0 and 4 inclusive
val_volt between 0 and 1 inclusive
i, if is val / val_volt == 0, between 0 and value_sum_maximum inclusive
i, if is volt / val_volt == 1, between 0 and voltorb_count_maximum inclusive

** this is dependent on num_rows and num_cols being equal
*/
func set_input(input *[num_rows][num_inputs]int, row_col int, val_volt int, i int) (bool, string) {
	// validate index
	if row_col < 0 || row_col >= len(*input) ||
		val_volt < 0 || val_volt >= len((*input)[0]) {
		return false, "Index out of bounds"
	}
	// validate i
	if val_volt < 0 ||
		(val_volt == 0 && i > value_sum_maximum) ||
		(val_volt == 1 && i > voltorb_count_maximum) {
		return false, "Invalid value"
	}

	// set
	(*input)[row_col][val_volt] = i
	return true, ""
}

func update_poss(board [num_rows][num_cols]int, board_poss *[num_rows][num_cols][num_poss]bool) {
	for row := range board {
		for col := range board[row] {
			val := board[row][col]
			// if the val is known
			if val != 0 {
				// set all other possibilities to false
				for poss := range (*board_poss)[row][col] {
					if poss != bposskey[val] {
						(*board_poss)[row][col][poss] = false
					}
				}
			}
		}
	}
}

func update_board(board *[num_rows][num_cols]int, board_poss [num_rows][num_cols][num_poss]bool) {

}

func validate_board(board [num_rows][num_cols]int, input_vertical [num_rows][num_inputs]int, input_horizontal [num_cols][num_inputs]int) bool {
	var output_vertical [num_rows][3]int
	var output_horizontal [num_cols][3]int

	for row := 0; row < len(board); row++ {
		for column := 0; column < len(board[row]); column++ {
			val := board[row][column]
			// if unknown, else if voltorb, else is value (1, 2, 3)
			if val == 0 {
				output_vertical[row][2] += 1
				output_horizontal[column][2] += 1
			} else if val == -1 {
				output_vertical[row][1] += 1
				output_horizontal[column][1] += 1
			} else {
				output_vertical[row][0] += val
				output_horizontal[column][0] += val
			}
		}
	}

	fmt.Printf("inp_ver: %v\n", input_vertical)
	fmt.Printf("out_ver: %v\n", output_vertical)
	fmt.Printf("inp_hor: %v\n", input_horizontal)
	fmt.Printf("out_hor: %v\n", output_horizontal)

	for row := 0; row < len(output_vertical); row++ {
		input_sum := input_vertical[row][0]
		input_voltorbs := input_vertical[row][1]

		output_sum := output_vertical[row][0]
		output_voltorbs := output_vertical[row][1]
		output_unknown := output_vertical[row][2]

		// if validate voltorbs, else if validate sum isn't too large with 1s, else if validate sum isn't too small with 3s
		if output_voltorbs > input_voltorbs || output_voltorbs+output_unknown < input_voltorbs {
			return false
		} else if output_sum+output_unknown-(input_voltorbs-output_voltorbs) > input_sum {
			return false
		} else if output_sum+(output_unknown*3) < input_sum {
			return false
		}
	}

	return true
}

func example_board_valid_complete() [5][5]int {
	return [5][5]int{
		{voltorb, val_one, val_one, val_two, val_one},
		{val_one, val_one, val_two, val_one, val_one},
		{voltorb, voltorb, voltorb, voltorb, val_two},
		{voltorb, val_one, val_two, val_one, val_two},
		{voltorb, val_one, val_two, voltorb, val_two},
	}
}

func example_input_vertical() [5][num_inputs]int {
	return [5][2]int{
		{5, 1},
		{6, 0},
		{2, 4},
		{6, 1},
		{5, 2},
	}
}

func example_input_horizontal() [5][num_inputs]int {
	return [5][2]int{
		{1, 4},
		{4, 1},
		{7, 1},
		{4, 2},
		{8, 0},
	}
}

func example_board_valid_incomplete() [5][5]int {
	return [5][5]int{
		{unknown, val_one, val_one, val_two, val_one},
		{val_one, val_one, unknown, val_one, val_one},
		{voltorb, voltorb, voltorb, voltorb, unknown},
		{voltorb, unknown, val_two, val_one, val_two},
		{voltorb, val_one, val_two, unknown, val_two},
	}
}

func example_board_invalid() [5][5]int {
	return [5][5]int{
		{0, 0, 1, 2, 2},
		{1, 1, 0, 1, 1},
		{-1, -1, -1, -1, 0},
		{-1, 0, 2, 1, 2},
		{-1, 1, 2, 0, 2},
	}
}
