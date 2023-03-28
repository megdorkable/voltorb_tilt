package main

import (
	"flag"
	"fmt"

	validator "gopkg.in/validator.v2"
)

const (
	verbose    = false
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
)

func main() {
	validator.SetValidationFunc("valid_input", validate_inputs)

	var generate_random bool
	flag.BoolVar(&generate_random, "g", false, "generate a random board")

	flag.Parse()

	var b Board

	if generate_random {
		generated := generate_random_board()
		b = generated[0]
		solution := generated[1]

		fmt.Println(string(colorGreen) + "Generated Board:")
		fmt.Println(solution, string(colorReset))
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
