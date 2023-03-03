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

	b := generate_board()
	b.Vertical = example_input_vertical()
	b.Horizontal = example_input_horizontal()

	if errs := b.validate(); errs != nil {
		fmt.Println(errs)
	}

	fmt.Println(b.update())
	fmt.Println(b)
}
