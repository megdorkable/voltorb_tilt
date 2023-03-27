package main

import (
	"fmt"

	validator "gopkg.in/validator.v2"
)

const (
	verbose = false
)

func main() {
	validator.SetValidationFunc("valid_input", validate_inputs)

	b := generate_random_board()
	// fmt.Println(b)

	// b := generate_board()
	// b.Vertical = example_input_vertical()
	// b.Horizontal = example_input_horizontal()

	// if errs := b.validate(); errs != nil {
	// 	fmt.Println(errs)
	// }

	solved := b.update()
	fmt.Println(b)

	if solved {
		fmt.Println("Solved!")
	} else {
		fmt.Println("Oops, try again.")
	}
}
