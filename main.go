package main

import (
	"fmt"

	validator "gopkg.in/validator.v2"
)

const (
	verbose         = false
	generate_random = false
)

func main() {
	validator.SetValidationFunc("valid_input", validate_inputs)

	var b Board

	if generate_random {
		generated := generate_random_board()
		b = generated[0]
		solution := generated[1]

		fmt.Println("Generated Board:")
		fmt.Println(solution)
	} else {
		b = generate_board()
		b.Vertical = example_input_vertical()
		b.Horizontal = example_input_horizontal()

		if errs := b.validate(); errs != nil {
			fmt.Println(errs)
		}
	}

	solved := b.update()
	fmt.Println(b)

	if solved {
		fmt.Println("Solved!")
	} else {
		fmt.Println("Oops, try again.")
	}
}
