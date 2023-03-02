package main

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
